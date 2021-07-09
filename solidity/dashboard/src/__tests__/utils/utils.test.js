import { displayPercentageValue, isString } from "../../utils/general.utils"

describe("Test the `displayPercentageValue`", () => {
  test("should display correctly if the value is bigger than `max` param", () => {
    const percentageValue = 150
    const min = 0.01
    const max = 149

    const result = displayPercentageValue(percentageValue, false, min, max)

    expect(result).toEqual(`>${max}%`)
  })

  test("should display correctly if the value is less than `min` param", () => {
    const percentageValue = 0.00123
    const min = 0.01

    const result = displayPercentageValue(percentageValue, true, min)

    expect(result).toEqual(`<${min}%`)
  })

  test("should display correctly if the value is between `min` and `max`", () => {
    const percentageValue = 100.47

    const result = displayPercentageValue(percentageValue)

    expect(result).toEqual(`${percentageValue}%`)
  })
})

describe("Test `isString`", () => {
  test("should return true if the value is primitive string", () => {
    expect(isString("string")).toBeTruthy()
  })

  test("should return true if the value is String object", () => {
    // eslint-disable-next-line no-new-wrappers
    expect(isString(new String("string"))).toBeTruthy()
  })

  test("should return false if the value is not a primitive string or String object", () => {
    expect(isString(1)).toBeFalsy()
    expect(isString(false)).toBeFalsy()
  })
})
