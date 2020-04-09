const blsData = require("../helpers/data.js")
const {expectRevert} = require("@openzeppelin/test-helpers")
const initContracts = require('../helpers/initContracts')
const assert = require('chai').assert
const mineBlocks = require("../helpers/mineBlocks")
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract, accounts} = require("@openzeppelin/test-environment")

describe("KeepRandomBeaconOperator/RelayEntryTimeout", function() {
  let operatorContract, serviceContract, fee;
  const blocksForward = 20;
  const requestCounter = 0;

  before(async() => {
    let contracts = await initContracts(
      contract.fromArtifact('KeepToken'),
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorStub')
    ); 

    operatorContract = contracts.operatorContract;
    serviceContract = contracts.serviceContract;

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

  it("should not throw an error when entry is in progress and " +
     "block number > relay entry timeout", async () => {
    await operatorContract.sign(
      requestCounter, blsData.previousEntry, {value: fee, from: accounts[0]}
    );

    mineBlocks(blocksForward)

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

    mineBlocks(await operatorContract.relayEntryTimeout());

    await expectRevert(
      operatorContract.relayEntry(blsData.groupSignature), 
      "Entry timed out"
    );
  });
});
