/* eslint-disable @typescript-eslint/no-unused-expressions */

import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { reimbursmentPoolDeployment, params } from "./fixtures"
import type { ReimbursementPool } from "../typechain"

const ZERO_ADDRESS = ethers.constants.AddressZero
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

      context("when receiver is zero address", () => {
        it("should revert", async () => {
          await expect(
            reimbursementPool.connect(owner).withdrawAll(ZERO_ADDRESS)
          ).to.be.revertedWith("Receiver's address cannot be zero")
        })
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

      context("when receiver is zero address", () => {
        it("should revert", async () => {
          await expect(
            reimbursementPool.connect(owner).withdraw(42, ZERO_ADDRESS)
          ).to.be.revertedWith("Receiver's address cannot be zero")
        })
      })

      context("when withdrawing more than the pool's balance", () => {
        it("should revert", async () => {
          await expect(
            reimbursementPool
              .connect(owner)
              .withdraw(ethers.utils.parseEther("42.0"), ZERO_ADDRESS)
          ).to.be.revertedWith("Insufficient contract balance")
        })
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
      beforeEach(async () => {
        await createSnapshot()

        await reimbursementPool
          .connect(owner)
          .authorize(contractToAuthorize.address)
      })

      afterEach(async () => {
        await restoreSnapshot()
      })

      context("when tx gas price is lower than the max gas price", () => {
        it("should refund based on tx.gasprice", async () => {
          const refundeeBalanceBefore = await provider.getBalance(
            refundee.address
          )

          const tx = await reimbursementPool
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

          await expect(tx).not.to.emit(reimbursementPool, "SendingEtherFailed")
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
          // gas spent + static gas => 50k + 37.5k
          expect(refundeeBalanceDiff).to.be.eq(
            ethers.utils.parseUnits("87500", "gwei")
          )
        })
      })

      context("when receiver address is zero", () => {
        it("should revert", async () => {
          await expect(
            reimbursementPool
              .connect(contractToAuthorize)
              .refund(50000, ZERO_ADDRESS)
          ).to.be.revertedWith("Receiver's address cannot be zero")
        })
      })

      context("when no funds available in the pool", () => {
        it("should emit SendingEtherFailed event", async () => {
          await reimbursementPool
            .connect(owner)
            .setMaxGasPrice(ethers.utils.parseUnits("1.0", "gwei"))

          await reimbursementPool.connect(owner).withdrawAll(thirdParty.address)

          const tx = await reimbursementPool
            .connect(contractToAuthorize)
            .refund(50000, refundee.address)

          // gas spent + static gas => 50k + 37.5k
          await expect(tx)
            .to.emit(reimbursementPool, "SendingEtherFailed")
            .withArgs(
              ethers.utils.parseUnits("87500", "gwei"),
              refundee.address
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
        const tx = await reimbursementPool
          .connect(owner)
          .authorize(contractToAuthorize.address)

        expect(
          await reimbursementPool.isAuthorized(contractToAuthorize.address)
        ).to.be.true

        await expect(tx)
          .to.emit(reimbursementPool, "AuthorizedContract")
          .withArgs(contractToAuthorize.address)
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
      it("should unauthorize a contract", async () => {
        const tx = await reimbursementPool
          .connect(owner)
          .unauthorize(contractToAuthorize.address)

        expect(
          await reimbursementPool.isAuthorized(contractToAuthorize.address)
        ).to.be.false

        await expect(tx)
          .to.emit(reimbursementPool, "UnauthorizedContract")
          .withArgs(contractToAuthorize.address)
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
        expect(await reimbursementPool.staticGas()).to.be.equal(
          params.reimbursmentPoolStaticGas
        )

        const tx = await reimbursementPool.connect(owner).setStaticGas(42000)

        await expect(tx)
          .to.emit(reimbursementPool, "StaticGasUpdated")
          .withArgs(42000)

        expect(await reimbursementPool.staticGas()).to.be.equal(42000)
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
        expect(await reimbursementPool.maxGasPrice()).to.be.equal(
          params.reimbursmentPoolMaxGasPrice
        )
        const newMaxGasPrice = ethers.utils.parseUnits("21", "gwei")

        const tx = await reimbursementPool
          .connect(owner)
          .setMaxGasPrice(newMaxGasPrice)

        await expect(tx)
          .to.emit(reimbursementPool, "MaxGasPriceUpdated")
          .withArgs(newMaxGasPrice)

        expect(await reimbursementPool.maxGasPrice()).to.be.equal(
          newMaxGasPrice
        )
      })
    })
  })
})
