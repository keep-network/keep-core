const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { expectRevert, expectEvent } = require("@openzeppelin/test-helpers")
const { time } = require("@openzeppelin/test-helpers")
const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")

const GasPriceOracle = contract.fromArtifact("GasPriceOracle");
const GasPriceOracleConsumerStub = contract.fromArtifact("GasPriceOracleConsumerStub");

const BN = web3.utils.BN
const chai = require('chai');
chai.use(require('bn-chai')(BN))
const expect = chai.expect
const assert = chai.assert

describe("GasPriceOracle", () => {

    const owner = accounts[1]
    const thirdParty = accounts[2]

    let oracle

    before(async () => {
        oracle = await GasPriceOracle.new({from: owner})
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("when updating gas price", async () => {
        it("does not let third party to begin the process", async () => {
            await expectRevert(
                oracle.beginGasPriceUpdate(123, {from: thirdParty}),
                "caller is not the owner"
            )
        })

        it("lets the owner to begin the process", async () => {
            await oracle.beginGasPriceUpdate(123, {from: owner})
            // ok, no revert
        })

        it("does not allow to finalize change without initiating it first", async () => {
            await expectRevert(
                oracle.finalizeGasPriceUpdate(),
                "Change not initiated"
            )
        })

        it("does not allow to finalize change before governance delay", async () => {
            await oracle.beginGasPriceUpdate(123, {from: owner})
            await time.increase(time.duration.minutes(59))
            await expectRevert(
                oracle.finalizeGasPriceUpdate(),
                "Governance delay has not elapsed"
            )
        })

        it("updates value when finalizing the change", async () => {
            const newValue = 9129111
            await oracle.beginGasPriceUpdate(newValue, {from: owner})
            await time.increase(time.duration.minutes(61))
            await oracle.finalizeGasPriceUpdate()

            expect(await oracle.gasPrice()).to.eq.BN(newValue)
        })

        it("does not allow to finalize the change twice", async () => {
            const newValue = 1111
            await oracle.beginGasPriceUpdate(newValue, {from: owner})
            await time.increase(time.duration.minutes(61))
            await oracle.finalizeGasPriceUpdate()   
            await expectRevert(
                oracle.finalizeGasPriceUpdate(),
                "Change not initiated"
            )
        })

        it("emits an event when finalizing the change", async () => {
            const newValue = web3.utils.toBN(55555)
            await oracle.beginGasPriceUpdate(newValue, {from: owner})
            await time.increase(time.duration.minutes(61))
            const receipt = await oracle.finalizeGasPriceUpdate()
            await expectEvent(receipt, "GasPriceUpdated", {
                newValue: newValue
            })
        })
    
        it("notifies consumer contracts when finalizing the change", async () => {
            const consumer = await GasPriceOracleConsumerStub.new(oracle.address)
            await oracle.addConsumerContract(consumer.address, {from: owner});

            const newValue = 545666
            await oracle.beginGasPriceUpdate(newValue, {from: owner})
            await time.increase(time.duration.minutes(61))
            await oracle.finalizeGasPriceUpdate()   

            expect(await consumer.gasPrice()).to.eq.BN(newValue)
        })

        it("does not notify removed consumer contracts when finalizing the change", async () => {
            const consumer1 = await GasPriceOracleConsumerStub.new(oracle.address)
            const consumer2 = await GasPriceOracleConsumerStub.new(oracle.address)
            const consumer3 = await GasPriceOracleConsumerStub.new(oracle.address)

            await oracle.addConsumerContract(consumer1.address, {from: owner})
            await oracle.addConsumerContract(consumer2.address, {from: owner})
            await oracle.addConsumerContract(consumer3.address, {from: owner})

            await oracle.removeConsumerContract(1, {from: owner})

            const newValue = 156444
            await oracle.beginGasPriceUpdate(newValue, {from: owner})
            await time.increase(time.duration.minutes(61))
            await oracle.finalizeGasPriceUpdate()   

            expect(await consumer1.gasPrice()).to.eq.BN(newValue)
            expect(await consumer2.gasPrice()).to.eq.BN(0)
            expect(await consumer3.gasPrice()).to.eq.BN(newValue)
        })

        it("lets to overwrite the pending update", async () => {
            await oracle.beginGasPriceUpdate(55555, {from: owner})
            await time.increase(time.duration.minutes(59))

            const overwritten = 66666
            await oracle.beginGasPriceUpdate(overwritten, {from: owner})
            
            expect(await oracle.newGasPrice()).to.eq.BN(overwritten)
            expect(await oracle.gasPriceChangeInitiated()).to.eq.BN(await time.latest())     
            
            await time.increase(time.duration.minutes(61))
            await oracle.finalizeGasPriceUpdate()   

            expect(await oracle.gasPrice()).to.eq.BN(overwritten)
        })
    })

    describe("when managing consumer contracts", async () => {
        it("does not allow third party to add new consumer contract", async () => {
            const consumer = await GasPriceOracleConsumerStub.new(oracle.address)
            await expectRevert(
                oracle.addConsumerContract(consumer.address, {from: thirdParty}),
                "Ownable: caller is not the owner."
            )
        })

        it("allows the owner to add new consumer contract", async () => {
            const consumer = await GasPriceOracleConsumerStub.new(oracle.address)
            await oracle.addConsumerContract(consumer.address, {from: owner});
            // ok, no revert
        })

        it("does not allow third party to remove consumer contract", async () => {
            const consumer = await GasPriceOracleConsumerStub.new(oracle.address)
            await oracle.addConsumerContract(consumer.address, {from: owner});
            await expectRevert(
                oracle.removeConsumerContract(0, {from: thirdParty}),
                "Ownable: caller is not the owner."
            ) 
        })

        it("allows the owner to remove consumer contract", async () => {
            const consumer = await GasPriceOracleConsumerStub.new(oracle.address)
            await oracle.addConsumerContract(consumer.address, {from: owner});
            await oracle.removeConsumerContract(0, {from: owner})
            // ok, no revert
        })

        it("allows to selectively remove consumer contracts", async () => {
            const consumer1 = await GasPriceOracleConsumerStub.new(oracle.address)
            const consumer2 = await GasPriceOracleConsumerStub.new(oracle.address)
            const consumer3 = await GasPriceOracleConsumerStub.new(oracle.address)

            await oracle.addConsumerContract(consumer1.address, {from: owner})
            await oracle.addConsumerContract(consumer2.address, {from: owner})
            await oracle.addConsumerContract(consumer3.address, {from: owner})

            await oracle.removeConsumerContract(1, {from: owner})

            const consumers = await oracle.getConsumerContracts()
            assert.equal(consumers.length, 2)
            assert.equal(consumers[0], consumer1.address)
            assert.equal(consumers[1], consumer3.address)
        })
    })
})