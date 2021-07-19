import axios from "axios"
import BigNumber from "bignumber.js"

import { BaseExchange, UniswapV2Exchange } from "../exchange-api"

jest.mock("axios")

describe("Test exchange api", () => {
  describe("Test BaseExchange", () => {
    const exchange = new BaseExchange()
    exchange._getUniswapPairData = jest.fn()
    exchange._getKeepTokenPriceInUSD = jest.fn()
    exchange._getBTCPriceInUSD = jest.fn()
    const mockedPairId = 1

    test("should call function that implements fetching uniswap pair data", () => {
      exchange.getUniswapPairData(mockedPairId)

      expect(exchange._getUniswapPairData).toHaveBeenCalledWith(mockedPairId)
    })

    test("should call function that implements fetching KEEP token price", () => {
      exchange.getKeepTokenPriceInUSD()

      expect(exchange._getKeepTokenPriceInUSD).toHaveBeenCalled()
    })

    test("should call function that implements fetching BTC price", () => {
      exchange.getBTCPriceInUSD()

      expect(exchange._getBTCPriceInUSD).toHaveBeenCalled()
    })
  })

  describe("Test UniswapV2Exchange", () => {
    const exchange = new UniswapV2Exchange()
    const mockedPairId = 1

    const mockedResponse = {
      data: {
        data: {
          pair: {
            reserveUSD: 30000,
            reserveETH: 1000,
            token0: {
              derivedETH: 0.2,
            },
          },
        },
      },
    }

    test("should fetch uniswap pair data correctly", async () => {
      axios.post.mockResolvedValue(mockedResponse)

      const result = await exchange.getUniswapPairData(mockedPairId)

      const mockCalls = axios.post.mock.calls
      expect(axios.post).toHaveBeenCalled()
      expect(mockCalls[0][0]).toEqual(exchange.UNISWAP_API_URL)
      expect(
        mockCalls[0][1].query.toString().includes(mockedPairId)
      ).toBeTruthy()

      expect(result).toStrictEqual(mockedResponse.data.data.pair)
    })

    test("should fetch keep token price correctly", async () => {
      const spy = jest.spyOn(exchange, "_getTokenPriceInUSD")
      const getUniswapPairDataSpy = jest.spyOn(exchange, "_getUniswapPairData")
      axios.post.mockResolvedValue(mockedResponse)
      const pairData = mockedResponse.data.data.pair
      const expectedPrice = new BigNumber(pairData.reserveUSD)
        .div(pairData.reserveETH)
        .multipliedBy(pairData.token0.derivedETH)

      const result = await exchange.getKeepTokenPriceInUSD()

      expect(spy).toHaveBeenCalledWith(
        "0xe6f19dab7d43317344282f803f8e8d240708174a"
      )
      expect(getUniswapPairDataSpy).toHaveBeenCalledWith(
        "0xe6f19dab7d43317344282f803f8e8d240708174a"
      )
      expect(result).toEqual(expectedPrice)
    })

    test("should fetch BTC price correctly", () => {
      const spy = jest
        .spyOn(exchange, "_getTokenPriceInUSD")
        .mockResolvedValue("300")

      exchange.getBTCPriceInUSD()

      expect(spy).toHaveBeenCalledWith(
        "0xe6f19dab7d43317344282f803f8e8d240708174a"
      )
    })
  })
})
