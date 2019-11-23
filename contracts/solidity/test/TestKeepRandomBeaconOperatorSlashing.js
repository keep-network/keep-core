import {initContracts} from './helpers/initContracts'
import stakeDelegate from './helpers/stakeDelegate'
import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"
import {bls} from './helpers/data'

contract('KeepRandomBeaconOperator', function(accounts) {
  let token, stakingContract, serviceContract, operatorContract, minimumStake,
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

    await operatorContract.registerNewGroup(bls.groupPubKey)
    await operatorContract.addGroupMember(bls.groupPubKey, operator1)
    await operatorContract.addGroupMember(bls.groupPubKey, operator2)
    await operatorContract.addGroupMember(bls.groupPubKey, operator3)

    minimumStake = await operatorContract.minimumStake()
    await stakeDelegate(stakingContract, token, owner, operator1, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator2, owner, minimumStake)
    await stakeDelegate(stakingContract, token, owner, operator3, owner, minimumStake)

    let entryFeeEstimate = await serviceContract.entryFeeEstimate(0)
    await serviceContract.requestRelayEntry(bls.seed, {value: entryFeeEstimate})
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to report unauthorized signing", async () => {
    await operatorContract.reportUnauthorizedSigning(
      bls.groupPubKey,
      bls.signedGroupPubKey,
      {from: tattletale}
    )

    assert.equal((await stakingContract.balanceOf(operator1)).isZero(), true, "Unexpected operator 1 balance")
    assert.equal((await stakingContract.balanceOf(operator2)).isZero(), true, "Unexpected operator 2 balance")
    assert.equal((await stakingContract.balanceOf(operator3)).isZero(), true, "Unexpected operator 3 balance")
    
    // Expecting 5% of all the seized tokens
    let expectedTattletaleReward = minimumStake.muln(3).muln(5).divn(100)
    assert.isTrue((await token.balanceOf(tattletale)).eq(expectedTattletaleReward), "Unexpected tattletale balance")
  })
})
