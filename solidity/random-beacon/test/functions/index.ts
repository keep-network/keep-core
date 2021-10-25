import { BigNumber } from "ethers"

// eslint-disable-next-line import/prefer-default-export
export function to1e18(n) {
  const decimalMultiplier = BigNumber.from(10).pow(18)
  return BigNumber.from(n).mul(decimalMultiplier)
}

export const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
