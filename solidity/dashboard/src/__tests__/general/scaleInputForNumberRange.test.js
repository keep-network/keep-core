import { scaleInputForNumberRange } from "../../utils/general.utils"
import BigNumber from "bignumber.js"

describe("Test scaleInputForNumberRange function", () => {
  it("should return proper values", () => {
    expect(scaleInputForNumberRange(2, 1, 3, 1, 5).toNumber()).toBe(3)
    expect(scaleInputForNumberRange(2, 1, 3, 1, 10).toNumber()).toBe(5.5)
    expect(scaleInputForNumberRange(1, 1, 50, 24, 90).toNumber()).toBe(24)
    expect(scaleInputForNumberRange(3, 1, 2, 5, 10).toNumber()).toBe(15)
    expect(scaleInputForNumberRange(5, 10, 15, 6, 9).toNumber()).toBe(3)
  })

  it("should work for numbers", () => {
    expect(scaleInputForNumberRange(2, 1, 3, 1, 5).toNumber()).toBe(3)
  })

  it("should work for strings", () => {
    expect(scaleInputForNumberRange("2", "1", "3", "1", "5").toNumber()).toBe(3)
  })

  it("should work for BigNumbers", () => {
    expect(
      scaleInputForNumberRange(
        new BigNumber(2),
        new BigNumber(1),
        new BigNumber(3),
        new BigNumber(1),
        new BigNumber(5)
      ).toNumber()
    ).toBe(3)
  })

  it("should return BigNumber", () => {
    expect(scaleInputForNumberRange(2, 1, 3, 1, 5)).toBeInstanceOf(BigNumber)
  })

  it("should return NaN from function call without all arguments passed", () => {
    expect(isNaN(scaleInputForNumberRange(2, 2))).toBe(true)
  })
})
