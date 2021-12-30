import { utils, BigNumber, BigNumberish } from "ethers"

// eslint-disable-next-line import/prefer-default-export
export function calculateDkgSeed(
  relayEntry: BigNumberish,
  blockNumber: BigNumberish
): BigNumber {
  return BigNumber.from(
    utils.keccak256(
      utils.solidityPack(
        ["uint256", "uint256"],
        [BigNumber.from(relayEntry), BigNumber.from(blockNumber)]
      )
    )
  )
}
