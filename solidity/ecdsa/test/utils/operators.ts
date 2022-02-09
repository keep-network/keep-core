/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"

// eslint-disable-next-line import/no-cycle
import { constants } from "../fixtures"

import type { Address } from "hardhat-deploy/types"
import type { BigNumber } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, T } from "../../typechain"

export type OperatorID = number
export type Operator = { id: OperatorID; address: Address }

export async function registerOperators(
  walletRegistry: WalletRegistry,
  tToken: T,
  addresses: Address[],
  stakeAmount: BigNumber = constants.minimumStake
): Promise<Operator[]> {
  const operators: Operator[] = []

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

  const sortitionPool = await ethers.getContractAt(
    "SortitionPool",
    await walletRegistry.sortitionPool()
  )

  const staking = await ethers.getContractAt(
    "TokenStaking",
    await walletRegistry.staking()
  )

  for (let i = 0; i < addresses.length; i++) {
    const stakingProvider: SignerWithAddress = await ethers.getSigner(
      addresses[i]
    )

    // TODO: Use unique addresses for each role.
    // const owner: SignerWithAddress = stakingProvider
    const operator: SignerWithAddress = stakingProvider
    // const beneficiary: string = operator
    const authorizer: SignerWithAddress = stakingProvider

    await tToken.connect(deployer).mint(operator.address, stakeAmount)

    await tToken.connect(stakingProvider).approve(staking.address, stakeAmount)

    // TODO: Uncomment when integrating with the real TokenStaking contract.
    // await staking
    //   .connect(owner)
    //   .stake(stakingProvider, beneficiary, authorizer, stakeAmount)

    await staking
      .connect(authorizer)
      .increaseAuthorization(
        stakingProvider.address,
        walletRegistry.address,
        stakeAmount
      )

    await walletRegistry.connect(operator).registerOperator()

    const id = await sortitionPool.getOperatorID(operator.address)

    operators.push({ id, address: operator.address })
  }

  return operators
}
