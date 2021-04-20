import BigNumber from "bignumber.js"
import { Token } from "../../utils/token.utils"

describe("Test `Token` class", () => {
  const name = "Test Token"
  const decimals = 18
  const symbol = "TT"
  const smallestPrecisionUnit = "smallestTT"
  const smallestPrecisionDecimals = 14
  const icon = null
  const decimalsToDisplay = 5
  const TestToken = new Token(
    name,
    decimals,
    symbol,
    smallestPrecisionUnit,
    smallestPrecisionDecimals,
    icon,
    decimalsToDisplay
  )

  it("should create the Token instance correctly", () => {
    expect(TestToken.name).toBe(name)
    expect(TestToken.decimals).toBe(decimals)
    expect(TestToken.symbol).toBe(symbol)
    expect(TestToken.smallestPrecisionUnit).toBe(smallestPrecisionUnit)
    expect(TestToken.smallestPrecisionDecimals).toBe(smallestPrecisionDecimals)
    expect(TestToken.icon).toBe(icon)
    expect(TestToken.decimalsToDisplay).toBe(decimalsToDisplay)
  })

  describe("displayAmount test", () => {
    it("should return 0 if an amount param was not provided", () => {
      const result = TestToken.displayAmount()

      expect(result).toBe("0")
    })

    it("should return small amount correctly", () => {
      const amount = "90000000000000" // 0,00009

      const result = TestToken.displayAmount(amount)

      expect(result).toBe("<0.0001")
    })

    it("should return the amount with separators and w/o trailing zeros", () => {
      const amountWithTrailingZeros = "100000200000000000000000" // 100 000.2
      const amountWithLongDecimals = "100000244332211550000000" // 100 000.24433221155

      const result1 = TestToken.displayAmount(amountWithTrailingZeros)
      const result2 = TestToken.displayAmount(amountWithLongDecimals)

      expect(result1).toBe("100,000.2")
      expect(result2).toBe("100,000.24433")
    })
  })

  it("should return amount with symbol", () => {
    const amount = "100000000000000000000000" // 100 000

    const result = TestToken.displayAmountWithSymbol(amount)

    expect(result).toBe(`100,000 ${TestToken.symbol}`)
  })

  describe("displayAmountWithMetricSuffix test", () => {
    it("should return amount with metric suffix", () => {
      const amount = "123456123400000000000000" // 123456,1234
      const milions = "1639000000000000000000000" // 1 639 000
      const milions2 = "2230639000000000000000000000" // 2 230 639 000

      const result = TestToken.displayAmountWithMetricSuffix(amount)
      const resultMilions = TestToken.displayAmountWithMetricSuffix(milions)
      const resultMilions2 = TestToken.displayAmountWithMetricSuffix(milions2)

      expect(result).toBe("123.45K")
      expect(resultMilions).toBe("1.63M")
      expect(resultMilions2).toBe("2,230.63M")
    })

    it("should return 0 if the amount is 0", () => {
      expect(TestToken.displayAmountWithMetricSuffix(0)).toBe("0")
    })

    it("should return 0 if the amount is null", () => {
      expect(TestToken.displayAmountWithMetricSuffix(null)).toBe("0")
    })

    it("should return 0 if the amount is undefined", () => {
      expect(TestToken.displayAmountWithMetricSuffix()).toBe("0")
      expect(TestToken.displayAmountWithMetricSuffix(undefined)).toBe("0")
    })

    it("should return a `<...` if the amount is samllest than smallestDecimalsPrercision", () => {
      const smallAmount = new BigNumber(8).times(
        new BigNumber(10).pow(TestToken.smallestPrecisionDecimals - 1) // 0.00008
      )

      expect(TestToken.displayAmountWithMetricSuffix(smallAmount)).toBe(
        "<0.0001"
      )
    })
  })

  it("should convert token amount to readable format", () => {
    const tokenAmount = "100000000000000000000000" // 100 000

    const result = TestToken.toTokenUnit(tokenAmount).toString()
    const expectedValue = new BigNumber(tokenAmount).div(
      new BigNumber(10).pow(new BigNumber(TestToken.decimals))
    )

    expect(result).toBe(expectedValue.toString())
    expect(TestToken.toTokenUnit().toString()).toBe("0")
  })

  it("should convert readable amount to min token unit", () => {
    const amount = 100000
    const expectedValue = new BigNumber(amount).times(
      new BigNumber(10).pow(TestToken.decimals)
    )

    expect(TestToken.fromTokenUnit(amount).toString()).toBe(
      expectedValue.toString()
    )
  })
})
