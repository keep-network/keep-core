import { ethers, waffle } from "hardhat"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import { expect } from "chai"
import blsData from "./data/bls"
import { to1e18, ZERO_ADDRESS } from "./functions"
import { randomBeaconDeployment } from "./fixtures"
import type {
  RandomBeacon,
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

  let randomBeacon: RandomBeacon
  let testToken: TestToken
  let callbackContract: CallbackContractStub

  // prettier-ignore
  before(async () => {
    [requester, submitter] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(randomBeaconDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
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
    context("when a callback contract is zero address", () => {
      it("should not emit a callback contract set event", async () => {
        const tx = await randomBeacon
          .connect(requester)
          .requestRelayEntry(ZERO_ADDRESS)

        await expect(tx).not.to.emit(randomBeacon, "CallbackSet")
      })
    })

    context("when a callback contract is passed", () => {
      it("should emit a callback contract set event", async () => {
        const tx = await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address)

        await expect(tx)
          .to.emit(randomBeacon, "CallbackSet")
          .withArgs(callbackContract.address)
      })
    })
  })

  describe("submitRelayEntry", () => {
    beforeEach(async () => {
      await approveTestToken()
    })
    context("when the callback is set", () => {
      context("when the callback was executed", () => {
        it("should emit a callback executed event", async () => {
          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          const tx = await randomBeacon
            .connect(submitter)
            .submitRelayEntry(16, blsData.groupSignature)

          await expect(tx)
            .to.emit(randomBeacon, "CallbackExecuted")
            .withArgs(blsData.groupSignatureUint256, tx.blockNumber)
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
