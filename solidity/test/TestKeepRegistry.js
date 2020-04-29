const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")
const { accounts, contract } = require("@openzeppelin/test-environment")
const { expectRevert, expectEvent } = require("@openzeppelin/test-helpers")
var assert = require('chai').assert

const KeepRegistry = contract.fromArtifact('KeepRegistry');

describe('KeepRegistry', () => {

    const owner = accounts[0]
    const governance = accounts[1]
    const defaultPanicButton = accounts[2]
    const registryKeeper = accounts[3]
    const operatorContractUpgrader = accounts[4]
    const individualContractPanicButton = accounts[5]
    const serviceContractUpgrader = accounts[6]

    const someoneElse = "0x524f2E0176350d950fA630D9A5a59A0a190DAf48"

    const serviceContract1 = "0xF2D3Af2495E286C7820643B963FB9D34418c871d"
    const serviceContract2 = "0x65EA55c1f10491038425725dC00dFFEAb2A1e28A"
    const operatorContract1 = "0x4566716c07617c5854fe7dA9aE5a1219B19CCd27"
    const operatorContract2 = "0x7020A5556Ba1ce5f92c81063a13d33512cf1305c"

    let registry

    before(async () => {
        registry = await KeepRegistry.new({ from: owner })
        await registry.setGovernance(governance, { from: owner })
        await registry.setDefaultPanicButton(defaultPanicButton, { from: governance })
        await registry.setRegistryKeeper(registryKeeper, { from: governance })
        await registry.setOperatorContractUpgrader(
            serviceContract1,
            operatorContractUpgrader,
            { from: governance }
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
            await registry.setGovernance(someoneElse, { from: governance })
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectRevert(
                registry.setGovernance(someoneElse, { from: owner }),
                "Not authorized"
            )
        })

        it("updates governance", async () => {
            await registry.setGovernance(someoneElse, { from: governance })
            assert.equal(
                await registry.governance(),
                someoneElse,
                "Unexpected governance"
            )
        })

        it("emits an event", async () => {
            const receipt = await registry.setGovernance(
                someoneElse, { from: governance }
            )
            expectEvent(receipt, "GovernanceUpdated", {
                governance: someoneElse
            })
        })
    })

    describe("setRegistryKeeper", async () => {
        it("can be called by governance", async () => {
            await registry.setRegistryKeeper(someoneElse, { from: governance })
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectRevert(
                registry.setRegistryKeeper(someoneElse, { from: owner }),
                "Not authorized"
            )
        })

        it("updates registry keeper", async () => {
            await registry.setRegistryKeeper(someoneElse, { from: governance })
            assert.equal(
                await registry.registryKeeper(),
                someoneElse,
                "Unexpected registry keeper"
            )
        })

        it("emits an event", async () => {
            const receipt = await registry.setRegistryKeeper(
                someoneElse, { from: governance }
            )
            expectEvent(receipt, "RegistryKeeperUpdated", {
                registryKeeper: someoneElse
            })
        })
    })

    describe("setDefaultPanicButton", async () => {
        it("can be called by governance", async () => {
            await registry.setDefaultPanicButton(someoneElse, { from: governance })
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectRevert(
                registry.setDefaultPanicButton(someoneElse, { from: owner }),
                "Not authorized"
            )
        })

        it("updates default panic button", async () => {
            await registry.setDefaultPanicButton(someoneElse, { from: governance })
            assert.equal(
                await registry.defaultPanicButton(),
                someoneElse,
                "Unexpected registry keeper"
            )
        })

        it("emits an event", async () => {
            const receipt = await registry.setDefaultPanicButton(
                someoneElse, { from: governance }
            )
            expectEvent(receipt, "DefaultPanicButtonUpdated", {
                defaultPanicButton: someoneElse
            })
        })
    })

    describe("setOperatorContractPanicButton", async () => {
        beforeEach(async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )
        })

        it("can be called by governance", async () => {
            await registry.setOperatorContractPanicButton(
                operatorContract1,
                someoneElse,
                { from: governance }
            )
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectRevert(
                registry.setOperatorContractPanicButton(
                    operatorContract1,
                    someoneElse,
                    { from: owner }
                ),
                "Not authorized"
            )
        })

        it("can not be called with zero panic button address", async () => {
            await expectRevert(
                registry.setOperatorContractPanicButton(
                    operatorContract1,
                    "0x0000000000000000000000000000000000000000",
                    { from: governance }
                ),
                "Panic button must be non-zero address"
            )
        })

        it("can not be called on contracts with disabled panic button", async () => {
            await registry.disableOperatorContractPanicButton(
                operatorContract1,
                { from: governance }
            )
            assert.equal(
                await registry.panicButtons(operatorContract1),
                "0x0000000000000000000000000000000000000000",
                "Panic button not disabled correctly"
            )
            await expectRevert(
                registry.setOperatorContractPanicButton(
                    operatorContract1,
                    someoneElse,
                    { from: governance }
                ),
                "Disabled panic button cannot be updated"
              )
        })

        it("updates contract panic button", async () => {
            await registry.setOperatorContractPanicButton(
                operatorContract1,
                someoneElse,
                { from: governance }
            )
            assert.equal(
                await registry.panicButtons(operatorContract1),
                someoneElse,
                "Unexpected operator contract panic button"
            )
        })

        it("does not update default panic button", async () => {
            await registry.setOperatorContractPanicButton(
                operatorContract1,
                someoneElse,
                { from: governance }
            )
            assert.equal(
                await registry.defaultPanicButton(),
                defaultPanicButton,
                "Unexpected default panic button"
            )
        })
    })

    describe("setOperatorContractUpgrader", async () => {
        it("can be called by governance", async () => {
            await registry.setOperatorContractUpgrader(
                serviceContract1,
                someoneElse,
                { from: governance }
            )
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectRevert(
                registry.setOperatorContractUpgrader(
                    serviceContract1,
                    someoneElse,
                    { from: owner }
                ),
                "Not authorized"
            )
        })

        it("updates operator contract upgrader", async () => {
            await registry.setOperatorContractUpgrader(
                serviceContract1,
                someoneElse,
                { from: governance }
            )

            await registry.setOperatorContractUpgrader(
                serviceContract2,
                operatorContractUpgrader,
                { from: governance }
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

    describe("setServiceContractUpgrader", async () => {
        it("can be called by governance", async () => {
            await registry.setServiceContractUpgrader(
                operatorContract1,
                serviceContractUpgrader,
                { from: governance }
            )
            // ok, no revert
        })

        it("can not be called by non-governance", async () => {
            await expectRevert(
                registry.setServiceContractUpgrader(
                    operatorContract1,
                    serviceContractUpgrader,
                    { from: owner }
                ),
                "Not authorized"
            )
        })

        it("updates service contract upgrader", async () => {
            await registry.setServiceContractUpgrader(
                operatorContract1,
                serviceContractUpgrader,
                { from: governance }
            )

            await registry.setServiceContractUpgrader(
                operatorContract2,
                someoneElse,
                { from: governance }
            )

            assert.equal(
                await registry.serviceContractUpgraderFor(operatorContract1),
                serviceContractUpgrader,
                "Unexpected service contract upgrader"
            )

            assert.equal(
                await registry.serviceContractUpgraderFor(operatorContract2),
                someoneElse,
                "Unexpected service contract upgrader"
            )
        })

        it("emits an event", async () => {
            const receipt = await registry.setServiceContractUpgrader(
                operatorContract1,
                serviceContractUpgrader,
                { from: governance }
            )

            expectEvent(receipt, "ServiceContractUpgraderUpdated", {
                operatorContract: operatorContract1,
                keeper: serviceContractUpgrader
            })
        })
    })

    describe("approveOperatorContract", async () => {
        it("can be called by registry keeper", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )
            // ok, no revert
        })

        it("can not be called by non-registry-keeper", async () => {
            await expectRevert(
                registry.approveOperatorContract(
                    operatorContract1,
                    { from: owner }
                ),
                "Not authorized"
            )
        })

        it("approves operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
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

        it("sets contract's panic button to the default one", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )

            assert.equal(
                await registry.panicButtons(operatorContract1),
                defaultPanicButton,
                "not a default panic button"
            )
        })

        it("cannot be called for already approved contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )

            await expectRevert(
                registry.approveOperatorContract(
                    operatorContract1,
                    { from: registryKeeper }

                ),
                "Not a new operator contract"
            )
        })

        it("cannot be called for disabled contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )
            await registry.disableOperatorContract(
                operatorContract1,
                { from: defaultPanicButton }
            )

            await expectRevert(
                registry.approveOperatorContract(
                    operatorContract1,
                    { from: registryKeeper }

                ),
                "Not a new operator contract"
            )
        })
    })

    describe("disableOperatorContract", async () => {
        beforeEach(async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )
        })

        it("can be called by default panic button", async () => {
            await registry.disableOperatorContract(
                operatorContract1,
                { from: defaultPanicButton }
            )
            // ok, no revert
        })

        it("cannot be called by default panic button if contract has its own", async () => {
            await registry.setOperatorContractPanicButton(
                operatorContract1,
                individualContractPanicButton,
                { from: governance }
            )

            await expectRevert(
                registry.disableOperatorContract(
                    operatorContract1,
                    { from: defaultPanicButton }
                ),
                "Not authorized"
            )
        })

        it("can not be called by non-registry-keeper", async () => {
            await expectRevert(
                registry.disableOperatorContract(
                    operatorContract1,
                    { from: owner }
                ),
                "Not authorized"
            )
        })

        it("disables operator contract", async () => {
            await registry.disableOperatorContract(
                operatorContract1,
                { from: defaultPanicButton }
            )

            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract1),
                "operator contract should not be approved"
            )
        })

        it("disables operator contract with individual panic button", async () => {
            await registry.setOperatorContractPanicButton(
                operatorContract1,
                individualContractPanicButton,
                { from: governance }
            )

            await registry.disableOperatorContract(
                operatorContract1,
                { from: individualContractPanicButton }
            )

            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract1),
                "operator contract should not be approved"
            )
        })

        it("cannot be called if panic button has been disabled", async () => {
            await registry.disableOperatorContractPanicButton(
                operatorContract1,
                { from: governance }
            )

            await expectRevert(
                registry.disableOperatorContract(
                    operatorContract1,
                    { from: defaultPanicButton }
                ),
                "Panic button disabled"
            )
        })

        it("cannot be called for already disabled contract", async () => {
            await registry.disableOperatorContract(
                operatorContract1,
                { from: defaultPanicButton }
            )

            await expectRevert(
                registry.disableOperatorContract(
                    operatorContract1,
                    { from: defaultPanicButton }
                ),
                "Not an approved operator contract"
            )
        })

        it("cannot be called for new operator contract", async () => {
            await expectRevert(
                registry.disableOperatorContract(
                    operatorContract2,
                    { from: defaultPanicButton }
                ),
                "Not an approved operator contract"
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
                { from: registryKeeper }
            )

            assert.isFalse(
                await registry.isNewOperatorContract(operatorContract1),
                "Expected false for approved operator contract"
            )
        })

        it("returns false for disabled operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )
            await registry.disableOperatorContract(
                operatorContract1,
                { from: defaultPanicButton }
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
                { from: registryKeeper }
            )

            assert.isTrue(
                await registry.isApprovedOperatorContract(operatorContract1),
                "Expected true for approved operator contract"
            )
        })

        it("returns false for disabled operator contract", async () => {
            await registry.approveOperatorContract(
                operatorContract1,
                { from: registryKeeper }
            )
            await registry.disableOperatorContract(
                operatorContract1,
                { from: defaultPanicButton }
            )

            assert.isFalse(
                await registry.isApprovedOperatorContract(operatorContract1),
                "Expected false for disabled operator contract"
            )
        })
    })
})
