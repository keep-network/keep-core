/* eslint-disable @typescript-eslint/no-unused-expressions */

import { ethers, waffle } from "hardhat"
import { expect } from "chai"

import { reimbursableDeployment } from "./fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { ReimbursableImplStub } from "../typechain"

const fixture = async () => reimbursableDeployment()

describe("Reimbursable", () => {
  let reimbursableImplStub: ReimbursableImplStub
  let owner: SignerWithAddress
  let thirdParty: SignerWithAddress
  let contractToUpdate: SignerWithAddress

  // prettier-ignore
  before(async () => {
    [owner, thirdParty, contractToUpdate] =
      await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(fixture)

    reimbursableImplStub =
      contracts.reimbursableImplStub as ReimbursableImplStub
  })

  describe("updateReimbursementPool", () => {
    context("when a caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursableImplStub
            .connect(thirdParty)
            .updateReimbursementPool(contractToUpdate.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when a caller is the owner", () => {
      it("should update a reimbursement contract", async () => {
        await reimbursableImplStub
          .connect(owner)
          .updateReimbursementPool(contractToUpdate.address)

        expect(await reimbursableImplStub.reimbursementPool()).to.be.equal(
          contractToUpdate.address
        )
      })

      it("should emit ReimbursementPoolUpdated event", async () => {
        await expect(
          reimbursableImplStub
            .connect(owner)
            .updateReimbursementPool(contractToUpdate.address)
        )
          .to.emit(reimbursableImplStub, "ReimbursementPoolUpdated")
          .withArgs(contractToUpdate.address)
      })
    })
  })

  describe("updateTransactionGas", () => {
    context("when a caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursableImplStub.connect(thirdParty).updateTransactionGas(22000)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when a caller is the owner", () => {
      it("should update transaction gas", async () => {
        await reimbursableImplStub.connect(owner).updateTransactionGas(22000)

        expect(await reimbursableImplStub.transactionGas()).to.be.equal(22000)
      })

      it("should emit TransactionGasUpdated event", async () => {
        await expect(
          reimbursableImplStub.connect(owner).updateTransactionGas(22000)
        )
          .to.emit(reimbursableImplStub, "TransactionGasUpdated")
          .withArgs(22000)
      })
    })
  })
})
