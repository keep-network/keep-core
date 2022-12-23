/* eslint-disable @typescript-eslint/no-extra-semi */

import { ethers, helpers } from "hardhat"
import { BigNumber } from "ethers"
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
        await createSnapshot()

        await randomBeaconChaosnet
          .connect(deployer)
          .setRequesterAuthorization(requester.address, true)
      })

      after(async () => {
        await restoreSnapshot()
      })

      context("when called once", () => {
        before(async () => {
          await createSnapshot()

          await randomBeaconChaosnet
            .connect(requester)
            .requestRelayEntry(callbackContract.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should be set to a callback contract address", async () => {
          expect(await randomBeaconChaosnet.getCallbackContract()).to.equal(
            callbackContract.address
          )
        })

        it("should execute callback with proper entry", async () => {
          expect(await callbackContract.lastEntry()).to.equal(
            // The entry is keccak-256 of the initial value stored in
            // the RandomBeaconChaosnet contract
            BigNumber.from(
              "86322480231844907215266847458792959757192550318770676212332984" +
                "332154459033029"
            )
          )
        })
      })

      context("when called twice", () => {
        before(async () => {
          await createSnapshot()

          await randomBeaconChaosnet
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          await randomBeaconChaosnet
            .connect(requester)
            .requestRelayEntry(callbackContract1.address)
        })

        after(async () => {
          await restoreSnapshot()
        })

        it("should be set to the latest callback contract address", async () => {
          expect(await randomBeaconChaosnet.getCallbackContract()).to.equal(
            callbackContract1.address
          )
        })

        it("should execute the first callback with proper entry", async () => {
          expect(await callbackContract.lastEntry()).to.equal(
            // The entry is keccak-256 of the initial value stored in
            // the RandomBeaconChaosnet contract
            BigNumber.from(
              "86322480231844907215266847458792959757192550318770676212332984" +
                "332154459033029"
            )
          )
        })

        it("should execute the second callback with proper entry", async () => {
          // The entry is keccak-256 of the previous entry
          expect(await callbackContract1.lastEntry()).to.equal(
            BigNumber.from(
              "45055825411044151981109535788320043556123542984485670123474642" +
                "322436340913380"
            )
          )
        })
      })
    })
  })
})
