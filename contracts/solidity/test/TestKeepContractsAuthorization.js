import {createSnapshot, restoreSnapshot} from './helpers/snapshot'
import expectThrowWithMessage from './helpers/expectThrowWithMessage'
import {initContracts} from './helpers/initContracts'
const KeepRandomBeaconOperator = artifacts.require('./stubs/KeepRandomBeaconOperatorStub.sol')

contract('Registry', function(accounts) {

  let registry, stakingContract, operatorContract, anotherOperatorContract, serviceContract,
    governance = accounts[0],
    panicButton = accounts[1],
    operatorContractUpgrader = accounts[2],
    registryKeeper = accounts[3]

  before(async () => {

    let contracts = await initContracts(
      artifacts.require('./KeepToken.sol'),
      artifacts.require('./TokenStaking.sol'),
      artifacts.require('./KeepRandomBeaconService.sol'),
      artifacts.require('./KeepRandomBeaconServiceImplV1.sol'),
      KeepRandomBeaconOperator
    );

    registry = contracts.registry
    stakingContract = contracts.stakingContract
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract
    await operatorContract.registerNewGroup("0x01")
    anotherOperatorContract = await KeepRandomBeaconOperator.new(serviceContract.address, stakingContract.address)
    await anotherOperatorContract.registerNewGroup("0x02")

    await registry.setRegistryKeeper(registryKeeper)
    await registry.setPanicButton(panicButton)
    await registry.setOperatorContractUpgrader(serviceContract.address, operatorContractUpgrader)
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  it("should be able to approve or disable operator contract", async() => {
    assert.isTrue((await registry.operatorContracts(anotherOperatorContract.address)).eqn(0), "Unexpected status of operator contract")

    await expectThrowWithMessage(
      registry.approveOperatorContract(anotherOperatorContract.address, {from: operatorContractUpgrader}),
      "Not authorized"
    );

    await registry.approveOperatorContract(anotherOperatorContract.address, {from: registryKeeper})
    assert.isTrue((await registry.operatorContracts(anotherOperatorContract.address)).eqn(1), "Unexpected status of operator contract")

    await expectThrowWithMessage(
      registry.disableOperatorContract(anotherOperatorContract.address, {from: governance}),
      "Not authorized"
    );

    await registry.disableOperatorContract(operatorContract.address, {from: panicButton})
    assert.isTrue((await registry.operatorContracts(operatorContract.address)).eqn(2), "Unexpected status of operator contract")
  })

  it("should be able to add or remove operator contracts from service contract", async() => {
    await registry.approveOperatorContract(anotherOperatorContract.address, {from: registryKeeper})

    await expectThrowWithMessage(
      serviceContract.addOperatorContract(anotherOperatorContract.address, {from: governance}),
      "Caller is not operator contract upgrader"
    )

    await serviceContract.addOperatorContract(anotherOperatorContract.address, {from: operatorContractUpgrader})
    assert.isTrue((await serviceContract.selectOperatorContract(0)) == operatorContract.address, "Unexpected operator contract address")
    assert.isTrue((await serviceContract.selectOperatorContract(1)) == anotherOperatorContract.address, "Unexpected operator contract address")

    await expectThrowWithMessage(
      serviceContract.removeOperatorContract(anotherOperatorContract.address, {from: governance}),
      "Caller is not operator contract upgrader"
    )

    await serviceContract.removeOperatorContract(anotherOperatorContract.address, {from: operatorContractUpgrader})
    assert.isTrue((await serviceContract.selectOperatorContract(1)) == operatorContract.address, "Unexpected operator contract address")
  })

  it("should be able to disable operator contract via panic button", async() => {
    await registry.approveOperatorContract(anotherOperatorContract.address, {from: registryKeeper})
    await serviceContract.addOperatorContract(anotherOperatorContract.address, {from: operatorContractUpgrader})

    assert.isTrue((await serviceContract.selectOperatorContract(1)) == anotherOperatorContract.address, "Unexpected operator contract address")
    await registry.disableOperatorContract(anotherOperatorContract.address, {from: panicButton})
    assert.isTrue((await serviceContract.selectOperatorContract(1)) == operatorContract.address, "Unexpected operator contract address")
  })
})
