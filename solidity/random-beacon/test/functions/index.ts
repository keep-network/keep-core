import { BigNumber } from "ethers";

export function to1e18(n) {
  const decimalMultiplier = BigNumber.from(10).pow(18)
  return BigNumber.from(n).mul(decimalMultiplier)
}