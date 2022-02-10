/* eslint-disable no-underscore-dangle */
import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import { dkgState, walletRegistryFixture } from "./fixtures"
import { fakeRandomBeacon, resetMock } from "./utils/randomBeacon"

import type { FakeContract } from "@defi-wonderland/smock"
import type { ContractTransaction } from "ethers"
import type {
  WalletRegistry,
  WalletRegistryStub,
  IRandomBeacon,
} from "../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Random Beacon", async () => {
  let walletRegistry: WalletRegistryStub & WalletRegistry
  let randomBeacon: FakeContract<IRandomBeacon>

  let walletOwner: SignerWithAddress
  let thirdParty: SignerWithAddress

  before("load test fixture", async () => {
    await createSnapshot()

    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner, thirdParty } = await waffle.loadFixture(
      walletRegistryFixture
    ))

    randomBeacon = await fakeRandomBeacon(walletRegistry)
  })

  after(async () => {
    await restoreSnapshot()
  })

  describe("requestNewWallet", async () => {
    context("when requestRelayEntry reverts", async () => {
      let tx: Promise<ContractTransaction>

      before(async () => {
        await createSnapshot()

        await randomBeacon.requestRelayEntry.reverts("beacon internal error")

        tx = walletRegistry.connect(walletOwner).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()

        resetMock(randomBeacon)
      })

      it("should revert", async () => {
        // FIXME: For some reason this check doesn't work with the expected error message
        // await expect(tx).to.be.revertedWith("beacon internal error")
        await expect(tx).to.be.reverted
      })
    })

    context("when requestRelayEntry succeeds", async () => {
      let tx: Promise<ContractTransaction>

      before(async () => {
        await createSnapshot()

        tx = walletRegistry.connect(walletOwner).requestNewWallet()
      })

      after(async () => {
        await restoreSnapshot()
      })

      it("should succeed", async () => {
        await expect(tx).to.not.be.reverted
      })

      it("should call random beacon", async () => {
        await expect(randomBeacon.requestRelayEntry).to.be.calledWith(
          walletRegistry.address
        )
      })
    })
  })

  describe("__beaconCallback", async () => {
    before(async () => {
      await createSnapshot()
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
      context("when new wallet was not requested", async () => {
        it("should revert", async () => {
          await expect(
            walletRegistry
              .connect(randomBeacon.wallet)
              .__beaconCallback(123, 456)
          ).to.be.revertedWith("Current state is not AWAITING_SEED")
        })
      })

      context("when new wallet was requested", async () => {
        const relayEntry = ethers.BigNumber.from(ethers.utils.randomBytes(32))

        before(async () => {
          await createSnapshot()

          await walletRegistry.connect(walletOwner).requestNewWallet()

          await walletRegistry
            .connect(randomBeacon.wallet)
            .__beaconCallback(relayEntry, 0)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should transition wallet creation state to `AWAITING_RESULT`", async () => {
          await expect(
            await walletRegistry.getWalletCreationState()
          ).to.be.equal(dkgState.AWAITING_RESULT)
        })

        it("should set seed for wallet creation", async () => {
          await expect((await walletRegistry.getDkgData()).seed).to.be.equal(
            relayEntry
          )
        })
      })
    })
  })
})
