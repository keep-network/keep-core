const {createSnapshot, restoreSnapshot} = require("./helpers/snapshot.js")
const expectThrowWithMessage = require('./helpers/expectThrowWithMessage.js')
const {accounts, contract} = require("@openzeppelin/test-environment")
var assert = require('chai').assert

const Registry = contract.fromArtifact('RegistryStub');

describe('Registry', () => {
    
    const owner = accounts[0]
    const governance = accounts[1]
    const panicButton = accounts[2]
    const registryKeeper = accounts[3]
    const operatorContractUpgrader = accounts[4]

    const someoneElse = "0x524f2E0176350d950fA630D9A5a59A0a190DAf48"

    const serviceContract1 = "0xF2D3Af2495E286C7820643B963FB9D34418c871d"
    const serviceContract2 = "0x65EA55c1f10491038425725dC00dFFEAb2A1e28A"
    const operatorContract1 = "0x4566716c07617c5854fe7dA9aE5a1219B19CCd27"
    const operatorContract2 = "0x7020A5556Ba1ce5f92c81063a13d33512cf1305c"

    let registry

    before(async () => {
        registry = await Registry.new({from: owner})
        await registry.setGovernance(governance, {from: owner})
        await registry.setPanicButton(panicButton, {from: governance})
        await registry.setRegistryKeeper(registryKeeper, {from: governance})
        await registry.setOperatorContractUpgrader(
            serviceContract1, 
            operatorContractUpgrader,
            {from: governance}
        )
    })

    beforeEach(async () => {
        await createSnapshot()
    })
    
    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("setGovernance", async () => {
        it("can be called by governance", async () => {
            await registry.setGovernance(someoneElse, {from: governance})
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectThrowWithMessage(
                registry.setGovernance(someoneElse, {from: owner}),
                "Not authorized"
            )
        })

        it("updates governance", async () => {
            await registry.setGovernance(someoneElse, {from: governance})
            assert.equal(
                await registry.getGovernance(),
                someoneElse,
                "Unexpected governance"
            )
        })
    })

    describe("setRegistryKeeper", async () => {
        it("can be called by governance", async () => {
            await registry.setRegistryKeeper(someoneElse, {from: governance})
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectThrowWithMessage(
                registry.setRegistryKeeper(someoneElse, {from: owner}),
                "Not authorized"
            )
        })

        it("updates registry keeper", async () => {
            await registry.setRegistryKeeper(someoneElse, {from: governance})
            assert.equal(
                await registry.getRegistryKeeper(),
                someoneElse,
                "Unexpected registry keeper"
            )
        })
    })

    describe("setPanicButton", async () => {
        it("can be called by governance", async () => {
            await registry.setPanicButton(someoneElse, {from: governance})
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectThrowWithMessage(
                registry.setPanicButton(someoneElse, {from: owner}),
                "Not authorized"
            )
        })

        it("updates panic button", async () => {
            await registry.setPanicButton(someoneElse, {from: governance})
            assert.equal(
                await registry.getPanicButton(),
                someoneElse,
                "Unexpected registry keeper"
            )
        })
    })

    describe("setOperatorContractUpgrader", async () => {
        it("can be called by governance", async () => {
            await registry.setOperatorContractUpgrader(
                serviceContract1, 
                someoneElse, 
                {from: governance}
            )
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectThrowWithMessage(
                registry.setOperatorContractUpgrader(
                    serviceContract1,
                    someoneElse, 
                    {from: owner}
                ),
                "Not authorized"
            )
        })

        it("updates operator contract upgrader", async () => {
            await registry.setOperatorContractUpgrader(
                serviceContract1,
                someoneElse,
                {from: governance}
            )

            await registry.setOperatorContractUpgrader(
                serviceContract2,
                operatorContractUpgrader,
                {from: governance}
            )

            assert.equal(
                await registry.operatorContractUpgraderFor(serviceContract1),
                someoneElse,
                "Unexpected operator contract upgrader"
            )

            assert.equal(
                await registry.operatorContractUpgraderFor(serviceContract2),
                operatorContractUpgrader,
                "Unexpected operator contract upgrader"
            )
        })
    })

    describe("approveOperatorContract", async () => {
        it("can be called by registry keeper", async () => {
            await registry.approveOperatorContract(
                operatorContract1, 
                {from: registryKeeper}
            )
            // ok, no revert
        })

        it("can not be called by non-registry-keeper", async () => {
            await expectThrowWithMessage(
                registry.approveOperatorContract(
                    operatorContract1,
                    {from: owner}
                ),
                "Not authorized"
            )
        })

        it("approves operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                {from: registryKeeper}
            )

            assert.isTrue(
                await registry.isApprovedOperatorContract(operatorContract1),
                "operator contract should be approved"
            )
            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract2),
                "operator contract should not be approved"
            )
        })

        it("cannot be called for already approved contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                {from: registryKeeper}
            )

            await expectThrowWithMessage(
                registry.approveOperatorContract(
                    operatorContract1,
                    {from: registryKeeper}

                ),
                "Only new operator contracts can be approved"
            )
        })

        it("cannot be called for disabled contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                {from: registryKeeper}
            )
            await registry.disableOperatorContract(
                operatorContract1,
                {from: panicButton}
            )

            await expectThrowWithMessage(
                registry.approveOperatorContract(
                    operatorContract1,
                    {from: registryKeeper}

                ),
                "Only new operator contracts can be approved"
            )
        })
    })

    describe("disableOperatorContract", async () => {
        beforeEach(async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                {from: registryKeeper}
            )
        })

        it("can be called by panic button", async () => {
            await registry.disableOperatorContract(
                operatorContract1, 
                {from: panicButton}
            )
            // ok, no revert
        })

        it("can not be called by non-registry-keeper", async () => {
            await expectThrowWithMessage(
                registry.disableOperatorContract(
                    operatorContract1,
                    {from: owner}
                ),
                "Not authorized"
            )
        })

        it("disables operator contract", async () => {
            await registry.disableOperatorContract(
                operatorContract1, 
                {from: panicButton}
            )

            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract1),
                "operator contract should not be approved"
            )
        })

        it("cannot be called for already disabled contract", async () => {
            await registry.disableOperatorContract(
                operatorContract1, 
                {from: panicButton}
            )

            await expectThrowWithMessage(
                registry.disableOperatorContract(
                    operatorContract1, 
                    {from: panicButton}
                ),
                "Only approved operator contracts can be disabled"
            )
        })

        it("cannot be called for new operator contract", async () => {
            await expectThrowWithMessage(
                registry.disableOperatorContract(
                    operatorContract2, 
                    {from: panicButton}
                ),
                "Only approved operator contracts can be disabled"
            )
        })
    })

    describe("isNewOperatorContract", async () => {
        it("returns true for new operator contracts", async () => {
            assert.isTrue(
                await registry.isNewOperatorContract(operatorContract1),
                "Expected true for new operator contract"
            )
        })

        it("returns false for approved operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1, 
                {from: registryKeeper}
            )

            assert.isFalse(
                await registry.isNewOperatorContract(operatorContract1),
                "Expected false for approved operator contract"
            )
        })

        it("returns false for disabled operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1, 
                {from: registryKeeper}
            ) 
            await registry.disableOperatorContract(
                operatorContract1,
                {from: panicButton}
            )

            assert.isFalse(
                await registry.isNewOperatorContract(operatorContract1),
                "Expected false for disabled operator contract"
            )
        })
    })

    describe("isApprovedOperatorContract", async () => {
        it("returns false for new operator contracts", async () => {
            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract1),
                "Expected false for new operator contract"
            )
        })

        it("returns true for approved operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1, 
                {from: registryKeeper}
            )

            assert.isTrue(
                await registry.isApprovedOperatorContract(operatorContract1),
                "Expected true for approved operator contract"
            )
        })

        it("returns false for disabled operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1, 
                {from: registryKeeper}
            ) 
            await registry.disableOperatorContract(
                operatorContract1,
                {from: panicButton}
            )

            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract1),
                "Expected false for disabled operator contract"
            )
        })
    })
})