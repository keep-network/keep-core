const { expectRevert } = require("@openzeppelin/test-helpers")
const { contract } = require("@openzeppelin/test-environment")
const ArrayUtilsStub = contract.fromArtifact("ArrayUtilsStub")

describe("ArrayUtils", () => {

  let arrayUtils

  before(async () => {
    arrayUtils = await ArrayUtilsStub.new()
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
