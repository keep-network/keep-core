import {createSnapshot, restoreSnapshot} from './helpers/snapshot'
import expectThrowWithMessage from './helpers/expectThrowWithMessage'
import {initContracts} from './helpers/initContracts'
const RegistryKeeper = artifacts.require('./RegistryKeeper.sol')

contract('RegistryKeeper', function(accounts) {

  let registryKeeper, operatorContract,
    governance = accounts[0],
    panicButton = accounts[1],
    operatorContractUpgrader = accounts[2]

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      artifacts.require('./KeepRandomBeaconOperator.sol')
    );

    stakingContract = contracts.stakingContract
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract

    registryKeeper = await RegistryKeeper.new(panicButton)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to approve or disable operator contract", async() => {
    assert.isTrue((await registryKeeper.operatorContracts(operatorContract.address)).eqn(0), "Unexpected status of operator contract")

    await expectThrowWithMessage(
      registryKeeper.approveOperatorContract(operatorContract.address, {from: operatorContractUpgrader}),
      "Ownable: caller is not the owner"
    );

    await registryKeeper.approveOperatorContract(operatorContract.address, {from: governance})
    assert.isTrue((await registryKeeper.operatorContracts(operatorContract.address)).eqn(1), "Unexpected status of operator contract")

    await expectThrowWithMessage(
      registryKeeper.disableOperatorContract(operatorContract.address, {from: governance}),
      "Not authorized"
    );

    await registryKeeper.disableOperatorContract(operatorContract.address, {from: panicButton})
    assert.isTrue((await registryKeeper.operatorContracts(operatorContract.address)).eqn(2), "Unexpected status of operator contract")
  })
})
