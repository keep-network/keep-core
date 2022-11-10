/* eslint-disable @typescript-eslint/no-extra-semi */

import { ethers, helpers } from "hardhat"
import { expect } from "chai"

import type {
  RandomBeaconChaosnetStub,
  CallbackContractStub,
} from "../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

describe("RandomBeaconChaosnet - Callback", () => {
  let requester: SignerWithAddress
  let thirdParty: SignerWithAddress
  let deployer: SignerWithAddress

  let randomBeaconChaosnet: RandomBeaconChaosnetStub
  let callbackContract: CallbackContractStub
  let callbackContract1: CallbackContractStub

  before(async () => {
    ;[requester, thirdParty] = await helpers.signers.getUnnamedSigners()
    ;({ deployer } = await helpers.signers.getNamedSigners())

    randomBeaconChaosnet = (await (
      await ethers.getContractFactory("RandomBeaconChaosnetStub")
    )
      .connect(deployer)
      .deploy()) as RandomBeaconChaosnetStub

    callbackContract = (await (
      await ethers.getContractFactory("CallbackContractStub")
    )
      .connect(deployer)
      .deploy()) as CallbackContractStub

    callbackContract1 = (await (
      await ethers.getContractFactory("CallbackContractStub")
    )
      .connect(deployer)
      .deploy()) as CallbackContractStub
  })

  describe("requestRelayEntry", () => {
    before(async () => {
      await createSnapshot()
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when requester is not authorized", () => {
      it("should revert", async () => {
        await expect(
          randomBeaconChaosnet
            .connect(thirdParty)
            .requestRelayEntry(callbackContract.address)
        ).to.be.revertedWith("Requester must be authorized")
      })
    })

    context("when requester is authorized", () => {
      before(async () => {
        await randomBeaconChaosnet
          .connect(deployer)
          .setRequesterAuthorization(requester.address, true)
      })
      context("when passed a callback address", () => {
        it("should be set to a callback contract address", async () => {
          await createSnapshot()

          await randomBeaconChaosnet
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          expect(await randomBeaconChaosnet.getCallbackContract()).to.equal(
            callbackContract.address
          )

          await restoreSnapshot()
        })

        it("should be set to the latest callback address", async () => {
          await createSnapshot()

          await randomBeaconChaosnet
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          await randomBeaconChaosnet
            .connect(requester)
            .requestRelayEntry(callbackContract1.address)

          expect(await randomBeaconChaosnet.getCallbackContract()).to.equal(
            callbackContract1.address
          )

          await restoreSnapshot()
        })
      })
    })
  })
})
