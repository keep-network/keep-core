import {createSnapshot, restoreSnapshot} from './helpers/snapshot'
import expectThrowWithMessage from './helpers/expectThrowWithMessage'
import {initContracts} from './helpers/initContracts'
const RegistryKeeper = artifacts.require('./RegistryKeeper.sol')
const KeepRandomBeaconOperator = artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')

contract('RegistryKeeper', function(accounts) {

  let registryKeeper, stakingContract, operatorContract, anotherOperatorContract, serviceContract,
    governance = accounts[0],
    panicButton = accounts[1],
    operatorContractUpgrader = accounts[2]

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      KeepRandomBeaconOperator
    );

    stakingContract = contracts.stakingContract
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract
    await operatorContract.registerNewGroup("0x01")
    anotherOperatorContract = await KeepRandomBeaconOperator.new(serviceContract.address, stakingContract.address)
    await anotherOperatorContract.registerNewGroup("0x02")

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

  it("should be able to add or remove operator contracts from service contract", async() => {
    // Transfer ownership from governance to operatorContractUpgrader
    serviceContract.transferOwnership(operatorContractUpgrader, {from: governance})

    await expectThrowWithMessage(
      serviceContract.addOperatorContract(anotherOperatorContract.address, {from: governance}),
      "Ownable: caller is not the owner"
    )

    await serviceContract.addOperatorContract(anotherOperatorContract.address, {from: operatorContractUpgrader})
    assert.isTrue((await serviceContract.selectOperatorContract(0)) == operatorContract.address, "Unexpected operator contract address")
    assert.isTrue((await serviceContract.selectOperatorContract(1)) == anotherOperatorContract.address, "Unexpected operator contract address")

    await expectThrowWithMessage(
      serviceContract.removeOperatorContract(anotherOperatorContract.address, {from: governance}),
      "Ownable: caller is not the owner"
    )

    await serviceContract.removeOperatorContract(anotherOperatorContract.address, {from: operatorContractUpgrader})
    assert.isTrue((await serviceContract.selectOperatorContract(1)) == operatorContract.address, "Unexpected operator contract address")
  })
})
