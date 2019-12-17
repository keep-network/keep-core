import {createSnapshot, restoreSnapshot} from "./helpers/snapshot"
import expectThrowWithMessage from './helpers/expectThrowWithMessage';
const RegistryKeeper = artifacts.require('./RegistryKeeper.sol')
const KeepToken = artifacts.require('./KeepToken.sol')
const TokenStaking = artifacts.require('./TokenStaking.sol')
const KeepRandomBeaconOperator = artifacts.require('./KeepRandomBeaconOperator.sol')
const KeepRandomBeaconService = artifacts.require('./KeepRandomBeaconService.sol')
const KeepRandomBeaconServiceImplV1 = artifacts.require('./KeepRandomBeaconServiceImplV1.sol')

contract('RegistryKeeper', function(accounts) {

  let registryKeeper, token, stakingContract, operatorContract,
    serviceContractImplV1, serviceContractProxy, serviceContract,
    governance = accounts[0],
    panicButton = accounts[1],
    operatorContractUpgrader = accounts[2]

  before(async () => {
    registryKeeper = await RegistryKeeper.new(panicButton)
    token = await KeepToken.new();
    stakingContract = await TokenStaking.new(token.address, 0)
    serviceContractImplV1 = await KeepRandomBeaconServiceImplV1.new({from: operatorContractUpgrader})
    serviceContractProxy = await KeepRandomBeaconService.new(serviceContractImplV1.address, {from: operatorContractUpgrader})
    serviceContract = await KeepRandomBeaconServiceImplV1.at(serviceContractProxy.address)
    operatorContract = await KeepRandomBeaconOperator.new(serviceContract.address, stakingContract.address)
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
