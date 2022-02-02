/* eslint-disable no-underscore-dangle */
import {
  ethers,
  waffle,
  helpers,
  deployments,
  getUnnamedAccounts,
} from "hardhat"
import { expect } from "chai"

import { walletRegistryFixture } from "./fixtures"

import type { ContractTransaction } from "ethers"
import type { MockContract } from "ethereum-waffle"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

const { deployMockContract } = waffle

describe("WalletRegistry - Random Beacon", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry
  let randomBeacon: MockContract

  let deployer: SignerWithAddress
  let walletOwner: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner, deployer, thirdParty } =
      await waffle.loadFixture(walletRegistryFixture))

    const randomBeaconMock = await deployMockContract(
      deployer,
      (
        await deployments.getArtifact("RandomBeacon")
      ).abi
    )

    await walletRegistry.updateRandomBeacon(randomBeaconMock.address)

    randomBeacon = randomBeaconMock
  })

  describe("requestNewWallet", async () => {
    context("when requestRelayEntry reverts", async () => {
      let tx: Promise<ContractTransaction>

      before(async () => {
        await createSnapshot()

        await randomBeacon.mock.requestRelayEntry.revertsWithReason(
          "beacon is busy"
        )

        tx = walletRegistry.connect(walletOwner).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should succeed", async () => {
        await expect(tx).to.not.be.reverted
      })

      it("should emit RelayEntryRequestFailed", async () => {
        await expect(tx).to.emit(walletRegistry, "RelayEntryRequestFailed")
      })
    })

    context("when requestRelayEntry succeeds", async () => {
      let tx: Promise<ContractTransaction>

      before(async () => {
        await createSnapshot()

        await randomBeacon.mock.requestRelayEntry.returns()

        tx = walletRegistry.connect(walletOwner).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should succeed", async () => {
        await expect(tx).to.not.be.reverted
      })

      it("should not emit RelayEntryRequestFailed", async () => {
        await expect(tx).to.not.emit(walletRegistry, "RelayEntryRequestFailed")
      })
    })
  })

  describe("__beaconCallback", async () => {
    let randomBeaconSigner: SignerWithAddress

    before(async () => {
      await createSnapshot()

      randomBeaconSigner = await ethers.getSigner(
        (
          await getUnnamedAccounts()
        )[1]
      )

      await walletRegistry.updateRandomBeacon(randomBeaconSigner.address)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when called by a third party", async () => {
      it("should revert", async () => {
        await expect(
          walletRegistry.connect(thirdParty).__beaconCallback(123, 456)
        ).to.be.revertedWith("Caller is not the Random Beacon")
      })
    })

    context("when called by the random beacon", async () => {
      it("should set new value", async () => {
        const newRelayEntry = 3121

        await walletRegistry
          .connect(randomBeaconSigner)
          .__beaconCallback(newRelayEntry, 456)

        await expect(await walletRegistry.randomRelayEntry()).to.be.equal(
          newRelayEntry
        )
      })
    })
  })
})
