import { ethers, waffle, helpers } from "hardhat"
import { expect } from "chai"

import secretData from "./data/secret"
import { secretDeployment } from "./fixtures"

import type { Secret } from "../typechain"

const { createSnapshot, restoreSnapshot } = helpers.snapshot

let abi = ethers.utils.defaultAbiCoder

function toBytes32(decimalInput: string): string {
    let n = ethers.BigNumber.from(decimalInput)
    return abi.encode(["bytes32"], [n])
}
  
function toG1(bigintStrs: string[]): Uint8Array {
    let x = toBytes32(bigintStrs[0]) 
    let y = toBytes32(bigintStrs[1])
    return ethers.utils.concat([x, y])
}
  
function toG2(bigintStrs: string[]): Uint8Array {
    let xx = toBytes32(bigintStrs[1]) 
    let xy = toBytes32(bigintStrs[0])
    let yx = toBytes32(bigintStrs[3]) 
    let yy = toBytes32(bigintStrs[2])
    return ethers.utils.concat([xx, xy, yx, yy])
}

describe.only("Secret", () => {
  let secret: Secret

  beforeEach("load test fixture", async () => {
    const contracts = await waffle.loadFixture(secretDeployment)

    secret = contracts.secret as Secret
  })

  context("registering a group", () => {
    beforeEach(async () => {
      await createSnapshot()
    })

    afterEach(async () => {
      await restoreSnapshot()
    })

    it("should be able to register a group", async () => {
        let sP = toG1(secretData.sP)
        let xP = toG1(secretData.xP)
        let xQ = toG2(secretData.xQ)
        await secret.registerGroup(sP, xP, xQ)
    })

    it("should not be able to register a bad group", async () => {
        await expect(
            secret.registerGroup(
                toG1(secretData.xP),
                toG1(secretData.sP),
                toG2(secretData.xQ)
            )
        ).to.be.reverted
    })
  })

  context("making requests", () => {
    beforeEach(async () => {
        await createSnapshot()
        let sP = toG1(secretData.sP)
        let xP = toG1(secretData.xP)
        let xQ = toG2(secretData.xQ)
        await secret.registerGroup(sP, xP, xQ)
      })
  
      afterEach(async () => {
        await restoreSnapshot()
      })

      it("should accept good requests", async () => {
        await secret.makeRequest(toG2(secretData.yQ), toG1(secretData.xyP))
        await secret.challengeRequest(toG2(secretData.yQ), toG1(secretData.xyP))
        expect(await secret.callStatic.challengeRequest(toG2(secretData.yQ), toG1(secretData.xyP))).to.be.true
      })

      it("should reject bad requests", async () => {
        await secret.makeRequest(toG2(secretData.sQ), toG1(secretData.xyP))
        await secret.challengeRequest(toG2(secretData.sQ), toG1(secretData.xyP))
        expect(await secret.callStatic.challengeRequest(toG2(secretData.sQ), toG1(secretData.xyP))).to.be.false
      })
  })

  context("response", () => {
    beforeEach(async () => {
        await createSnapshot()
        let sP = toG1(secretData.sP)
        let xP = toG1(secretData.xP)
        let xQ = toG2(secretData.xQ)
        await secret.registerGroup(sP, xP, xQ)
        await secret.makeRequest(toG2(secretData.yQ), toG1(secretData.xyP))
      })
  
      afterEach(async () => {
        await restoreSnapshot()
      })

      it("should accept good responses", async () => {
        await secret.respond(toG2(secretData.yQ), toG1(secretData.xyP), toG2(secretData.s_xyQ))
        expect(await secret.callStatic.respond(toG2(secretData.yQ), toG1(secretData.xyP), toG2(secretData.s_xyQ))).to.be.true
      })

      it("should reject bad responses", async () => {
        await secret.respond(toG2(secretData.yQ), toG1(secretData.xyP), toG2(secretData.sQ))
        expect(await secret.callStatic.respond(toG2(secretData.yQ), toG1(secretData.xyP), toG2(secretData.sQ))).to.be.false
      })
  })
})