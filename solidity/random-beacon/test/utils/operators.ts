/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"

import { params } from "../fixtures"

import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type { Address } from "hardhat-deploy/types"
import type { BigNumber, BigNumberish } from "ethers"
import type {
  RandomBeacon,
  RandomBeaconStub,
  T,
  TokenStaking,
} from "../../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

export type OperatorID = number
export type Operator = { id: OperatorID; signer: SignerWithAddress }

export async function registerOperators(
  randomBeacon: RandomBeacon,
  t: T,
  addresses: Address[],
  stakeAmount: BigNumber = params.minimumAuthorization
): Promise<Operator[]> {
  const operators: Operator[] = []

  const sortitionPool = await ethers.getContractAt(
    "SortitionPool",
    await randomBeacon.sortitionPool()
  )

  const staking: TokenStaking = await ethers.getContractAt(
    "TokenStaking",
    await randomBeacon.staking()
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

    await stake(
      t,
      staking,
      randomBeacon,
      owner,
      stakingProvider,
      stakeAmount,
      operator,
      beneficiary,
      authorizer
    )

    await randomBeacon
      .connect(stakingProvider)
      .registerOperator(operator.address)

    await randomBeacon.connect(operator).joinSortitionPool()

    const id = await sortitionPool.getOperatorID(operator.address)

    operators.push({ id, signer: operator })
  }

  return operators
}

export async function stake(
  t: T,
  staking: TokenStaking,
  randomBeacon: RandomBeacon | RandomBeaconStub,
  owner: SignerWithAddress,
  stakingProvider: SignerWithAddress,
  stakeAmount: BigNumberish,
  operator = stakingProvider,
  beneficiary = stakingProvider,
  authorizer = stakingProvider
): Promise<void> {
  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

  await t.connect(deployer).mint(operator.address, stakeAmount)
  await t.connect(stakingProvider).approve(staking.address, stakeAmount)

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
      randomBeacon.address,
      stakeAmount
    )
}
