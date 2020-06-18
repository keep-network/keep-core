const blsData = require("../helpers/data.js")
const {expectRevert, time} = require("@openzeppelin/test-helpers")
const initContracts = require('../helpers/initContracts')
const assert = require('chai').assert
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts, web3} = require("@openzeppelin/test-environment")

describe("KeepRandomBeaconOperator/RelayEntryTimeout", function() {
  let operatorContract, serviceContract, fee;
  const blocksForward = web3.utils.toBN(20);
  const requestCounter = 0;

  before(async() => {
    let contracts = await initContracts(
      contract.fromArtifact('KeepToken'),
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorStub')
    ); 

    registryContract = contracts.registry
    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;

    await registryContract.setServiceContractUpgrader(
      operatorContract.address, accounts[0], {from: accounts[0]}
    )
    await operatorContract.addServiceContract(accounts[0], {from: accounts[0]})  

    await operatorContract.registerNewGroup(blsData.groupPubKey);
    await operatorContract.setGroupMembers(blsData.groupPubKey, [accounts[0]]);

    fee = await serviceContract.entryFeeEstimate(0);
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should set fields for groups library", async () => {
    let relayEntryTimeout = await operatorContract.getRelayEntryTimeout()
    assert.equal(
      relayEntryTimeout.toNumber(), 
      384,
      "relay entry should have been set to (groupSize * resultPublicationBlockStep)"
    )

    let groupActiveTimeout = await operatorContract.getGroupActiveTime()
    assert.equal(
      groupActiveTimeout.toNumber(), 
      80640,
      "group active time should have been set"
    )
  })

  it("should not throw an error when entry is in progress and " +
     "block number > relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}
    );

    await time.advanceBlockTo(blocksForward.addn(await web3.eth.getBlockNumber()))

    await operatorContract.sign(
      requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}
    );

    assert.equal(
      (await operatorContract.getPastEvents())[0].event, 
      "RelayEntryRequested", 
      "RelayEntryRequested event should occur on operator contract"
    );
  });

  it("should throw an error when entry is in progress and " + 
     "block number <= relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}
    );

    await expectRevert(
      operatorContract.sign(requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}), 
      "Beacon is busy"
    );
  });

  it("should not throw an error when entry is not in progress and " + 
     "block number > relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}
      );

    assert.equal(
      (await operatorContract.getPastEvents())[0].event, 
      "RelayEntryRequested", 
      "RelayEntryRequested event should occur on operator contract."
    );
  });

  it("should not allow to submit relay entry after timeout", async () => {
    await operatorContract.sign(
      requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}
    );

    const relayEntryTimeout = await operatorContract.relayEntryTimeout()
    await time.advanceBlockTo(relayEntryTimeout.addn(await web3.eth.getBlockNumber()))

    await expectRevert(
      operatorContract.relayEntry(blsData.groupSignature), 
      "Entry timed out"
    );
  });
});
