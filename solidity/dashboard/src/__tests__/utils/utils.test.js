import { displayPercentageValue } from "../../utils/general.utils"

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
