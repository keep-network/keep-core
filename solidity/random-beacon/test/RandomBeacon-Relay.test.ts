import { ethers, waffle } from "hardhat"
import { testDeployment } from "./helpers/fixtures"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { RandomBeacon } from "../typechain"

describe("RandomBeacon - Relay", function () {
  let requester: SignerWithAddress

  let randomBeacon: RandomBeacon

  before(async function () {
    [requester] = await ethers.getSigners()
  })

  beforeEach("load test fixture", async function () {
    const contracts = await waffle.loadFixture(testDeployment)

    randomBeacon = contracts.randomBeacon as RandomBeacon
  })

  describe("requestRelayEntry", () => {
    context("when groups exist", () => {
      beforeEach(async () => {
        // TODO: Add group.
      })

      context("when there is no other relay entry in progress", () => {
        context("when the requester pays the relay request fee", () => {
          beforeEach(async () => {
            // TODO: Allow fee to be pulled.
          })

          it("should deposit relay request fee to the maintenance pool", async () => {

          })

          it("should set correct current request info", async () => {

          })

          it("should emit RelayEntryRequested event", async () => {

          })
        })

        context("when the requester doesn't pay the relay request fee", () => {
          it("should revert", () => {

          })
        })
      })

      context("when there is an other relay entry in progress", () => {
        beforeEach(async () => {
          // TODO: Create a request.
        })

        it("should revert", () => {

        })
      })
    })

    context("when no groups exist", () => {
      it("should revert", () => {

      })
    })
  })
})