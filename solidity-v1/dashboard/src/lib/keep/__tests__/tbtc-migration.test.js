import { TBTC } from "../../../utils/token.utils"
import TBTCV2Migration from "../tbtc-migration"

const createMockedContract = (address) => ({
  makeCall: jest.fn(),
  address: address,
  getPastEvents: jest.fn(),
})

describe("Test CoveragePoolV1 lib", () => {
  /** @type {TBTCV2Migration} */
  let tbtcV2Migration
  beforeEach(() => {
    const tbtcv1 = createMockedContract("0x9")
    const tbtcv2 = createMockedContract("0x7")
    const vendingMachine = createMockedContract("0x6")

    const web3 = {}

    tbtcV2Migration = new TBTCV2Migration(tbtcv1, tbtcv2, vendingMachine, web3)
  })

  test.each`
    tokenVersion | mockedBalance
    ${"V1"}      | ${1000}
    ${"V2"}      | ${2000}
  `(
    "should fetch the tbtc $tokenVersion balance",
    async ({ tokenVersion, mockedBalance }) => {
      const mockedResult = mockedBalance
      const address = "0x1"
      tbtcV2Migration[`tbtc${tokenVersion}`].makeCall.mockResolvedValue(
        mockedResult
      )

      const result = await tbtcV2Migration[`tbtc${tokenVersion}BalanceOf`](
        address
      )

      expect(
        tbtcV2Migration[`tbtc${tokenVersion}`].makeCall
      ).toHaveBeenCalledWith("balanceOf", address)
      expect(result).toEqual(mockedResult)
    }
  )

  it("should fetch unmint fee", async () => {
    const mockedResult = 1
    tbtcV2Migration.vendingMachine.makeCall.mockResolvedValue(mockedResult)

    const result = await tbtcV2Migration.unmintFee()

    expect(tbtcV2Migration.vendingMachine.makeCall).toHaveBeenCalledWith(
      "unmintFee"
    )
    expect(result).toEqual(mockedResult)
  })

  it("should fetch unmint fee for a given amount", async () => {
    const mockedResult = 0.1
    const amount = 100
    tbtcV2Migration.vendingMachine.makeCall.mockResolvedValue(mockedResult)

    const result = await tbtcV2Migration.unmintFeeFor(amount)

    expect(tbtcV2Migration.vendingMachine.makeCall).toHaveBeenCalledWith(
      "unmintFeeFor",
      amount
    )
    expect(result).toEqual(mockedResult)
  })

  it("should calculate unmint fee for a given amount if pass unmint fee param", async () => {
    const mockedUnmintFee = TBTC.fromTokenUnit(0.001).toString()
    const amount = TBTC.fromTokenUnit(2).toString()
    const expectedResult = TBTC.fromTokenUnit(0.002).toString()

    const result = await tbtcV2Migration.unmintFeeFor(amount, mockedUnmintFee)

    expect(tbtcV2Migration.vendingMachine.makeCall).not.toHaveBeenCalled()
    expect(result.toString()).toEqual(expectedResult.toString())
  })
})
