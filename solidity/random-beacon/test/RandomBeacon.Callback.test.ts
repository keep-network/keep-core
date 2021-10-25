import { ethers, waffle } from "hardhat"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { expect } from "chai"
import blsData from "./data/bls"
import { to1e18, ZERO_ADDRESS } from "./functions"
import { randomBeaconDeployment } from "./fixtures"
import type {
  RandomBeaconStub,
  TestToken,
  CallbackContractStub,
} from "../typechain"

describe("RandomBeacon - Callback", () => {
  const relayRequestFee = to1e18(100)
  const relayEntryHardTimeout = 5760
  const relayEntrySubmissionEligibilityDelay = 10
  let callbackGasLimit = 50000

  let requester: SignerWithAddress
  let submitter: SignerWithAddress

  let randomBeacon: RandomBeaconStub
  let testToken: TestToken
  let callbackContract: CallbackContractStub

  // prettier-ignore
  before(async () => {
    [requester, submitter] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(randomBeaconDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeaconStub
    testToken = contracts.testToken as TestToken
    callbackContract = contracts.callbackContractStub as CallbackContractStub

    await randomBeacon.updateRelayEntryParameters(
      to1e18(100),
      relayEntrySubmissionEligibilityDelay,
      relayEntryHardTimeout,
      callbackGasLimit
    )
  })

  describe("requestRelayEntry", () => {
    beforeEach(async () => {
      await approveTestToken()
    })
    context("when passed zero callback address", () => {
      it("should be set to zero address after request", async () => {
        const tx = await randomBeacon
          .connect(requester)
          .requestRelayEntry(ZERO_ADDRESS)

        const callbackData = await randomBeacon.getCallbackData()

        await expect(callbackData.callbackContract).to.equal(ZERO_ADDRESS)
      })
    })

    context("when passed non-zero callback address", () => {
      it("should be set to a callback contract address", async () => {
        const tx = await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address)

        const callbackData = await randomBeacon.getCallbackData()

        await expect(callbackData.callbackContract).to.equal(
          callbackContract.address
        )
      })
    })
  })

  describe("submitRelayEntry", () => {
    beforeEach(async () => {
      await approveTestToken()
    })
    context("when the callback is set", () => {
      context("when the callback was executed", () => {
        it("should set callback contract entry", async () => {
          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          await randomBeacon
            .connect(submitter)
            .submitRelayEntry(16, blsData.groupSignature)

          const lastEntry = await callbackContract.lastEntry()
          await expect(lastEntry).to.equal(blsData.groupSignatureUint256)
        })
      })

      context("when the callback failed", () => {
        it("should emit a callback failed event", async () => {
          callbackGasLimit = 40000
          await randomBeacon.updateRelayEntryParameters(
            to1e18(100),
            relayEntrySubmissionEligibilityDelay,
            relayEntryHardTimeout,
            callbackGasLimit
          )
          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          const tx = await randomBeacon
            .connect(submitter)
            .submitRelayEntry(16, blsData.groupSignature)

          await expect(tx)
            .to.emit(randomBeacon, "CallbackFailed")
            .withArgs(blsData.groupSignatureUint256, tx.blockNumber)
        })
      })
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, relayRequestFee)
    await testToken
      .connect(requester)
      .approve(randomBeacon.address, relayRequestFee)
  }
})
