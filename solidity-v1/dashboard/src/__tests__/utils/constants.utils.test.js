import { renderDynamicConstant } from "../../utils/constants.utils"

describe("Test renderDynamicConstant function", () => {
  const exampleString =
    "My favourite fruits are ${fruit1}, ${fruit2} and ${fruit3}. ${This should not be replaced"
  const expectedResult =
    "My favourite fruits are apple, banana and strawberry. ${This should not be replaced"

  it("Should properly replace arguments inside string", () => {
    const result = renderDynamicConstant(
      exampleString,
      "apple",
      "banana",
      "strawberry"
    )

    expect(result).toBe(expectedResult)
  })

  it("Should properly replace arguments inside string even with too many arguments given", () => {
    const result = renderDynamicConstant(
      exampleString,
      "apple",
      "banana",
      "strawberry",
      "pineapple",
      "orange"
    )

    expect(result).toBe(expectedResult)
  })

  it("Should throw an error when too few arguments are given", () => {
    const errorMsg = "Too few arguments were given!"

    expect(() => renderDynamicConstant(exampleString, "apple")).toThrow(
      new Error(errorMsg)
    )
  })
})
