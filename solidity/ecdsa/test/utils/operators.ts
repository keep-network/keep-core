/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"

// eslint-disable-next-line import/no-cycle
import { params } from "../fixtures"

import type { Address } from "hardhat-deploy/types"
import type { BigNumber } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { WalletRegistry, T } from "../../typechain"

export type OperatorID = number
export type Operator = {
  id: OperatorID
  signer: SignerWithAddress
}

export async function registerOperators(
  walletRegistry: WalletRegistry,
  tToken: T,
  addresses: Address[],
  stakeAmount: BigNumber = params.minimumAuthorization
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
    const owner: SignerWithAddress = stakingProvider
    const operator: SignerWithAddress = stakingProvider
    const beneficiary: SignerWithAddress = stakingProvider
    const authorizer: SignerWithAddress = stakingProvider

    await tToken.connect(deployer).mint(operator.address, stakeAmount)

    await tToken.connect(stakingProvider).approve(staking.address, stakeAmount)

    await staking
      .connect(owner)
      .stake(
        stakingProvider.address,
        beneficiary.address,
        authorizer.address,
        stakeAmount
      )

    await staking
      .connect(authorizer)
      .increaseAuthorization(
        stakingProvider.address,
        walletRegistry.address,
        stakeAmount
      )

    await walletRegistry
      .connect(stakingProvider)
      .registerOperator(operator.address)

    await walletRegistry.connect(operator).joinSortitionPool()

    const id = await sortitionPool.getOperatorID(operator.address)

    operators.push({ id, signer: operator })
  }

  return operators
}
