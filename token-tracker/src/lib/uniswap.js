import axios from "axios"

const UNISWAP_API_URL =
  "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"

export const getPairData = async (pairId) => {
  const response = await axios.post(UNISWAP_API_URL, {
    query: `query pairquery {
        pair(id: "${pairId}") {
            token0 {
                symbol,
                derivedETH,
            },
            token1 {
                symbol,
                derivedETH,
            },
            reserve0,
            reserve1,
            reserveETH,
            reserveUSD,
            token0Price,
            token1Price,
            totalSupply,
        }
    }`,
  })

  if (response.data?.errors) {
    const error = new Error(
      "Failed fetching data from Uniswap V2 subgraph API."
    )
    error.errors = response.data.errors
    throw error
  }

  return response.data.data.pair
}
