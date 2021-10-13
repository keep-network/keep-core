import { ethers, waffle } from "hardhat"
import { expect } from "chai"
import { blsData } from "./helpers/data"
import { to1e18 } from "./helpers/functions"
import { testDeployment } from "./helpers/fixtures"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { RandomBeacon, TestToken, MaintenancePool } from "../typechain"

describe("RandomBeacon - Relay", function () {
  const relayRequestFee = to1e18(100)

  let requester: SignerWithAddress

  let randomBeacon: RandomBeacon
  let testToken: TestToken
  let maintenancePool: MaintenancePool

  before(async function () {
    [requester] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async function () {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
    testToken = contracts.testToken as TestToken
    maintenancePool = contracts.maintenancePool as MaintenancePool
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      beforeEach(async () => {
        // TODO: Currently `selectGroup` returns a hardcoded group. Once
        //       proper implementation is ready, add the group manually here.
      })

      context("when there is no other relay entry in progress", () => {
        context("when the requester pays the relay request fee", () => {
          let tx
          let previousMaintenancePoolBalance

          beforeEach(async () => {
            previousMaintenancePoolBalance = await testToken.balanceOf(maintenancePool.address)
            await approveTestToken()
            tx = await randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
          })

          it("should deposit relay request fee to the maintenance pool", async () => {
            const actualMaintenancePoolBalance = await testToken.balanceOf(maintenancePool.address)
            expect(actualMaintenancePoolBalance.sub(previousMaintenancePoolBalance)).to.be.equal(relayRequestFee)
          })

          it("should emit RelayEntryRequested event", async () => {
            await expect(tx).to
              .emit(randomBeacon, "RelayEntryRequested")
              .withArgs(1, blsData.groupPubKey, blsData.previousEntry)
          })
        })

        context("when the requester doesn't pay the relay request fee", () => {
          it.only("should revert", async () => {
            await expect(
              randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
            ).to.be.revertedWith("Transfer amount exceeds allowance")
          })
        })
      })

      context("when there is an other relay entry in progress", () => {
        beforeEach(async () => {
          await approveTestToken()
          await randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
        })

        it("should revert", async () => {
          await expect(
            randomBeacon.connect(requester).requestRelayEntry(blsData.previousEntry)
          ).to.be.revertedWith("Another relay request in progress")
        })
      })
    })

    context("when no groups exist", () => {
      it("should revert", () => {
        // TODO: Implement once proper `selectGroup` is ready.
      })
    })
  })

  async function approveTestToken() {
    await testToken.mint(requester.address, relayRequestFee)
    await testToken.connect(requester).approve(randomBeacon.address, relayRequestFee)
  }
})