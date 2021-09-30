import { ethers, waffle } from "hardhat"
import { expect } from "chai"
import { increaseTime, lastBlockTime } from "./helpers/contract-test-helpers"

import { BigNumber, ContractTransaction } from "ethers"
import type { DKG, RandomBeacon } from "../typechain"

describe("RandomBeacon contract", function () {
  async function fixture() {
    const DKG = await ethers.getContractFactory("DKG")
    const dkg = (await DKG.deploy()) as DKG

    const RandomBeacon = await ethers.getContractFactory("RandomBeacon", {
      libraries: {
        DKG: dkg.address,
      },
    })

    const randomBeacon = (await RandomBeacon.deploy()) as RandomBeacon

    await randomBeacon.deployed()

    return randomBeacon
  }

  let randomBeacon: RandomBeacon

  beforeEach("load test fixture", async function () {
    randomBeacon = await waffle.loadFixture(fixture)
  })

  describe("genesis function call", async function () {
    let tx: ContractTransaction

    beforeEach("run genesis", async () => {
      tx = await randomBeacon.genesis()
    })

    it("emits DkgStarted event", async function () {
      await expect(tx)
        .to.emit(randomBeacon, "DkgStarted")
        .withArgs(
          await randomBeacon.GENESIS_SEED(),
          await randomBeacon.GROUP_SIZE(),
          await randomBeacon.DKG_TIMEOUT()
        )
    })

    it("sets values", async function () {
      const dkgData = await randomBeacon.dkg()

      expect(dkgData.seed).to.eq(await randomBeacon.GENESIS_SEED())
      expect(dkgData.groupSize).to.eq(await randomBeacon.GROUP_SIZE())
      expect(dkgData.timeoutDuration).to.eq(await randomBeacon.DKG_TIMEOUT())
      expect(dkgData.startTimestamp).to.be.eq(await lastBlockTime())
    })
  })

  describe("isDkgInProgress function call", async function () {
    context("with initial contract state", async function () {
      it("returns false", async function () {
        expect(await randomBeacon.isDkgInProgress()).to.be.false
      })
    })

    context("when genesis dkg started", async function () {
      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
      })

      it("returns false", async function () {
        expect(await randomBeacon.isDkgInProgress()).to.be.true
      })

      context("when dkg timeout was notified", async function () {
        beforeEach("notify dkg timeout", async () => {
          await increaseTime(await randomBeacon.DKG_TIMEOUT())
          await randomBeacon.notifyDkgTimeout()
        })

        it("returns false", async function () {
          expect(await randomBeacon.isDkgInProgress()).to.be.false
        })
      })
    })
  })

  describe("notifyDkgTimeout function call", async function () {
    context("with group creation in progress", async function () {
      let seed: BigNumber

      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
        seed = (await randomBeacon.dkg()).seed
      })

      context("with group creation not timed out", async function () {
        it("reverts with NotTimedOut error", async function () {
          const expectedTimeout = (await randomBeacon.dkg()).startTimestamp.add(
            await randomBeacon.DKG_TIMEOUT()
          )

          // FIXME: I don't like it. Need to rerun many times to see if it really works.
          const currentTimestamp = await increaseTime(10)

          await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
            `NotTimedOut(${expectedTimeout}, ${currentTimestamp + 1})`
          )
        })
      })

      context("with group creation timed out", async function () {
        let tx: ContractTransaction

        beforeEach("notify dkg timeout", async () => {
          await increaseTime(await randomBeacon.DKG_TIMEOUT())
          tx = await randomBeacon.notifyDkgTimeout()
        })

        it("emits an event", async function () {
          await expect(tx).to.emit(randomBeacon, "DkgTimedOut").withArgs(seed)
        })

        it("cleans up dkg data", async function () {
          await assertDkgCleanData(randomBeacon)
        })
      })
    })

    context("with group creation not in progress", async function () {
      it("reverts with InvalidState error", async function () {
        await expect(randomBeacon.notifyDkgTimeout()).to.be.revertedWith(
          "InvalidInProgressState(true, false)"
        )
      })
    })
  })

  describe("submitDkgResult function call", async function () {
    context("with initial contract state", async function () {
      it("reverts with InvalidInProgressState error", async function () {
        await expect(randomBeacon.submitDkgResult()).to.be.revertedWith(
          "InvalidInProgressState(true, false)"
        )
      })
    })

    context("with group creation in progress", async function () {
      beforeEach("run genesis", async () => {
        await randomBeacon.genesis()
      })

      context("with group creation not timed out", async function () {
        it("cleans up dkg data", async function () {
          await randomBeacon.submitDkgResult()

          // TODO: This test will be enhanced to clean up only after all DKG
          // results were submitted.
          await assertDkgCleanData(randomBeacon)
        })
      })

      context("with group creation timed out", async function () {
        beforeEach("increase time", async () => {
          await increaseTime(await randomBeacon.DKG_TIMEOUT())
        })

        context("with timeout not notified", async function () {
          it("cleans up dkg data", async function () {
            await randomBeacon.submitDkgResult()

            // TODO: This test will be enhanced to clean up only after all DKG
            // results were submitted.
            await assertDkgCleanData(randomBeacon)
          })
        })

        context("with timeout notified", async function () {
          beforeEach("notify dkg timeout", async () => {
            await randomBeacon.notifyDkgTimeout()
          })

          it("reverts with InvalidInProgressState error", async function () {
            await expect(randomBeacon.submitDkgResult()).to.be.revertedWith(
              "InvalidInProgressState(true, false)"
            )
          })
        })
      })
    })
  })
})

async function assertDkgCleanData(randomBeacon: RandomBeacon) {
  const dkgData = await randomBeacon.dkg()

  expect(dkgData.seed).to.eq(0)
  expect(dkgData.groupSize).to.eq(0)
  expect(dkgData.timeoutDuration).to.eq(0)
  expect(dkgData.startTimestamp).to.eq(0)
}
