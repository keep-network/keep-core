import { ethers, waffle, helpers, getUnnamedAccounts } from "hardhat"
import { expect } from "chai"
import type { Signer } from "ethers"
import blsData from "./data/bls"
import { constants, params, testDeployment } from "./fixtures"
import {
  getDkgGroupSigners,
  genesis,
  signAndSubmitDkgResult,
} from "./utils/dkg"
import type { DkgGroupSigners } from "./utils/dkg"
import type { TestRandomBeacon, RandomBeacon } from "../typechain"

const { mineBlocks, mineBlocksTo } = helpers.time

describe("RandomBeacon", () => {
  const groupPublicKey: string = ethers.utils.hexValue(blsData.groupPubKey)

  let thirdParty: Signer
  let signers: DkgGroupSigners

  let randomBeacon: TestRandomBeacon & RandomBeacon

  before(async () => {
    thirdParty = await ethers.getSigner((await getUnnamedAccounts())[1])

    // Accounts offset provided to getDkgGroupSigners have to include number of
    // unnamed accounts that were already used.
    signers = await getDkgGroupSigners(constants.groupSize, 1)
  })

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as TestRandomBeacon & RandomBeacon
  })

  describe("submitDkgResult", async () => {
    let startBlock: number

    beforeEach("run genesis", async () => {
      const [genesisTx] = await genesis(randomBeacon)

      startBlock = genesisTx.blockNumber
    })

    // TODO: These tests are work in progress to quickly test proof of concept.
    // We need to rewrite them to match our style.
    it("with approval", async () => {
      await mineBlocksTo(startBlock + constants.offchainDkgTime)

      const { transaction: tx, members } = await signAndSubmitDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock
      )

      // TODO: Test with misbehaved and threshold members.
      const expectedMembers = members

      let storedGroup = await randomBeacon.getGroup(groupPublicKey)
      let storedGroups = await randomBeacon.getGroups()

      expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
      expect(storedGroup.activationTimestamp).to.be.equal(0)
      expect(storedGroup.members).to.be.deep.equal(members)

      expect(storedGroups).to.be.lengthOf(1)
      expect(storedGroups[0]).to.deep.equal(storedGroup)

      await mineBlocks(params.dkgResultChallengePeriodLength)

      const approveResultTx = await randomBeacon.approveDkgResult()
      // FIXME: Unclear why `approveResultTx.timestamp` is undefined
      const approveResultTxTimestamp = (
        await ethers.provider.getBlock(approveResultTx.blockHash)
      ).timestamp

      storedGroup = await randomBeacon.getGroup(groupPublicKey)
      storedGroups = await randomBeacon.getGroups()

      expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
      expect(storedGroup.activationTimestamp).to.be.equal(
        approveResultTxTimestamp
      )
      expect(storedGroup.members).to.be.deep.equal(expectedMembers)

      expect(storedGroups).to.be.lengthOf(1)
      expect(storedGroups[0]).to.deep.equal(storedGroup)
    })

    it("with challenge", async () => {
      await mineBlocksTo(startBlock + constants.offchainDkgTime)

      const { transaction: tx, members } = await signAndSubmitDkgResult(
        randomBeacon,
        groupPublicKey,
        signers,
        startBlock
      )

      // TODO: Test with misbehaved and threshold members.
      const expectedMembers = members

      let storedGroup = await randomBeacon.getGroup(groupPublicKey)
      let storedGroups = await randomBeacon.getGroups()

      expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
      expect(storedGroup.activationTimestamp).to.be.equal(0)
      expect(storedGroup.members).to.be.deep.equal(expectedMembers)

      expect(storedGroups).to.be.lengthOf(1)
      expect(storedGroups[0]).to.be.deep.equal(storedGroup)

      await randomBeacon.challengeDkgResult()

      storedGroup = await randomBeacon.getGroup(groupPublicKey)
      storedGroups = await randomBeacon.getGroups()

      expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
      expect(storedGroup.activationTimestamp).to.be.equal(0)
      expect(storedGroup.members).to.be.deep.equal(expectedMembers)

      expect(storedGroups).to.be.lengthOf(1)
      expect(storedGroups[0]).to.deep.equal(storedGroup)

      // SUBMIT ANOTHER RESULT WITH THE SAME GROUP PUBLIC KEY
      const { transaction: tx2, members: members2 } =
        await signAndSubmitDkgResult(
          randomBeacon,
          groupPublicKey,
          signers,
          startBlock
        )

      await expect(tx2)
        .to.emit(randomBeacon, "PendingGroupCreated")
        .withArgs(groupPublicKey)

      // TODO: Test with misbehaved and threshold members.
      const expectedMembers2 = members2

      storedGroup = await randomBeacon.getGroup(groupPublicKey)
      storedGroups = await randomBeacon.getGroups()

      expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
      expect(storedGroup.activationTimestamp).to.be.equal(0)
      expect(storedGroup.members).to.be.deep.equal(expectedMembers2)

      expect(storedGroups).to.be.lengthOf(2)
      expect(storedGroups[0]).to.deep.equal(storedGroup)

      await mineBlocks(params.dkgResultChallengePeriodLength)

      const approveResultTx = await randomBeacon.approveDkgResult()
      // FIXME: Unclear why `approveResultTx.timestamp` is undefined
      const approveResultTxTimestamp = (
        await ethers.provider.getBlock(approveResultTx.blockHash)
      ).timestamp

      storedGroup = await randomBeacon.getGroup(groupPublicKey)
      storedGroups = await randomBeacon.getGroups()

      expect(storedGroup.groupPubKey).to.be.equal(groupPublicKey)
      expect(storedGroup.activationTimestamp).to.be.equal(
        approveResultTxTimestamp
      )
      expect(storedGroup.members).to.be.deep.equal(expectedMembers)

      expect(storedGroups).to.be.lengthOf(2)
      expect(storedGroups[1]).to.deep.equal(storedGroup)
    })
  })
})
