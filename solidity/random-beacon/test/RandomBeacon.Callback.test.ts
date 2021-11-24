import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import blsData from "./data/bls"
import { constants, randomBeaconDeployment } from "./fixtures"
import { createGroup } from "./utils/groups"
import type { DeployedContracts } from "./fixtures"
import type {
  RandomBeaconStub,
  TestToken,
  CallbackContractStub,
} from "../typechain"
import { registerOperators, Operator } from "./utils/operators"

const ZERO_ADDRESS = ethers.constants.AddressZero

const { to1e18 } = helpers.number

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
    deployment.randomBeacon as RandomBeaconStub,
    (await getUnnamedAccounts()).slice(1, 1 + constants.groupSize)
  )

  await createGroup(contracts.randomBeacon as RandomBeaconStub, signers)

  return { contracts, signers }
}

describe("RandomBeacon - Callback", () => {
  const relayRequestFee = to1e18(100)
  const relayEntryHardTimeout = 5760
  const relayEntrySubmissionEligibilityDelay = 10
  const groupCreationFrequency = 100
  const groupLifetime = 200

  // When determining the eligibility queue, the `(blsData.groupSignature % 64) + 1`
  // equation points member`16` as the first eligible one. This is why we use that
  // index as `submitRelayEntry` parameter. The `submitter` signer represents that
  // member too.
  const submitterMemberIndex = 16

  let callbackGasLimit = 50000

  let requester: SignerWithAddress
  let submitter: SignerWithAddress
  let signers: Operator[]

  let randomBeacon: RandomBeaconStub
  let testToken: TestToken
  let callbackContract: CallbackContractStub
  let callbackContract1: CallbackContractStub

  before(async () => {
    requester = await ethers.getSigner((await getUnnamedAccounts())[1])
  })

  beforeEach("load test fixture", async () => {
    let contracts
      // eslint-disable-next-line @typescript-eslint/no-extra-semi
    ;({ contracts, signers } = await waffle.loadFixture(fixture))

    randomBeacon = contracts.randomBeacon as RandomBeaconStub
    testToken = contracts.testToken as TestToken
    callbackContract = contracts.callbackContractStub as CallbackContractStub
    callbackContract1 = contracts.callbackContractStub1 as CallbackContractStub

    await randomBeacon.updateRelayEntryParameters(
      to1e18(100),
      relayEntrySubmissionEligibilityDelay,
      relayEntryHardTimeout,
      callbackGasLimit
    )
    await randomBeacon.updateGroupCreationParameters(
      groupCreationFrequency,
      groupLifetime
    )

    submitter = await ethers.getSigner(
      signers[submitterMemberIndex - 1].address
    )
  })

  describe("requestRelayEntry", () => {
    beforeEach(async () => {
      await approveTestToken()
    })

    context("when passed non-zero and zero callback addresses", () => {
      it("should be set to a non-zero callback contract address", async () => {
        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address)

        const callbackData = await randomBeacon.getCallbackData()

        await expect(callbackData.callbackContract).to.equal(
          callbackContract.address
        )
      })

      it("should be reset to zero callback address", async () => {
        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address)

        await randomBeacon
          .connect(submitter)
          .submitRelayEntry(16, blsData.groupSignature)

        await approveTestToken()

        await randomBeacon.connect(requester).requestRelayEntry(ZERO_ADDRESS)

        const callbackData = await randomBeacon.getCallbackData()
        await expect(callbackData.callbackContract).to.equal(ZERO_ADDRESS)
      })

      it("should be set to the latest non-zero callback address", async () => {
        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract.address)

        await randomBeacon
          .connect(submitter)
          .submitRelayEntry(16, blsData.groupSignature)

        await approveTestToken()

        await randomBeacon
          .connect(requester)
          .requestRelayEntry(callbackContract1.address)

        const callbackData = await randomBeacon.getCallbackData()
        await expect(callbackData.callbackContract).to.equal(
          callbackContract1.address
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
        it("should set callback contract params", async () => {
          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          await randomBeacon
            .connect(submitter)
            .submitRelayEntry(16, blsData.groupSignature)

          const lastEntry = await callbackContract.lastEntry()
          await expect(lastEntry).to.equal(blsData.groupSignatureUint256)

          const blockNumber = await callbackContract.blockNumber()
          const latestBlock = await ethers.provider.getBlock("latest")

          await expect(blockNumber).to.equal(latestBlock.number)
        })
      })

      context("when the callback failed", () => {
        it("should emit a callback failed event because of the gas limit", async () => {
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

        it("should emit a callback failed event because of the internal error", async () => {
          await randomBeacon
            .connect(requester)
            .requestRelayEntry(callbackContract.address)

          await callbackContract.setFailureFlag(true)

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
