const {initContracts} = require('../helpers/initContracts')
const assert = require('chai').assert
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {contract} = require("@openzeppelin/test-environment")

describe("KeepRandomBeaconOperator/Initialization", function() {
  let operatorContract;

  before(async() => {
    let contracts = await initContracts(
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorInitializationStub')
    ); 

    operatorContract = contracts.operatorContract;
  });

  beforeEach(async () => {
    await createSnapshot()
  });

  afterEach(async () => {
    await restoreSnapshot()
  });

  it("should set groups relay entry timeout", async () => {
    let relayEntryTimeout = await operatorContract.getGroupsRelayEntryTimeout()
    assert.equal(
      relayEntryTimeout.toNumber(), 
      384,
      "groups relay entry should have been set to (groupSize * resultPublicationBlockStep)"
    )
  })

  it("should set groups active time", async () => {
    let groupActiveTimeout = await operatorContract.getGroupsActiveTime()
    assert.equal(
      groupActiveTimeout.toNumber(), 
      80640,
      "group active time should have been set"
    )
  })
});
