/* eslint-disable no-underscore-dangle */
import { expect } from "chai"
import { ethers, helpers } from "hardhat"

import type { ContractTransaction } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { GovernableImpl, GovernableImpl__factory } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("Governable", () => {
  let governable: GovernableImpl
  let deployer: SignerWithAddress
  let governance: SignerWithAddress
  let thirdParty: SignerWithAddress

  before(async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ deployer, governance } = await ethers.getNamedSigners())
    ;[thirdParty] = await ethers.getUnnamedSigners()

    const GovernableFactory: GovernableImpl__factory =
      await ethers.getContractFactory("GovernableImpl", deployer)
    governable = await GovernableFactory.deploy()
  })

  describe("constructor", () => {
    it("sets governance to default zero address", async () => {
      expect(await governable.governance()).to.be.equal(
        ethers.constants.AddressZero
      )
    })
  })

  describe("transferGovernance", () => {
    describe("when governance was not initialized", () => {
      describe("when called by the deployer", () => {
        it("reverts", async () => {
          await expect(
            governable
              .connect(deployer)
              .transferGovernance(ethers.Wallet.createRandom().address)
          ).to.be.revertedWith("Caller is not the governance")
        })
      })

      describe("when called by the governance", () => {
        it("reverts", async () => {
          await expect(
            governable
              .connect(governance)
              .transferGovernance(ethers.Wallet.createRandom().address)
          ).to.be.revertedWith("Caller is not the governance")
        })
      })

      describe("when called by a third party", () => {
        it("reverts", async () => {
          await expect(
            governable
              .connect(thirdParty)
              .transferGovernance(ethers.Wallet.createRandom().address)
          ).to.be.revertedWith("Caller is not the governance")
        })
      })
    })

    describe("when governance was initialized", () => {
      before(async () => {
        await governable._transferGovernanceExposed(governance.address)
      })

      describe("when called by the deployer", () => {
        it("reverts", async () => {
          await expect(
            governable
              .connect(deployer)
              .transferGovernance(ethers.Wallet.createRandom().address)
          ).to.be.revertedWith("Caller is not the governance")
        })
      })

      describe("when called by the governance", () => {
        const newGovernance: string = ethers.Wallet.createRandom().address
        let tx: ContractTransaction

        before(async () => {
          await createSnapshot()

          tx = await governable
            .connect(governance)
            .transferGovernance(newGovernance)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("updates governance address", async () => {
          expect(await governable.governance()).to.be.equal(newGovernance)
        })

        it("emits GovernanceTransferred event", async () => {
          await expect(tx)
            .to.emit(governable, "GovernanceTransferred")
            .withArgs(governance.address, newGovernance)
        })
      })

      describe("when called by a third party", () => {
        it("reverts", async () => {
          await expect(
            governable
              .connect(thirdParty)
              .transferGovernance(ethers.Wallet.createRandom().address)
          ).to.be.revertedWith("Caller is not the governance")
        })
      })

      describe("when new governance is zero address", () => {
        it("reverts", async () => {
          await expect(
            governable
              .connect(governance)
              .transferGovernance(ethers.constants.AddressZero)
          ).to.be.revertedWith("New governance is the zero address")
        })
      })
    })
  })

  describe("_transferGovernance", () => {
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    describe("when called by the deployer", () => {
      it("succeeds", async () => {
        const newGovernance = ethers.Wallet.createRandom().address

        await governable
          .connect(deployer)
          ._transferGovernanceExposed(newGovernance)

        expect(await governable.governance()).to.be.equal(newGovernance)
      })
    })

    describe("when called by the governance", () => {
      it("succeeds", async () => {
        const newGovernance = ethers.Wallet.createRandom().address

        await governable
          .connect(governance)
          ._transferGovernanceExposed(newGovernance)

        expect(await governable.governance()).to.be.equal(newGovernance)
      })
    })

    describe("when called by a third party", () => {
      it("succeeds", async () => {
        const newGovernance = ethers.Wallet.createRandom().address

        await governable
          .connect(thirdParty)
          ._transferGovernanceExposed(newGovernance)

        expect(await governable.governance()).to.be.equal(newGovernance)
      })
    })
  })
})
