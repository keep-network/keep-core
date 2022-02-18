/* eslint-disable @typescript-eslint/no-unused-expressions */

import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { reimbursmentPoolDeployment } from "./fixtures"
import type { ReimbursementPool } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot
const { provider } = waffle
const fixture = async () => reimbursmentPoolDeployment()

describe("ReimbursementPool - Pool", () => {
  let owner: SignerWithAddress
  let thirdParty: SignerWithAddress
  let refundee: SignerWithAddress
  let contractToAuthorize: SignerWithAddress
  let reimbursementPool: ReimbursementPool

  // prettier-ignore
  before(async () => {
    [owner, thirdParty, contractToAuthorize, refundee] =
      await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(fixture)

    reimbursementPool = contracts.reimbursementPool as ReimbursementPool
  })

  describe("transfer ETH", () => {
    context("when a third party funds a reimbursment pool", () => {
      it("should send ETH to the Reimbursment Pool", async () => {
        let reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )
        expect(reimbursementPoolBalance).to.be.equal(0)

        await thirdParty.sendTransaction({
          to: reimbursementPool.address,
          value: ethers.utils.parseEther("1.0"), // Send 1.0 ETH
        })

        reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )

        expect(reimbursementPoolBalance).to.be.equal(
          ethers.utils.parseEther("1.0")
        )
      })
    })

    context("when the owner funds a reimbursment pool", () => {
      it("should withdraw entire ETH balance", async () => {
        let reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )
        expect(reimbursementPoolBalance).to.be.equal(0)

        await owner.sendTransaction({
          to: reimbursementPool.address,
          value: ethers.utils.parseEther("1.0"), // Send 1.0 ETH
        })

        reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )

        expect(reimbursementPoolBalance).to.be.equal(
          ethers.utils.parseEther("1.0")
        )
      })
    })
  })

  describe("withdrawAll", () => {
    beforeEach(async () => {
      await createSnapshot()

      await thirdParty.sendTransaction({
        to: reimbursementPool.address,
        value: ethers.utils.parseEther("10.0"), // Send 10.0 ETH
      })
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    context("when withdrawing all the funds as a non owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool
            .connect(thirdParty)
            .withdrawAll(contractToAuthorize.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when widhrawing all the funds as an owner", () => {
      it("should withdraw entire ETH balance", async () => {
        let reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )
        expect(reimbursementPoolBalance).to.be.equal(
          ethers.utils.parseEther("10.0")
        )

        const thirdPartyBalanceBefore = await provider.getBalance(
          thirdParty.address
        )

        await reimbursementPool.connect(owner).withdrawAll(thirdParty.address)

        reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )
        expect(reimbursementPoolBalance).to.be.equal(0)

        const thirdPartyBalanceAfter = await provider.getBalance(
          thirdParty.address
        )
        const thirdPartyBalanceDiff = thirdPartyBalanceAfter.sub(
          thirdPartyBalanceBefore
        )
        expect(thirdPartyBalanceDiff).to.be.equal(
          ethers.utils.parseEther("10.0")
        )
      })
    })
  })

  describe("withdraw", () => {
    beforeEach(async () => {
      await createSnapshot()

      await thirdParty.sendTransaction({
        to: reimbursementPool.address,
        value: ethers.utils.parseEther("10.0"), // Send 10.0 ETH
      })
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    context("when withdrawing funds as a non owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool
            .connect(thirdParty)
            .withdraw(
              ethers.utils.parseEther("2.0"),
              contractToAuthorize.address
            )
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when widhrawing funds as an owner", () => {
      it("should withdraw ETH balance", async () => {
        let reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )
        expect(reimbursementPoolBalance).to.be.equal(
          ethers.utils.parseEther("10.0")
        )

        const thirdPartyBalanceBefore = await provider.getBalance(
          thirdParty.address
        )

        await reimbursementPool
          .connect(owner)
          .withdraw(ethers.utils.parseEther("2.0"), thirdParty.address)

        reimbursementPoolBalance = await provider.getBalance(
          reimbursementPool.address
        )
        expect(reimbursementPoolBalance).to.be.equal(
          ethers.utils.parseEther("8.0")
        )

        const thirdPartyBalanceAfter = await provider.getBalance(
          thirdParty.address
        )
        const thirdPartyBalanceDiff = thirdPartyBalanceAfter.sub(
          thirdPartyBalanceBefore
        )
        expect(thirdPartyBalanceDiff).to.be.equal(
          ethers.utils.parseEther("2.0")
        )
      })
    })
  })

  describe("refund", () => {
    beforeEach(async () => {
      await createSnapshot()

      await thirdParty.sendTransaction({
        to: reimbursementPool.address,
        value: ethers.utils.parseEther("10.0"), // Send 10.0 ETH
      })

      await reimbursementPool.connect(owner).setStaticGas(23000)
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    context("when contract is not authorized", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool
            .connect(thirdParty)
            .refund(ethers.utils.parseEther("2.0"), thirdParty.address)
        ).to.be.revertedWith("Contract is not authorized for a refund")
      })
    })

    context("when contract is authorized", () => {
      context("when tx gas price is lower than the max gas price", () => {
        it("should refund based on tx.gasprice", async () => {
          await reimbursementPool
            .connect(owner)
            .authorize(contractToAuthorize.address)

          await reimbursementPool
            .connect(owner)
            .setMaxGasPrice(ethers.utils.parseUnits("100.0", "gwei"))

          const refundeeBalanceBefore = await provider.getBalance(
            refundee.address
          )

          await reimbursementPool
            .connect(contractToAuthorize)
            .refund(50000, refundee.address)

          const refundeeBalanceAfter = await provider.getBalance(
            refundee.address
          )
          const refundeeBalanceDiff = refundeeBalanceAfter.sub(
            refundeeBalanceBefore
          )
          expect(refundeeBalanceDiff).to.be.gt(
            ethers.utils.parseUnits("73000", "gwei")
          )

          expect(refundeeBalanceDiff).to.be.lt(
            ethers.utils.parseUnits("146000", "gwei")
          )
        })
      })

      context("when tx gas price is higher than the max gas price", () => {
        it("should refund based on max gas price", async () => {
          await reimbursementPool
            .connect(owner)
            .authorize(contractToAuthorize.address)

          await reimbursementPool
            .connect(owner)
            .setMaxGasPrice(ethers.utils.parseUnits("1.0", "gwei"))

          const refundeeBalanceBefore = await provider.getBalance(
            refundee.address
          )

          await reimbursementPool
            .connect(contractToAuthorize)
            .refund(50000, refundee.address)

          const refundeeBalanceAfter = await provider.getBalance(
            refundee.address
          )
          const refundeeBalanceDiff = refundeeBalanceAfter.sub(
            refundeeBalanceBefore
          )
          expect(refundeeBalanceDiff).to.be.eq(
            ethers.utils.parseUnits("73000", "gwei")
          )
        })
      })
    })
  })

  describe("authorize", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool
            .connect(thirdParty)
            .authorize(contractToAuthorize.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      it("should authorize a contract", async () => {
        await reimbursementPool
          .connect(owner)
          .authorize(contractToAuthorize.address)

        expect(
          await reimbursementPool.isAuthorized(contractToAuthorize.address)
        ).to.be.true
      })
    })
  })

  describe("unuthorize", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool
            .connect(thirdParty)
            .unauthorize(contractToAuthorize.address)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      it("should authorize a contract", async () => {
        await reimbursementPool
          .connect(owner)
          .unauthorize(contractToAuthorize.address)

        expect(
          await reimbursementPool.isAuthorized(contractToAuthorize.address)
        ).to.be.false
      })
    })
  })

  describe("setStaticGas", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool.connect(thirdParty).setStaticGas(42)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      it("should set the static gas cost", async () => {
        expect(await reimbursementPool.staticGas()).to.be.equal(0)

        await reimbursementPool.connect(owner).setStaticGas(42)

        expect(await reimbursementPool.staticGas()).to.be.equal(42)
      })
    })
  })

  describe("setMaxGasPrice", () => {
    context("when the caller is not the owner", () => {
      it("should revert", async () => {
        await expect(
          reimbursementPool.connect(thirdParty).setMaxGasPrice(42)
        ).to.be.revertedWith("Ownable: caller is not the owner")
      })
    })

    context("when the caller is the owner", () => {
      it("should set the max gas price", async () => {
        expect(await reimbursementPool.maxGasPrice()).to.be.equal(0)

        await reimbursementPool.connect(owner).setMaxGasPrice(42)

        expect(await reimbursementPool.maxGasPrice()).to.be.equal(42)
      })
    })
  })
})
