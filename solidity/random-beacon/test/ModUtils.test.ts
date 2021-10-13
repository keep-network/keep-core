import { waffle } from "hardhat"
import { ethers } from "hardhat"
import type { TestModUtils } from "../typechain"

describe("ModUtils", () => {
  let testModUtils: TestModUtils

  const fixture = async function() {
    const TestModUtils = await ethers.getContractFactory("TestModUtils")
    const testModUtils = await TestModUtils.deploy()
    await testModUtils.deployed()

    return testModUtils
  }

  beforeEach("load test fixture", async function () {
    testModUtils = await waffle.loadFixture(fixture)
  })

  it("runModExponentTest()", async () => {
    await testModUtils.runModExponentTest()
  })

  it("runLegendreRangeTest()", async () => {
    await testModUtils.runLegendreRangeTest()
  })

  it("runLegendreListTest()", async () => {
    await testModUtils.runLegendreListTest()
  })

  it("runModSqrtOf0Test()", async () => {
    await testModUtils.runModSqrtOf0Test()
  })

  it("runModSqrtMultipleOfPTest()", async () => {
    await testModUtils.runModSqrtMultipleOfPTest()
  })

  it("runModSqrtAgainstListTest()", async () => {
    await testModUtils.runModSqrtAgainstListTest()
  })

  it("runModSqrtAgainstNonSquaresTest()", async () => {
    await testModUtils.runModSqrtAgainstNonSquaresTest()
  })

  it("runModSqrtALessThanPTest()", async () => {
    await testModUtils.runModSqrtALessThanPTest()
  })

  it("runModSqrtAGreaterThanPTest()", async () => {
    await testModUtils.runModSqrtAGreaterThanPTest()
  })
})
