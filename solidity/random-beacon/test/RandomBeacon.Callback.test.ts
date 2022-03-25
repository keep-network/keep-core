import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"

import blsData from "./data/bls"
import { constants, params, randomBeaconDeployment } from "./fixtures"
import { createGroup } from "./utils/groups"
import { registerOperators } from "./utils/operators"

import type { DeployedContracts } from "./fixtures"
import type {
  RandomBeaconStub,
  TestToken,
  CallbackContractStub,
  RandomBeacon,
} from "../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

const ZERO_ADDRESS = ethers.constants.AddressZero
const { createSnapshot, restoreSnapshot } = helpers.snapshot

const fixture = async () => {
  const deployment = await randomBeaconDeployment()

  const contracts: DeployedContracts = {
    randomBeacon: deployment.randomBeacon,
    testToken: deployment.testToken,
    callbackContractStub: await (
      await ethers.getContractFactory("CallbackContractStub")
    ).deploy(),
    callbackContractStub1: await (
      await ethers.getContractFactory("CallbackContractStub")
    ).deploy(),
  }

  // Accounts offset provided to slice getUnnamedAccounts have to include number
  // of unnamed accounts that were already used.
  const signers = await registerOperators(
    contracts.randomBeacon as RandomBeacon,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  await createGroup(contracts.randomBeacon as RandomBeacon, signers)

  return { contracts, signers }
}

describe("RandomBeacon - Callback", () => {
  let requester: SignerWithAddress
  let submitter: SignerWithAddress

  let randomBeacon: RandomBeaconStub
  let testToken: TestToken
  let callbackContract: CallbackContractStub
  let callbackContract1: CallbackContractStub

  before(async () => {
    requester = await ethers.getSigner((await getUnnamedAccounts())[1])
    submitter = await ethers.getSigner((await getUnnamedAccounts())[2])

    const { contracts } = await waffle.loadFixture(fixture)

    randomBeacon = contracts.randomBeacon as RandomBeaconStub
    testToken = contracts.testToken as TestToken
    callbackContract = contracts.callbackContractStub as CallbackContractStub
    callbackContract1 = contracts.callbackContractStub1 as CallbackContractStub
  })

  describe("requestRelayEntry", () => {
    before(async () => {
      await createSnapshot()

      await approveTestToken()
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when passed non-zero and zero callback addresses", () => {
      it("should be set to a non-zero callback contract address", async () => {
        await createSnapshot()

        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address, {
            value: params.relayRequestFee,
          })

        const callbackData = await randomBeacon.getCallbackData()

        await expect(callbackData.callbackContract).to.equal(
          callbackContract.address
        )

        await restoreSnapshot()
      })

      it("should reset to zero callback address", async () => {
        await createSnapshot()

        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address, {
            value: params.relayRequestFee,
          })

        await randomBeacon
          .connect(submitter)
          ["submitRelayEntry(bytes)"](blsData.groupSignature)

        await approveTestToken()

        await randomBeacon
          .connect(requester)
          .requestRelayEntry(ZERO_ADDRESS, { value: params.relayRequestFee })

        const callbackData = await randomBeacon.getCallbackData()
        await expect(callbackData.callbackContract).to.equal(ZERO_ADDRESS)

        await restoreSnapshot()
      })

      it("should be set to the latest non-zero callback address", async () => {
        await createSnapshot()

        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address, {
            value: params.relayRequestFee,
          })

        await randomBeacon
          .connect(submitter)
          ["submitRelayEntry(bytes)"](blsData.groupSignature)

        await approveTestToken()

        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract1.address, {
            value: params.relayRequestFee,
          })

        const callbackData = await randomBeacon.getCallbackData()
        await expect(callbackData.callbackContract).to.equal(
          callbackContract1.address
        )

        await restoreSnapshot()
      })
    })
  })

  describe("submitRelayEntry", () => {
    before(async () => {
      await createSnapshot()

      await approveTestToken()
    })

    after(async () => {
      await restoreSnapshot()
    })

    context("when the callback is set", () => {
      context("when the callback was executed", () => {
        it("should set callback contract params", async () => {
          await createSnapshot()

          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address, {
              value: params.relayRequestFee,
            })

          await randomBeacon
            .connect(submitter)
            ["submitRelayEntry(bytes)"](blsData.groupSignature)

          const lastEntry = await callbackContract.lastEntry()
          await expect(lastEntry).to.equal(blsData.groupSignatureUint256)

          const blockNumber = await callbackContract.blockNumber()
          const latestBlock = await ethers.provider.getBlock("latest")

          await expect(blockNumber).to.equal(latestBlock.number)

          await restoreSnapshot()
        })
      })

      context("when the callback failed", () => {
        it("should emit a callback failed event because of the gas limit", async () => {
          await createSnapshot()

          await randomBeacon.updateRelayEntryParameters(
            params.relayRequestFee,
            params.relayEntrySoftTimeout,
            params.relayEntryHardTimeout,
            40000
          )
          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address, {
              value: params.relayRequestFee,
            })

          const tx = await randomBeacon
            .connect(submitter)
            ["submitRelayEntry(bytes)"](blsData.groupSignature)

          await expect(tx)
            .to.emit(randomBeacon, "CallbackFailed")
            .withArgs(blsData.groupSignatureUint256, tx.blockNumber)

          await restoreSnapshot()
        })

        it("should emit a callback failed event because of the internal error", async () => {
          await createSnapshot()

          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address, {
              value: params.relayRequestFee,
            })

          await callbackContract.setFailureFlag(true)

          const tx = await randomBeacon
            .connect(submitter)
            ["submitRelayEntry(bytes)"](blsData.groupSignature)

          await expect(tx)
            .to.emit(randomBeacon, "CallbackFailed")
            .withArgs(blsData.groupSignatureUint256, tx.blockNumber)

          await restoreSnapshot()
        })
      })
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, params.relayRequestFee)
    await testToken
      .connect(requester)
      .approve(randomBeacon.address, params.relayRequestFee)
  }
})
