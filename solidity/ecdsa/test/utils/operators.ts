/* eslint-disable no-await-in-loop */

import { ethers } from "hardhat"
import type { Address } from "hardhat-deploy/types"
import type { BigNumber } from "ethers"
import { constants } from "../fixtures"
import type { WalletFactory } from "../../typechain"
import type { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers"

export type OperatorID = number
export type Operator = { id: OperatorID; address: Address }

export async function registerOperators(
  walletFactory: WalletFactory,
  addresses: Address[],
  stakeAmount: BigNumber = constants.minimumStake
): Promise<Operator[]> {
  const operators: Operator[] = []

  const deployer: SignerWithAddress = await ethers.getNamedSigner("deployer")

  const sortitionPool = await ethers.getContractAt(
    "SortitionPool",
    await walletFactory.sortitionPool()
  )

  const tToken = await ethers.getContractAt("T", await walletFactory.tToken())

  const staking = await ethers.getContractAt(
    "TokenStaking",
    await walletFactory.staking()
  )

  for (let i = 0; i < addresses.length; i++) {
    const operator: string = addresses[i]
    // const beneficiary: string = operator
    const authorizer: string = operator

    await tToken.connect(deployer).mint(operator, stakeAmount)

    await tToken
      .connect(await ethers.getSigner(operator))
      .approve(staking.address, stakeAmount)

    // TODO: Uncomment when integrating with the real TokenStaking contract.
    // await staking
    //   .connect(await ethers.getSigner(operator))
    //   .stake(operator, beneficiary, authorizer, stakeAmount)

    await staking
      .connect(await ethers.getSigner(authorizer))
      .increaseAuthorization(operator, walletFactory.address, stakeAmount)

    await walletFactory
      .connect(await ethers.getSigner(operator))
      .registerOperator()

    const id = await sortitionPool.getOperatorID(operator)

    operators.push({ id, address: operator })
  }

  return operators
}
