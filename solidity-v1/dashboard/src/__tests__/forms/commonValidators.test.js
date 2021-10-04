import { isNumeric } from "../../forms/common-validators"

describe("Test isNumeric function", () => {
  it("should return proper values", () => {
    expect(isNumeric(2)).toBe(true)
    expect(isNumeric("2")).toBe(true)
    expect(isNumeric(2.254)).toBe(true)
    expect(isNumeric(2.12345678901234567890123456789)).toBe(true)
    expect(isNumeric("2,432,432")).toBe(false)
    expect(isNumeric("123.123.12")).toBe(false)
    expect(isNumeric(".23")).toBe(true)
    expect(isNumeric(-1)).toBe(true)
    expect(isNumeric(0)).toBe(true)
    expect(isNumeric(" ")).toBe(false)
    expect(isNumeric("\t\t")).toBe(false)
    expect(isNumeric("\n\r")).toBe(false)
  })
})
