const {contract, accounts} = require("@openzeppelin/test-environment")
const {expectRevert} = require("@openzeppelin/test-helpers")
const initContracts = require('../helpers/initContracts')
const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const blsData = require("../helpers/data.js")

describe('KeepRandomBeaconOperator/ManageServiceContracts', () => {
  let serviceContract
  let operatorContract
  let registry
  let serviceContract2 = accounts[1]
  let serviceContract3 = accounts[2]
  let serviceContractKeeper = accounts[3]
  let someoneElse = accounts[4]

  let groupProfitAndEntryVerificationFee;

  before(async () => {
    let contracts = await initContracts(
      contract.fromArtifact('KeepToken'),
      contract.fromArtifact('TokenStaking'),
      contract.fromArtifact('KeepRandomBeaconService'),
      contract.fromArtifact('KeepRandomBeaconServiceImplV1'),
      contract.fromArtifact('KeepRandomBeaconOperatorStub')
    )
            
    serviceContract = contracts.serviceContract
    operatorContract = contracts.operatorContract
    registry = contracts.registry
    
    await registry.setServiceContractKeeper(
        operatorContract.address, 
        serviceContractKeeper,
        {from: accounts[0]}
    )

    groupProfitFee = await operatorContract.groupProfitFee()
    entryVerificationFee = await operatorContract.entryVerificationFee()

    groupProfitAndEntryVerificationFee = groupProfitFee.add(
        entryVerificationFee
    )

    await operatorContract.registerNewGroup(blsData.groupPubKey)
  });

  beforeEach(async () => {
    await createSnapshot()
  });
  
  afterEach(async () => {
    await restoreSnapshot()
  });

  describe("addServiceContract", async () => {
    it("can be called by service contract keeper", async () => {
      await operatorContract.addServiceContract(
        serviceContract2, 
          {from: serviceContractKeeper}
      )
      // ok, no revert
    })

    it("cannot be called by non-keeper", async () => {
      await expectRevert(
        operatorContract.addServiceContract(
          serviceContract2, 
          {from: someoneElse}
        ),
        "Not authorized" 
      )
    })

    it("adds service contract to the list", async () => {
      await operatorContract.addServiceContract(
        serviceContract2, 
        {from: serviceContractKeeper}
      )

      await operatorContract.sign(
        1,
        blsData.previousEntry, 
        {
          value: groupProfitAndEntryVerificationFee, 
          from: serviceContract2
        }
      )
      // ok, no revert
    })
  })

  describe("removeServiceContract", async () => {
    it("can be called by service contract keeper", async () => {
      await operatorContract.removeServiceContract(
        serviceContract.address, 
        {from: serviceContractKeeper}
      )
      // ok, no revert
    })

    it("cannot be called by non-keeper", async () => {
      await expectRevert(
        operatorContract.removeServiceContract(
          serviceContract.address,
          {from: someoneElse}
        ),
        "Not authorized"
      )
    })

    it("removes service contract from the list", async () => {
      await operatorContract.addServiceContract(
        serviceContract2,
        {from: serviceContractKeeper}
      )
      await operatorContract.addServiceContract(
        serviceContract3,
        {from: serviceContractKeeper}
      )

      await operatorContract.removeServiceContract(
        serviceContract2,
        {from: serviceContractKeeper}
      )

      await expectRevert(
        operatorContract.sign(
          1, 
          blsData.previousEntry, 
          {
            value: groupProfitAndEntryVerificationFee, 
            from: serviceContract2
          }
        ),
        "Caller is not a service contract"
      )

      await operatorContract.sign(
        1, 
        blsData.previousEntry, 
        {
          value: groupProfitAndEntryVerificationFee, 
          from: serviceContract3
        }
      )
      // ok, no revert - the second service contract is still there
    })
  })
})