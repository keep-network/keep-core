/* eslint-disable no-underscore-dangle */
import { helpers } from "hardhat"
import { expect } from "chai"

import { walletRegistryFixture } from "./fixtures"
import { upgradeRandomBeacon } from "./utils/governance"

import type { Contract, BigNumber } from "ethers"
import type { IWalletOwner } from "../typechain/IWalletOwner"
import type { FakeContract } from "@defi-wonderland/smock"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, WalletRegistryStub } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("WalletRegistry - Random Beacon Chaosnet", async () => {
  let relayEntry: BigNumber
  let walletRegistry: WalletRegistryStub & WalletRegistry
  let randomBeaconChaosnet: Contract
  let walletOwner: FakeContract<IWalletOwner>
  let deployer: SignerWithAddress

  before("load test fixture", async () => {
    // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ walletRegistry, walletOwner, randomBeaconChaosnet } =
      await walletRegistryFixture())
    ;({ deployer } = await helpers.signers.getNamedSigners())

    relayEntry = await randomBeaconChaosnet.entry()
  })

  describe("requestNewWallet", async () => {
    before(async () => {
      await createSnapshot()

      await randomBeaconChaosnet
        .connect(deployer)
        .setRequesterAuthorization(walletRegistry.address, true)

      await upgradeRandomBeacon(walletRegistry, randomBeaconChaosnet.address)
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when a new wallet was requested", async () => {
      it("should set seed for wallet creation by the __beaconCallback function", async () => {
        await walletRegistry.connect(walletOwner.wallet).requestNewWallet()

        expect((await walletRegistry.getDkgData()).seed).to.be.equal(relayEntry)
      })
    })
  })
})
