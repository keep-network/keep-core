const { contract } = require("@openzeppelin/test-environment")
const TestArrayUtils = contract.fromArtifact("TestArrayUtils")

describe("TestArrayUtils", () => {

  let arrayUtils

  before(async () => {
    arrayUtils = await TestArrayUtils.new()
  })

  it("runCanHandleEmptyArrayTest()", async () => {
    await arrayUtils.runCanHandleEmptyArrayTest();
  })

  it("runCanRemoveAddressFromSingleElementArrayTest()", async () => {
    await arrayUtils.runCanRemoveAddressFromSingleElementArrayTest();
  })

  it("runCanRemoveIdenticalAddressesTest()", async () => {
    await arrayUtils.runCanRemoveIdenticalAddressesTest();
  })

  it("runCanRemoveAddressTest()", async () => {
    await arrayUtils.runCanRemoveAddressTest();
  })

  it("runCanHandleEmptyValueArrayTest()", async () => {
    await arrayUtils.runCanHandleEmptyValueArrayTest();
  })

  it("runCanRemoveValueFromSingleElementArrayTest()", async () => {
    await arrayUtils.runCanRemoveValueFromSingleElementArrayTest();
  })

  it("runCanRemoveIdenticalValuesTest()", async () => {
    await arrayUtils.runCanRemoveIdenticalValuesTest();
  })

  it("runCanRemoveValueTest()", async () => {
    await arrayUtils.runCanRemoveValueTest();
  })
})
