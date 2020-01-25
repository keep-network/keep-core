import {initContracts} from './helpers/initContracts'
import stakeDelegate from './helpers/stakeDelegate'
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"
import {bls} from './helpers/data'
import expectThrowWithMessage from './helpers/expectThrowWithMessage'

contract('KeepRandomBeaconOperator', function(accounts) {
  let token, stakingContract, serviceContract, operatorContract, minimumStake, largeStake, entryFeeEstimate, groupIndex,
    owner = accounts[0],
    operator1 = accounts[1],
    operator2 = accounts[2],
    operator3 = accounts[3],
    tattletale = accounts[4]

  before(async () => {
    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')
    )

    token = contracts.token
    stakingContract = contracts.stakingContract
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract

    groupIndex = 0
    await operatorContract.registerNewGroup(bls.groupPubKey)
    await operatorContract.addGroupMember(bls.groupPubKey, operator1)
    await operatorContract.addGroupMember(bls.groupPubKey, operator2)
    await operatorContract.addGroupMember(bls.groupPubKey, operator3)

    minimumStake = await operatorContract.minimumStake()
    largeStake = minimumStake.muln(2)
    await stakeDelegate(stakingContract, token, owner, operator1, owner, owner, largeStake)
    await stakeDelegate(stakingContract, token, owner, operator2, owner, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, owner, owner, minimumStake)

    entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to report unauthorized signing", async () => {
    await operatorContract.reportUnauthorizedSigning(
      groupIndex,
      bls.signedGroupPubKey,
      {from: tattletale}
    )

    assert.isTrue((await stakingContract.balanceOf(operator1)).eq(largeStake.sub(minimumStake)),"Unexpected operator 1 balance")
    assert.isTrue((await stakingContract.balanceOf(operator2)).isZero(), "Unexpected operator 2 balance")
    assert.isTrue((await stakingContract.balanceOf(operator3)).isZero(), "Unexpected operator 3 balance")
    
    // Expecting 5% of all the seized tokens
    let expectedTattletaleReward = minimumStake.muln(3).muln(5).divn(100)
    assert.isTrue((await token.balanceOf(tattletale)).eq(expectedTattletaleReward), "Unexpected tattletale balance")

    // Group should be terminated, expecting total number of groups to become 0
    await expectThrowWithMessage(
      serviceContract.methods['requestRelayEntry()']({value: entryFeeEstimate}),
      "Total number of groups must be greater than zero."
    )
  })

  it("should ignore invalid report of unauthorized signing", async () => {
    await operatorContract.reportUnauthorizedSigning(
      groupIndex,
      bls.nextGroupSignature, // Wrong signature
      {from: tattletale}
    )

    assert.isTrue((await stakingContract.balanceOf(operator1)).eq(largeStake), "Unexpected operator 1 balance")
    assert.isTrue((await stakingContract.balanceOf(operator2)).eq(minimumStake), "Unexpected operator 2 balance")
    assert.isTrue((await stakingContract.balanceOf(operator3)).eq(minimumStake), "Unexpected operator 3 balance")

    assert.isTrue((await token.balanceOf(tattletale)).isZero(), "Unexpected tattletale balance")
  })
})
