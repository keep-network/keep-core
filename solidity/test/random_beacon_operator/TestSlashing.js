const blsData = require("../helpers/data.js")
const {initContracts} = require('../helpers/initContracts')
const { createSnapshot, restoreSnapshot } = require("../helpers/snapshot.js")
const { contract, accounts, web3 } = require("@openzeppelin/test-environment")
const { expectRevert, time } = require("@openzeppelin/test-helpers")
const stakeDelegate = require('../helpers/stakeDelegate')
const BLS = contract.fromArtifact('BLS');

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('KeepRandomBeaconOperator/Slashing', function () {
  let token, stakingContract, serviceContract, operatorContract, entryFeeEstimate, groupIndex,
    registry, bls,
    owner = accounts[0],
    operator1 = accounts[1],
    operator2 = accounts[2],
    operator3 = accounts[3],
    tattletale = accounts[4],
    authorizer = accounts[5],
    anotherOperatorContract = accounts[6],
    registryKeeper = accounts[7];

  let scheduleStart
  let relayRequestStartBlock

  const largeStake = web3.utils.toBN("50000000000000000000000000") // 50 000 000 KEEP
  const mediumStake = web3.utils.toBN("500000000000000000000000") // 500 000 KEEP
  const smallStake = web3.utils.toBN("100000000000000000000000") // 100 000 KEEP

  before(async () => {

    let contracts = await initContracts(
      contract.fromArtifact('TokenStakingStub'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorSlashingStub')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract
    registry = contracts.registry
    bls = await BLS.new()      

    await registry.setRegistryKeeper(registryKeeper, { from: accounts[0] })
    await registry.approveOperatorContract(anotherOperatorContract, { from: registryKeeper })

    await stakeDelegate(stakingContract, token, owner, operator1, owner, authorizer, largeStake)
    await stakeDelegate(stakingContract, token, owner, operator2, owner, authorizer, mediumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, owner, authorizer, smallStake)
    await stakingContract.authorizeOperatorContract(operator1, operatorContract.address, { from: authorizer })
    await stakingContract.authorizeOperatorContract(operator2, operatorContract.address, { from: authorizer })
    await stakingContract.authorizeOperatorContract(operator3, operatorContract.address, { from: authorizer })

    scheduleStart = await stakingContract.minimumStakeScheduleStart()

    time.increase((await stakingContract.initializationPeriod()).addn(1))

    groupIndex = 0
    await operatorContract.registerNewGroup(blsData.groupPubKey, [operator1, operator2, operator3])
    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    let tx = await serviceContract.methods['requestRelayEntry()']({ value: entryFeeEstimate, from: accounts[0] })
    relayRequestStartBlock = web3.utils.toBN(tx.receipt.blockNumber)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("reportUnauthorizedSigning", async () => {
    it("seizes 100% of minimum stake", async () => {
      let tattletaleSignature = await bls.sign(tattletale, blsData.secretKey);
  
      await operatorContract.reportUnauthorizedSigning(
        groupIndex,
        tattletaleSignature,
        { from: tattletale }
      )
  
      expect(await stakingContract.balanceOf(operator1)).to.eq.BN(
        "49900000000000000000000000" // 50000000000000000000000000 - 100000000000000000000000
      )
      expect(await stakingContract.balanceOf(operator2)).to.eq.BN(
        "400000000000000000000000" // 500000000000000000000000 - 100000000000000000000000
      )
      expect(await stakingContract.balanceOf(operator3)).to.eq.BN(
        "0" // 100000000000000000000000 - 100000000000000000000000
      )
  
      // Expecting 5% of all the seized tokens
      //
      // minimum stake = 100000000000000000000000
      // 3 * 100000000000000000000000 * 5% = 15000000000000000000000
      expect(await token.balanceOf(tattletale)).to.eq.BN("15000000000000000000000")
  
      // Group should be terminated, expecting total number of groups to become 0
      await expectRevert(
        serviceContract.methods['requestRelayEntry()']({ value: entryFeeEstimate, from: accounts[0] }),
        "Total number of groups must be greater than zero."
      )
    })

    it("reverts for invalid signature", async () => {
      await expectRevert(
        operatorContract.reportUnauthorizedSigning(
          groupIndex,
          blsData.nextGroupSignature, // Wrong signature
          { from: tattletale }
        ),
        "Invalid signature"
      )
    })

    it("reverts when already reported for the group", async () => {
      let tattletaleSignature = await bls.sign(tattletale, blsData.secretKey);
  
      await operatorContract.reportUnauthorizedSigning(
        groupIndex,
        tattletaleSignature,
        { from: tattletale }
      )
        
      await expectRevert(
        operatorContract.reportUnauthorizedSigning(
          groupIndex,
          tattletaleSignature,
          { from: tattletale }
        ),
        "Group has been already terminated"
      )
    })
  })

  describe("reportRelayEntryTimeout", async () => {
    it("reverts if entry did not time out", async () => {
      await expectRevert(
        operatorContract.reportRelayEntryTimeout({ from: tattletale }),
        "Entry did not time out."
      )

      await time.advanceBlockTo(relayRequestStartBlock.addn(9));

      await expectRevert(
        operatorContract.reportRelayEntryTimeout({ from: tattletale }),
        "Entry did not time out."
      )
    })

    // There is only one active group in the system and that group did not
    // produce relay entry on time. Relay entry timeout is reported but since
    // there is no other group in the system, we do not retry with another
    // group. This way, relay entry is timed out until someone requests for
    // another entry and we need to make sure the group cannot be slashed more
    // than one time.
    it("reverts when already reported for the last active group", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })
      await expectRevert(
        operatorContract.reportRelayEntryTimeout({ from: tattletale }),
        "Group has been already terminated"
      )
    })

    it("does not revert in the first block relay entry timed out", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })
      // ok, no reverts
    })

    it("seizes 1% of minimum stake from operators at the beginning", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // minimum stake = 100000000000000000000000
      expect(await stakingContract.balanceOf(operator1)).to.eq.BN(
        "49999000000000000000000000" // 50000000000000000000000000 - 1% * 100000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator2)).to.eq.BN(
        "499000000000000000000000"  // 500000000000000000000000 - 1% * 100000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator3)).to.eq.BN(
        "99000000000000000000000"  // 100000000000000000000000 - 1% * 100000000000000000000000 
      )
    })

    it("rewards tattletale with 1% stake adjustment at the beginning", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // Expecting 5% of all the seized tokens with reward adjustment of (20 / 64) = 31%.
      // And "all of the seized tokens" are 3 * minimum stake with 1% adjustment
      // for the first three months:
      //
      // minimum stake = 100000000000000000000000
      // 3 * 100000000000000000000000 * 1% * 5% * 31% = 46500000000000000000
      expect(await token.balanceOf(tattletale)).to.eq.BN("46500000000000000000")
    })

    it("seizes 1% of minimum stake from operators before the first 3 months end", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 89))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // minimum stake = 90000000000000000000000
      expect(await stakingContract.balanceOf(operator1)).to.eq.BN(
        "49999100000000000000000000" // 50000000000000000000000000 - 1% * 90000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator2)).to.eq.BN(
        "499100000000000000000000"  // 500000000000000000000000 - 1% * 90000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator3)).to.eq.BN(
        "99100000000000000000000"  // 100000000000000000000000 - 1% * 90000000000000000000000
      )
    })

    it("rewards tattletale with 1% stake adjustment before the first 3 months end", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 89))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // Expecting 5% of all the seized tokens with reward adjustment of (20 / 64) = 31%.
      // And "all of the seized tokens" are 3 * minimum stake with 1% adjustment
      // for the first three months:
      //
      // minimum stake = 90000000000000000000000
      // 3 * 90000000000000000000000 * 1% * 5% * 31% = 41850000000000000000
      expect(await token.balanceOf(tattletale)).to.eq.BN("41850000000000000000")
    })

    it("seizes 50% of minimum stake from operators after the first 3 months", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 90))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // minimum stake = 90000000000000000000000
      expect(await stakingContract.balanceOf(operator1)).to.eq.BN(
        "49955000000000000000000000" // 50000000000000000000000000 - 50% * 90000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator2)).to.eq.BN(
        "455000000000000000000000"  // 500000000000000000000000 - 50% * 90000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator3)).to.eq.BN(
        "55000000000000000000000"  // 100000000000000000000000 - 50% * 90000000000000000000000
      )
    })

    it("rewards tattletale with 50% stake adjustment after the first 3 months", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 90))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // Expecting 5% of all the seized tokens with reward adjustment of (20 / 64) = 31%.
      // And "all of the seized tokens" are 3 * minimum stake with 50% adjustment
      // after the first 3 months end
      //
      // minimum stake = 90000000000000000000000
      // 3 * 90000000000000000000000 * 50% * 5% * 31% = 2092500000000000000000
      expect(await token.balanceOf(tattletale)).to.eq.BN("2092500000000000000000")
    })
    
    it("seizes 50% of minimum stake from operators before the first 6 months end", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 179))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // minimum stake = 80000000000000000000000
      expect(await stakingContract.balanceOf(operator1)).to.eq.BN(
        "49960000000000000000000000" // 50000000000000000000000000 - 50% * 80000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator2)).to.eq.BN(
        "460000000000000000000000"  // 500000000000000000000000 - 50% * 80000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator3)).to.eq.BN(
        "60000000000000000000000"  // 100000000000000000000000 - 50% * 80000000000000000000000
      )
    })

    it("rewards tattletale with 50% stake adjustment before the first 6 months end", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 179))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // Expecting 5% of all the seized tokens with reward adjustment of (20 / 64) = 31%.
      // And "all of the seized tokens" are 3 * minimum stake with 1% adjustment
      // for the first three months:
      //
      // minimum stake = 80000000000000000000000
      // 3 * 80000000000000000000000 * 50% * 5% * 31% = 1860000000000000000000
      expect(await token.balanceOf(tattletale)).to.eq.BN("1860000000000000000000")
    })

    it("seizes 100% of minimum stake from operators after the first 6 months end", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 180))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // minimum stake = 80000000000000000000000
      expect(await stakingContract.balanceOf(operator1)).to.eq.BN(
        "49920000000000000000000000" // 50000000000000000000000000 - 100% * 80000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator2)).to.eq.BN(
        "420000000000000000000000"  // 500000000000000000000000 - 100% * 80000000000000000000000 
      )
      expect(await stakingContract.balanceOf(operator3)).to.eq.BN(
        "20000000000000000000000"  // 100000000000000000000000 - 100% * 80000000000000000000000
      )
    })

    it("rewards tattletale with 100% stake adjustment after the first 6 months end", async () => {
      await time.advanceBlockTo(relayRequestStartBlock.addn(10));
      await time.increaseTo(scheduleStart.addn(86400 * 180))
      await operatorContract.reportRelayEntryTimeout({ from: tattletale })

      // Expecting 5% of all the seized tokens with reward adjustment of (20 / 64) = 31%.
      // And "all of the seized tokens" are 3 * minimum stake with 1% adjustment
      // for the first three months:
      //
      // minimum stake = 80000000000000000000000
      // 3 * 80000000000000000000000 * 100% * 5% * 31% = 3720000000000000000000
      expect(await token.balanceOf(tattletale)).to.eq.BN("3720000000000000000000")
    })
  })
})
