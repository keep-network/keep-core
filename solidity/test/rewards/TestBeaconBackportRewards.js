const { initContracts } = require('../helpers/initContracts')
const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
const crypto = require("crypto")
const stakeDelegate = require('../helpers/stakeDelegate')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('BeaconBackportRewards', () => {

    // accounts
    const owner = accounts[0]
    let operators
    let beneficiaries

    // contracts
    let rewardsContract
    let token, stakingContract, operatorContract

    // system parameters
    const tokenDecimalMultiplier = web3.utils.toBN(10).pow(web3.utils.toBN(18))
    // 1,000,000,000 - total KEEP supply
    //   200,000,000 - 20% of the total supply goes to staker rewards
    //    20,000,000 - 10% of staker rewards goes to the random beacon stakers
    //       200,000 - 1% of staker rewards for the beacon goes to May genesis groups
    const totalBeaconRewards = web3.utils.toBN(200000).mul(tokenDecimalMultiplier)

    before(async() => {
        let contracts = await initContracts(
            contract.fromArtifact('TokenStaking'),
            contract.fromArtifact('KeepRandomBeaconService'),
            contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
            contract.fromArtifact('KeepRandomBeaconOperatorBeaconRewardsStub')
        )

        token = contracts.token
        stakingContract = contracts.stakingContract
        operatorContract = contracts.operatorContract

        rewardsContract = await contract.fromArtifact('BeaconBackportRewards').new(
            token.address,
            operatorContract.address,
            stakingContract.address,
            { from: owner }
        )

        await token.approveAndCall(
            rewardsContract.address,
            totalBeaconRewards,
            "0x0",
            { from: owner }
        )
        await rewardsContract.markAsFunded({from: owner})

        // create 64 operators and beneficiaries, delegate stake for them
        const minimumStake = await stakingContract.minimumStake()
        operators = []
        beneficiaries = []
        for (i = 0; i < 64; i++) {
            const operator = accounts[i]
            const beneficiary = accounts[64 + i]
            const authorizer = operator

            operators.push(operator)
            beneficiaries.push(beneficiary)
            await stakeDelegate(
                stakingContract, 
                token, 
                owner, 
                operator, 
                beneficiary, 
                authorizer, 
                minimumStake
            )
        }

        // 3 groups created in an interval
        const startOf = await rewardsContract.startOf(0)
        await registerNewGroup(startOf.addn(1))
        await registerNewGroup(startOf.addn(2))
        await registerNewGroup(startOf.addn(3))
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("interval allocation", async() => {
        it("should equal the full allocation", async() => {
            const expectedAllocation = 200000

            await timeJumpToEndOfInterval(0)
            await rewardsContract.allocateRewards(0)

            const allocated = await rewardsContract.getAllocatedRewards(0)                                
            const allocatedKeep = allocated.div(tokenDecimalMultiplier)
            
            expect(allocatedKeep).to.eq.BN(expectedAllocation)
        })
    })

    describe("rewards withdrawal", async () => {
        it("should correctly distribute rewards to beneficiaries", async () => {
            await timeJumpToEndOfInterval(0)

            await rewardsContract.receiveReward(0)
            // each beneficiary receives 200000 / 3 / 64 = 1041 KEEP
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(1041))
         
            await rewardsContract.receiveReward(1)
              // each beneficiary receives 200000 / 3 / 64 = 1041 KEEP
            // they should have 1041 + 1041 = 2082 KEEP now
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(2082))

            await rewardsContract.receiveReward(2)
            // each beneficiary receives 200000 / 3 / 64 = 1041 KEEP
            // they should have 1041 + 1041 + 1041 = 3123 KEEP now
            await assertKeepBalanceOfBeneficiaries(web3.utils.toBN(3123))
        })

        it("should fail for non-existing group", async () => {
            await expectRevert(
                rewardsContract.receiveReward(3),
                "Keep not recognized by factory"
            )
        })
    })

    async function timeJumpToEndOfInterval(intervalNumber) {
        const endOf = await rewardsContract.endOf(intervalNumber)
        const now = await time.latest()

        if (now.lt(endOf)) {
            await time.increaseTo(endOf.addn(1))
        }
    }

    async function registerNewGroup(creationTimestamp) {
        const groupPublicKey = crypto.randomBytes(128)
        await operatorContract.registerNewGroup(groupPublicKey, operators, creationTimestamp)
    }

    async function assertKeepBalanceOfBeneficiaries(expectedBalance) {
        // solidity is not very good when it comes to floating point precision,
        // we are allowing for ~2 KEEP difference margin between expected and
        // actual value
        const precision = 2

        for (let i = 0; i < beneficiaries.length; i++) {
            const balance = await token.balanceOf(beneficiaries[i])
            const balanceInKeep = balance.div(tokenDecimalMultiplier)

            expect(balanceInKeep).to.gte.BN(expectedBalance.subn(precision))
            expect(balanceInKeep).to.lte.BN(expectedBalance.addn(precision))
        }
    }
})
