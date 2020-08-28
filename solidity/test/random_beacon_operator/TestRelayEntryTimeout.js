
const {expectRevert, expectEvent, time} = require("@openzeppelin/test-helpers")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts, web3} = require("@openzeppelin/test-environment")
const blsData = require("../helpers/data.js")
const stakeDelegate = require('../helpers/stakeDelegate')
const {initContracts} = require('../helpers/initContracts')

describe("KeepRandomBeaconOperator/RelayEntryTimeout", function() {
  const deployer = accounts[0],
    serviceContractUpgrader = accounts[1]
    serviceContract = accounts[2],
    operator1 = accounts[3],
    operator2 = accounts[4],
    operator3 = accounts[5],
    thirdParty = accounts[6]

  let operatorContract, entryFee;

  before(async() => {
    let contracts = await initContracts(
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorStub')
    ); 

    operatorContract = contracts.operatorContract;
  
    //
    // register 'serviceContract' account as a new service contract so that
    // we can hit the operator contract from this account for tests' simplicity
    //
    await contracts.registry.setServiceContractUpgrader(
      operatorContract.address,
      serviceContractUpgrader,
      {from: deployer}
    )
    await operatorContract.addServiceContract(
      serviceContract,
      {from: serviceContractUpgrader}
    )  

    //
    // stake 3 operators, authorize operator contract for all of them,
    // and wait for the stake initialization period to complete
    //
    let token = contracts.token
    let tokenStaking = contracts.stakingContract
    const stake = web3.utils.toBN("500000000000000000000000") // 500 000 KEEP
    await stakeDelegate(tokenStaking, token, deployer, operator1, operator1, operator1, stake)
    await stakeDelegate(tokenStaking, token, deployer, operator2, operator2, operator2, stake)
    await stakeDelegate(tokenStaking, token, deployer, operator3, operator3, operator3, stake)
    await tokenStaking.authorizeOperatorContract(operator1, operatorContract.address, { from: operator1 })
    await tokenStaking.authorizeOperatorContract(operator2, operatorContract.address, { from: operator2 })
    await tokenStaking.authorizeOperatorContract(operator3, operatorContract.address, { from: operator3 })
    await time.increase((await tokenStaking.initializationPeriod()).addn(1))
    
    //
    // register two groups with operators staked in the previous step
    //
    await operatorContract.registerNewGroup("0x111")
    await operatorContract.setGroupMembers("0x111", [operator1, operator2, operator3])
    await operatorContract.registerNewGroup(blsData.groupPubKey)
    await operatorContract.setGroupMembers(blsData.groupPubKey, [operator3, operator2, operator1])

    entryFee = await contracts.serviceContract.entryFeeEstimate(0);
  });

  async function requestRelayEntry() {
    return operatorContract.sign(
      0,
      blsData.previousEntry,
      {value: entryFee, from: serviceContract}
    )
  }

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("request for a new relay entry", async () => {
    it("should be accepted when no other request is in progress", async () => {
      const receipt = await requestRelayEntry()
      await expectEvent(receipt, 'RelayEntryRequested') 
    })

    it("should be rejected when another request is in progress", async () => {
      await requestRelayEntry()
      await expectRevert(requestRelayEntry(), "Beacon is busy");
    })

    it("should be rejected when another request is in progress even if it timed out", async () => {
      const timeout = await operatorContract.relayEntryTimeout()      
      await requestRelayEntry()
      await time.advanceBlockTo((await time.latestBlock()).add(timeout))    
      await expectRevert(requestRelayEntry(), "Beacon is busy");
    })

    it("should be retried when another request timed out and it's been reported", async () => {
      const timeout = await operatorContract.relayEntryTimeout()      
      await requestRelayEntry()
      await time.advanceBlockTo((await time.latestBlock()).add(timeout))
      const receipt = await operatorContract.reportRelayEntryTimeout({from: thirdParty})
      await expectEvent(receipt, 'RelayEntryRequested') 
    })
  })

  describe("relay entry submission", async () => {
    it("should be rejected after timeout", async() => {
      const timeout = await operatorContract.relayEntryTimeout()      
      await requestRelayEntry()
      await time.advanceBlockTo((await time.latestBlock()).add(timeout))
      await expectRevert(
        operatorContract.relayEntry(blsData.groupSignature), 
        "Entry timed out"
      );
    })

    it("should be accepted when the previous request timed out and it's been reported", async () => {
      const timeout = await operatorContract.relayEntryTimeout()      
      await requestRelayEntry()
      await time.advanceBlockTo((await time.latestBlock()).add(timeout))
      await operatorContract.reportRelayEntryTimeout({from: thirdParty})
      // 0x111 group gets terminated and bls.groupPubKey group is now asked
      // to provide the signature

      await operatorContract.relayEntry(blsData.groupSignature)
      // ok, no revert
    })
  })
});
