const { contract } = require("@openzeppelin/test-environment")
const TestModUtils = contract.fromArtifact("TestModUtils")

describe("ModUtils", () => {

  let modUtils

  before(async () => {
    modUtils = await TestModUtils.new()
  })

  it("runModExponentTest()", async () => {
    await modUtils.runModExponentTest();
  })

  it("runLegendreRangeTest()", async () => {
    await modUtils.runLegendreRangeTest();
  })

  it("runLegendreListTest()", async () => {
    await modUtils.runLegendreListTest();
  })

  it("runModSqrtOf0Test()", async () => {
    await modUtils.runModSqrtOf0Test();
  })

  it("runModSqrtMultipleOfPTest()", async () => {
    await modUtils.runModSqrtMultipleOfPTest();
  })

  it("runModSqrtAgainstListTest()", async () => {
    await modUtils.runModSqrtAgainstListTest();
  })

  it("runModSqrtAgainstNonSquaresTest()", async () => {
    await modUtils.runModSqrtAgainstNonSquaresTest();
  })

  it("runModSqrtALessThanPTest()", async () => {
    await modUtils.runModSqrtALessThanPTest();
  })

  it("runModSqrtAGreaterThanPTest()", async () => {
    await modUtils.runModSqrtAGreaterThanPTest();
  })
})
