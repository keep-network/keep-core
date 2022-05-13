// eslint-disable-next-line import/no-extraneous-dependencies
import { deployments, ethers } from "hardhat"

import type { BaseContract } from "ethers"

// TODO: Move these utils to hardhat-helpers plugin.

// eslint-disable-next-line import/prefer-default-export
export async function getContract<T extends BaseContract>(
  deploymentName: string,
  contractName: string = deploymentName
): Promise<T> {
  return (await ethers.getContractAt(
    contractName,
    (
      await deployments.get(deploymentName)
    ).address
  )) as T
}
