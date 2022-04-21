/* eslint-disable @typescript-eslint/no-unused-expressions */

import { ethers } from "hardhat"
import { expect } from "chai"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { ReimbursableImplStub } from "../typechain"

describe("Reimbursable", () => {
  let reimbursableImplStub: ReimbursableImplStub
  let owner: SignerWithAddress
  let thirdParty: SignerWithAddress
  let contractToUpdate: SignerWithAddress

  // prettier-ignore
  before(async () => {
    const ReimbursableImplStub = await ethers.getContractFactory(
      "ReimbursableImplStub"
    )
    reimbursableImplStub = await ReimbursableImplStub.deploy() as ReimbursableImplStub

    [owner, thirdParty, contractToUpdate] =
      await ethers.getSigners()
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
})
