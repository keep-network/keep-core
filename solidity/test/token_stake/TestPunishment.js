const { contract, accounts, web3 } = require("@openzeppelin/test-environment")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
const { createSnapshot, restoreSnapshot } = require('../helpers/snapshot');
const {initTokenStaking} = require('../helpers/initContracts')
const stakeDelegate = require('../helpers/stakeDelegate')

const KeepToken = contract.fromArtifact('KeepToken');
const TokenGrant = contract.fromArtifact('TokenGrant');
const KeepRegistry = contract.fromArtifact("KeepRegistry");

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('TokenStaking/Punishment', () => {
    let token, registry, stakingContract;

    let owner = accounts[0],
        registryKeeper = accounts[1],
        operator = accounts[2],
        authorizer = accounts[3],
        operatorContract = accounts[4],
        tattletale = accounts[5];

    let largeStake, minimumStake;

    const initializationPeriod = time.duration.seconds(10)
    let undelegationPeriod

    before(async () => {
        token = await KeepToken.new({ from: owner })
        tokenGrant = await TokenGrant.new(token.address,  {from: owner})
        registry = await KeepRegistry.new({ from: owner })
        const stakingContracts = await initTokenStaking(
            token.address,
            tokenGrant.address,
            registry.address,
            initializationPeriod,
            contract.fromArtifact('TokenStakingEscrow'),
            contract.fromArtifact('TokenStakingStub')
        )
        stakingContract = stakingContracts.tokenStaking;

        undelegationPeriod = await stakingContract.undelegationPeriod()

        await registry.setRegistryKeeper(registryKeeper, { from: owner })

        minimumStake = await stakingContract.minimumStake()
        largeStake = minimumStake.muln(2)

        await registry.approveOperatorContract(
            operatorContract,
            { from: registryKeeper }
        )

        await stakeDelegate(
            stakingContract, token, owner, operator,
            owner, authorizer, largeStake
        )

        await stakingContract.authorizeOperatorContract(
            operator,
            operatorContract,
            { from: authorizer }
        )
    });

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("slash", () => {
        it("should slash token amount from stake", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))

            let amountToSlash = web3.utils.toBN(42000000);

            let balanceBeforeSlashing = await stakingContract.balanceOf(operator)
            await stakingContract.slash(amountToSlash, [operator], { from: operatorContract })
            let balanceAfterSlashing = await stakingContract.balanceOf(operator)

            expect(balanceAfterSlashing).to.eq.BN(balanceBeforeSlashing.sub(amountToSlash))
        })

        it("should slash no more than available on stake", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))

            let amountToSlash = largeStake.add(web3.utils.toBN(100))
            await stakingContract.slash(amountToSlash, [operator], { from: operatorContract })
            let balanceAfterSlashing = await stakingContract.balanceOf(operator)

            expect(balanceAfterSlashing).to.eq.BN(0)
        })

        it("should not fail if operator is slashed to zero", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))

            let amountToSlash = largeStake

            // the first slash will slash to 0, the second one has nothing
            // to slash; it should not fail
            await stakingContract.slash(amountToSlash, [operator], { from: operatorContract })
            await stakingContract.slash(amountToSlash, [operator], { from: operatorContract })

            let balanceAfterSlashing = await stakingContract.balanceOf(operator)
            expect(balanceAfterSlashing).to.eq.BN(0)
        })

        it("should fail when operator stake is not active yet", async () => {
            let amountToSlash = web3.utils.toBN(1000)
            await expectRevert(
                stakingContract.slash(amountToSlash, [operator], { from: operatorContract }),
                "Inactive stake"
            )
        })

        it("should fail when operator stake is released", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))
            await stakingContract.undelegate(operator, { from: owner })
            time.increase((await stakingContract.undelegationPeriod()).addn(1))

            let amountToSlash = web3.utils.toBN(100);
            await expectRevert(
                stakingContract.slash(amountToSlash, [operator], { from: operatorContract }),
                "Stake is released"
            )
        })
    })

    describe("seize", () => {
        it("should seize token amount from stake", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))

            let operatorBalanceBeforeSeizing = await stakingContract.balanceOf(operator)
            let tattletaleBalanceBeforeSeizing = await token.balanceOf(tattletale)

            let amountToSeize = web3.utils.toBN(42000000)
            let rewardMultiplier = web3.utils.toBN(25)
            await stakingContract.seize(
                amountToSeize, rewardMultiplier, tattletale,
                [operator], { from: operatorContract }
            )

            let operatorBalanceAfterSeizing = await stakingContract.balanceOf(operator)
            let tattletaleBalanceAfterSeizing = await token.balanceOf(tattletale)

            expect(operatorBalanceAfterSeizing).to.eq.BN(
                operatorBalanceBeforeSeizing.sub(amountToSeize)
            )

            // 525000 = (42000000 * 5 / 100) * 25 / 100
            let expectedTattletaleReward = web3.utils.toBN(525000)
            expect(tattletaleBalanceAfterSeizing).to.eq.BN(
                tattletaleBalanceBeforeSeizing.add(expectedTattletaleReward)
            )
        })

        it("should seize no more than available on stake", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))

            let tattletaleBalanceBeforeSeizing = await token.balanceOf(tattletale)

            // we test with a higher excess to ensure that the tattletale reward
            // is calculated from the applied penalty, not the requested penalty            
            let amountToSeize = largeStake.muln(2) // 400000000000000000000000
            let rewardMultiplier = web3.utils.toBN(10)
            await stakingContract.seize(
                amountToSeize, rewardMultiplier, tattletale,
                [operator], { from: operatorContract }
            )

            let operatorBalanceAfterSeizing = await stakingContract.balanceOf(operator)
            let tattletaleBalanceAfterSeizing = await token.balanceOf(tattletale)

            expect(operatorBalanceAfterSeizing).to.eq.BN(0)

            // 1000000000000000000000 = (200000000000000000000000 * 5 / 100) * 10 / 100
            let expectedTattletaleReward = web3.utils.toBN("1000000000000000000000")
            expect(tattletaleBalanceAfterSeizing).to.eq.BN(
                tattletaleBalanceBeforeSeizing.add(expectedTattletaleReward)
            )
        })

        it("should not fail if operator is slashed to zero", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))

            let amountToSlash = largeStake
            let amountToSeize = largeStake
            let rewardMultiplier = web3.utils.toBN(10)

            let tattletaleBalanceBeforeSeizing = await token.balanceOf(tattletale)

            // the first slash will slash to 0, the seize happening later
            // should not fail
            await stakingContract.slash(amountToSlash, [operator], { from: operatorContract })
            await stakingContract.seize(
                amountToSeize, rewardMultiplier, tattletale,
                [operator], { from: operatorContract }
            )

            let operatorBalanceAfterSeizing = await stakingContract.balanceOf(operator)
            let tattletaleBalanceAfterSeizing = await token.balanceOf(tattletale)

            expect(operatorBalanceAfterSeizing).to.eq.BN(0)
            expect(tattletaleBalanceAfterSeizing).to.eq.BN(
                tattletaleBalanceBeforeSeizing
            )
        })

        it("should fail when operator stake is not active yet", async () => {
            let amountToSeize = web3.utils.toBN(42000000)
            let rewardMultiplier = web3.utils.toBN(25)
            await expectRevert(
                stakingContract.seize(
                    amountToSeize, rewardMultiplier, tattletale,
                    [operator], { from: operatorContract }
                ),
                "Inactive stake"
            )
        })

        it("should fail when operator stake is released", async () => {
            time.increase((await stakingContract.initializationPeriod()).addn(1))
            await stakingContract.undelegate(operator, { from: owner })
            time.increase((await stakingContract.undelegationPeriod()).addn(1))

            let amountToSeize = web3.utils.toBN(10000);
            let rewardMultiplier = web3.utils.toBN(25)
            await expectRevert(
                stakingContract.seize(
                    amountToSeize, rewardMultiplier, tattletale,
                    [operator], { from: operatorContract }
                ),
                "Stake is released"
            )
        })
    })
})