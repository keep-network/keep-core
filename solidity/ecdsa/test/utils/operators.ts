/* eslint-disable no-await-in-loop */

import { ethers, helpers } from "hardhat"

// eslint-disable-next-line import/no-cycle
import { params } from "../fixtures"
import { testConfig } from "../../hardhat.config"

import type { BigNumber, BigNumberish } from "ethers"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"
import type {
  WalletRegistry,
  T,
  SortitionPool,
  TokenStaking,
} from "../../typechain"

export type OperatorID = number
export type Operator = {
  id: OperatorID
  signer: SignerWithAddress
}

export async function registerOperators(
  walletRegistry: WalletRegistry,
  t: T,
  numberOfOperators = testConfig.operatorsCount,
  unnamedSignersOffset = testConfig.nonStakingAccountsCount,
  stakeAmount: BigNumber = params.minimumAuthorization
): Promise<Operator[]> {
  const operators: Operator[] = []

  const sortitionPool: SortitionPool = await ethers.getContractAt(
    "SortitionPool",
    await walletRegistry.sortitionPool()
  )

  const staking: TokenStaking = await ethers.getContractAt(
    "TokenStaking",
    await walletRegistry.staking()
  )

  const signers = (await helpers.signers.getUnnamedSigners()).slice(
    unnamedSignersOffset
  )

  // We use unique accounts for each staking role for each operator.
  if (signers.length < numberOfOperators * 5) {
    throw new Error(
      "not enough unnamed signers; update hardhat network's configuration account count"
    )
  }

  for (let i = 0; i < numberOfOperators; i++) {
    const owner: SignerWithAddress = signers[i]
    const stakingProvider: SignerWithAddress =
      signers[1 * numberOfOperators + i]
    const operator: SignerWithAddress = signers[2 * numberOfOperators + i]
    const beneficiary: SignerWithAddress = signers[3 * numberOfOperators + i]
    const authorizer: SignerWithAddress = signers[4 * numberOfOperators + i]

    await stake(
      t,
      staking,
      walletRegistry,
      owner,
      stakingProvider,
      stakeAmount,
      beneficiary,
      authorizer
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

export async function stake(
  t: T,
  staking: TokenStaking,
  randomBeacon: WalletRegistry,
  owner: SignerWithAddress,
  stakingProvider: SignerWithAddress,
  stakeAmount: BigNumberish,
  beneficiary = stakingProvider,
  authorizer = stakingProvider
): Promise<void> {
  const { deployer } = await helpers.signers.getNamedSigners()

  await t.connect(deployer).mint(owner.address, stakeAmount)
  await t.connect(owner).approve(staking.address, stakeAmount)

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
