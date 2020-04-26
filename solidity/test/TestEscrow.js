const { accounts, contract, web3 } = require("@openzeppelin/test-environment")
const { createSnapshot, restoreSnapshot } = require("./helpers/snapshot.js")
const { expectRevert, expectEvent } = require("@openzeppelin/test-helpers")

const KeepToken = contract.fromArtifact('KeepToken')
const Escrow = contract.fromArtifact('Escrow')

const BN = web3.utils.BN
const chai = require('chai')
chai.use(require('bn-chai')(BN))
const expect = chai.expect
const assert = require('chai').assert

describe('Escrow', () => {
    let owner = accounts[0]
    let updatedOwner = accounts[1]

    let beneficiary = accounts[2]
    let updatedBeneficiary = accounts[3]
    let notBeneficiary = accounts[4]

    let token
    let escrow

    before(async () => {
        token = await KeepToken.new({from: owner})
        escrow = await Escrow.new(token.address, {from: owner})
    })

    beforeEach(async () => {
        await createSnapshot()
    })

    afterEach(async () => {
        await restoreSnapshot()
    })

    describe("setBeneficiary", async () => {
        it("can be called by owner", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            // ok, no revert
        })

        it("can be called by updated owner", async () => {
            await escrow.transferOwnership(updatedOwner, {from: owner})

            await expectRevert(
                escrow.setBeneficiary(beneficiary, {from: owner}),
                "Ownable: caller is not the owner"
            )
            await escrow.setBeneficiary(beneficiary, {from: updatedOwner})
            // ok, no revert
        })

        it("can not be called by non-owner", async () => {
            await expectRevert(
                escrow.setBeneficiary(beneficiary, {from: beneficiary}),
                "Ownable: caller is not the owner"
            )
        })

        it("sets beneficiary", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})

            assert.equal(
                await escrow.beneficiary(),
                beneficiary,
                "Unexpected beneficiary"
            )
        })

        it("allows to update beneficiary", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            await escrow.setBeneficiary(updatedBeneficiary, {from: owner})

            assert.equal(
                await escrow.beneficiary(),
                updatedBeneficiary,
                "Unexpected beneficiary"
            )
        })

        it("emits an event", async () => {
            let receipt = await escrow.setBeneficiary(
                beneficiary, 
                {from: owner}
            )

            await expectEvent(receipt, 'BeneficiaryUpdated', {
                beneficiary: beneficiary
            })
        })
    })

    describe("withdraw", async () => {
        it("can not be called if beneficiary wasn't set", async () => {
            await token.transfer(escrow.address, 100, {from: owner});
            await expectRevert(
                escrow.withdraw({from: beneficiary}),
                "Beneficiary not assigned"                
            )
        })

        it("can not be called by non-beneficiary", async () => {
            await token.transfer(escrow.address, 100, {from: owner});
            await escrow.setBeneficiary(beneficiary, {from: owner})
            await expectRevert(
                escrow.withdraw({from: notBeneficiary}),
                "Caller is not the beneficiary."                
            )
        })

        it("fails when escrow is empty", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            await expectRevert(
                escrow.withdraw({from: beneficiary}),
                "No tokens to withdraw"                
            )
        })

        it("can be called by beneficiary", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            await token.transfer(escrow.address, 100, {from: owner});
            await escrow.withdraw({from: beneficiary})
            // ok, no reverts
        })

        it("can be called by updated beneficiary", async () => {
            await token.transfer(escrow.address, 100, {from: owner});                        
            await escrow.setBeneficiary(updatedBeneficiary, {from: owner})

            await expectRevert(
                escrow.withdraw({from: beneficiary}),
                "Caller is not the beneficiary."                
            )
            await escrow.withdraw({from: updatedBeneficiary})
            // ok, no reverts
        })

        it("withdraws all tokens to beneficiary", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            let amount = web3.utils.toBN(123456789)
            await token.transfer(escrow.address, amount, {from: owner});

            await escrow.withdraw({from: beneficiary})
            let beneficiaryBalanceAfter = await token.balanceOf(beneficiary)

            expect(beneficiaryBalanceAfter).to.eq.BN(
                amount,
                "Unexpected amount withdrawn"
            )
        })

        it("withdraws all tokens to updated beneficiary", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            await escrow.setBeneficiary(updatedBeneficiary, {from: owner})

            let amount = web3.utils.toBN(987654321)
            await token.transfer(escrow.address, amount, {from: owner});

            await escrow.setBeneficiary(updatedBeneficiary, {from: owner})
            await escrow.withdraw({from: updatedBeneficiary})
            let beneficiaryBalanceAfter = await token.balanceOf(updatedBeneficiary)

            expect(beneficiaryBalanceAfter).to.eq.BN(
                amount,
                "Unexpected amount withdrawn"
            )
        })

        it("emits an event", async () => {
            await escrow.setBeneficiary(beneficiary, {from: owner})
            let amount = web3.utils.toBN(100)
            await token.transfer(escrow.address, amount, {from: owner});

            let receipt = await escrow.withdraw({from: beneficiary})

            await expectEvent(receipt, "TokensWithdrawn", {
                beneficiary: beneficiary,
                amount: amount
            })
        })
    })
})