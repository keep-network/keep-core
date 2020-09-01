const {createSnapshot, restoreSnapshot} = require("../helpers/snapshot.js")
const {BN, constants, expectEvent, expectRevert, time} = require("@openzeppelin/test-helpers")
const {contract, accounts} = require("@openzeppelin/test-environment")
const assert = require('chai').assert

const ServiceContractProxy = contract.fromArtifact('KeepRandomBeaconService')
const ServiceContractImplV1 = contract.fromArtifact('KeepRandomBeaconServiceImplV1')
const ServiceContractImplV2 = contract.fromArtifact('KeepRandomBeaconServiceUpgradeExample')

const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect

describe('KeepRandomBeaconService/Upgrade', function() {

  let proxy
  let implementationV1
  let implementationV2
  
  let initializeCallData

  const admin = accounts[1]
  const nonAdmin = accounts[2]
  const newAdmin = accounts[3]

  before(async () => {
    implementationV1 = await ServiceContractImplV1.new({from: admin})
    implementationV2 = await ServiceContractImplV2.new({from: admin})
    
    initializeCallData = implementationV1.contract.methods.initialize(
      100, '0x0000000000000000000000000000000000000001'
    ).encodeABI()

    proxy = await ServiceContractProxy.new(
      implementationV1.address, initializeCallData, {from: admin}
    )
  })

  beforeEach(async () => {
    await createSnapshot()
  })

  afterEach(async () => {
    await restoreSnapshot()
  })

  describe("constructor", async () => {
    it("sets admin", async () => {
      assert.equal(
        await proxy.admin(), 
        admin, 
        "Unexpected admin"
      )
    })

    it("sets upgrade time delay to one day", async () => {
      assert.equal(
        (await proxy.upgradeTimeDelay()).toNumber(),
        86400, // 1 day
        "Upgrade time delay should be one day"
      )
    })

    it("initializes implementation", async () => {
      assert.isTrue(
        await implementationV1.initialized(),
        "Implementation contract should be initialized"
      )
    })

    it("sets implementation", async () => {
      assert.equal(
        await proxy.implementation(),
        implementationV1.address,
        "Unexpected implementation contract address"
      )
    })
  })

  describe("upgradeTo", async () => {
    it("sets timestamp", async () => {
      await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData,
        {from: admin}
      )

      const expectedTimestamp = await time.latest()

      expect(await proxy.upgradeInitiatedTimestamp()).to.eq.BN(
        expectedTimestamp
      )
    })

    it("sets new implementation", async () => {
      await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData,
        {from: admin}
      )
  
      assert.equal(
        await proxy.newImplementation(),
        implementationV2.address,
        "Unexpected new implementation contract address"
      )
      assert.equal(
        await proxy.implementation(),
        implementationV1.address,
        "Unexpected implementation contract address"
      )
    })
  
    it("sets initialization call data", async () => {
      await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData,
        {from: admin}
      )
  
      assert.equal(
        await proxy.initializationData(implementationV2.address),
        initializeCallData,
        "Unexpected initialization call data"
      )
    })
  
    it("supports empty initialization call data", async () => {
      await proxy.upgradeTo(implementationV2.address, [], {from: admin})
      assert.notExists(await proxy.initializationData.call(implementationV2.address));
    })
  
    it("emits an event", async () => {
      const receipt = await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData,
        {from: admin}
      )
  
      const expectedTimestamp = await time.latest()
      expectEvent(receipt, "UpgradeStarted", {
        implementation: implementationV2.address,
        timestamp: expectedTimestamp
      })
    })
  
    it("allows implementation overwrite", async () => {
      const address3 = '0x4566716c07617c5854fe7dA9aE5a1219B19CCd27'
      await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData,
        {from: admin}
      )
      await proxy.upgradeTo(
        address3, 
        initializeCallData,
        {from: admin}
      )
  
      assert.equal(
        await proxy.newImplementation(),
        address3,
        "Unexpected new implementation contract address"
      )
    })
  
    it("allows implementation data overwrite", async () => {
      const initializeCallData2 = '0x123456'
      await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData,
        {from: admin}
      )
      await proxy.upgradeTo(
        implementationV2.address, 
        initializeCallData2,
        {from: admin}
      )
  
      assert.equal(
        await proxy.initializationData.call(implementationV2.address), 
        initializeCallData2,
        "unexpected initialization call data"
      )
    })
  
    it("reverts on zero address", async () => {
      await expectRevert(
        proxy.upgradeTo(
          constants.ZERO_ADDRESS, 
          initializeCallData, 
          {from: admin}
        ),
        "Implementation address can't be zero."
      )
    })
  
    it("reverts on the same address", async () => {
      await expectRevert(
        proxy.upgradeTo(
          implementationV1.address,
          initializeCallData,
          {from: admin}
        ), 
        "Implementation address must be different from the current one."
      )
    })
  
    it("reverts when called by non-admin", async () => {
      await expectRevert(
        proxy.upgradeTo(
          implementationV2.address,
          initializeCallData,
          {from: nonAdmin}
        ),
        "Caller is not the admin."
      )
    })
  })

  describe("completeUpgrade", async () => {
    it("reverts for non-initiated upgrade", async () => {
      await expectRevert(
        proxy.completeUpgrade({from: admin}),
        "Upgrade not initiated"
      )
    })

    it("reverts for non-elapsed timer", async () => {
      await proxy.upgradeTo(
        implementationV2.address,
        initializeCallData,
        {from: admin}
      )

      await time.increase((await proxy.upgradeTimeDelay()).subn(2))

      await expectRevert(
        proxy.completeUpgrade({ from: admin }), 
        "Timer not elapsed"
      )
    })

    it("clears timestamp", async () => {
      await proxy.upgradeTo(
        implementationV2.address,
        initializeCallData,
        {from: admin}
      )

      await time.increase(await proxy.upgradeTimeDelay())

      await proxy.completeUpgrade({from: admin})

      expect(await proxy.upgradeInitiatedTimestamp()).to.eq.BN(0)
    })

    it("sets implementation", async () => {
      await proxy.upgradeTo(
        implementationV2.address,
        initializeCallData,
        {from: admin}
      )

      await time.increase(await proxy.upgradeTimeDelay())

      await proxy.completeUpgrade({from: admin})

      assert.equal(
        await proxy.implementation(),
        implementationV2.address,
        "Unexpected new implementation address"
      )
    })

    it("emits an event", async () => {
      await proxy.upgradeTo(
        implementationV2.address,
        initializeCallData,
        {from: admin}
      )

      await time.increase(await proxy.upgradeTimeDelay())

      const receipt = await proxy.completeUpgrade({from: admin})

      await expectEvent(receipt, "UpgradeCompleted", {
        implementation: implementationV2.address
      })
    })

    it("supports empty initialization call data", async () => {
      const address3 = '0x4566716c07617c5854fe7dA9aE5a1219B19CCd27'
      await proxy.upgradeTo(address3, [], {from: admin})
      await time.increase(await proxy.upgradeTimeDelay());

      await proxy.completeUpgrade({from: admin});
    });

    it("reverts when called by non-admin", async () => {
      await expectRevert(
        proxy.completeUpgrade({from: nonAdmin}),
        "Caller is not the admin."
      )
    })

    it("reverts when initialization fails", async () => {
      const failingData = implementationV1.contract.methods.initialize(
        100, constants.ZERO_ADDRESS
      ).encodeABI()

      await proxy.upgradeTo(
        implementationV2.address,
        failingData,
        {from: admin}
      )

      await time.increase(await proxy.upgradeTimeDelay())

      await expectRevert(
        proxy.completeUpgrade({from: admin}),
        "Incorrect registry address"
      )
    })

    it("finalizes upgrade procedure", async () => {
      await proxy.upgradeTo(
        implementationV2.address,
        initializeCallData,
        {from: admin}
      )

      await time.increase(await proxy.upgradeTimeDelay())

      await proxy.completeUpgrade({from: admin})

      let v2 = await ServiceContractImplV2.at(proxy.address)
      assert.equal(
        await v2.getNewVar(),
        1234,
        "Should be able to get new data from upgraded contract"
      )
    })
  })

  describe("updateAdmin", async () => {
    it("sets new admin when called by admin", async () => {
      await proxy.updateAdmin(newAdmin, { from: admin })

      assert.equal(await proxy.admin(), newAdmin, "Unexpected admin")
    })

    it("reverts when called by non-admin", async () => {
      await expectRevert(
        proxy.updateAdmin(newAdmin, { from: nonAdmin }),
        "Caller is not the admin"
      )
    })

    it("reverts when called by admin after role transfer", async () => {
      await proxy.updateAdmin(newAdmin, { from: admin })

      await expectRevert(
        proxy.updateAdmin(nonAdmin, { from: admin }),
        "Caller is not the admin"
      )
    })
  })
});
