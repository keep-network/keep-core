import { ContractsLoaded } from "../contracts"
import { lte, gt } from "../utils/arithmetics.utils"

export class StakeOwnedStartegy {
  async stake(stakingContractAddress, amount, delegationExtraData) {
    const { stakingContract } = await ContractsLoaded

    return await stakingContract.methods
      .approveAndCall(stakingContractAddress, amount, delegationExtraData)
      .send()
  }
}

export class StakeGrantStrategy {
  constructor(tokenGrantId) {
    this.tokenGrantId = tokenGrantId
  }

  async stake(stakingContractAddress, amount, delegationExtraData) {
    const { tokenGrant } = await ContractsLoaded

    const amountLeft = await stakeFirstFromEscrow(
      this.tokenGrantId,
      amount,
      extraData
    )

    if (lte(amountLeft, 0)) {
      return
    }

    await tokenGrant.methods
      .stake(
        this.tokenGrantId,
        stakingContractAddress,
        amount,
        delegationExtraData
      )
      .send()
  }
}

export class StakeMangedGrantStrategy {
  constructor(contract, tokenGrantId) {
    this.managedGrantContract = contract
    this.tokenGrantId = tokenGrantId
  }

  async stake(stakingContractAddress, amount, delegationExtraData) {
    const amountLeft = await stakeFirstFromEscrow(
      this.tokenGrantId,
      amount,
      extraData
    )

    if (lte(amountLeft, 0)) {
      return
    }

    await this.managedGrantContract.methods
      .stake(stakingContractAddress, amountLeft, delegationExtraData)
      .send()
      .onTransactionHashCallback()
  }
}

const stakeFirstFromEscrow = async (grantId, amount, extraData) => {
  const { tokenStakingEscrow } = await ContractsLoaded

  const escrowDeposits = await tokenStakingEscrow.getPastEvents("Deposited", {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
    filter: { grantId },
  })

  let amountLeft = amount

  for (const deposit of escrowDeposits) {
    const {
      returnValues: { operator },
    } = deposit

    const availableAmount = await tokenStakingEscrow.methods
      .availableAmount(operator)
      .call()

    if (gt(amountLeft, 0) && gt(availableAmount, 0)) {
      try {
        const amountToRedelegate = gt(amountLeft, availableAmount)
          ? availableAmount
          : amountLeft
        await tokenStakingEscrow.methods
          .redelegate(operator, availableAmount, extraData)
          .send()
        amountLeft = sub(amountLeft, amountToRedelegate)
      } catch (err) {
        continue
      }
    }
  }

  return amountLeft
}

export class StakingManager {
  static async stake(data, stakingStrategy) {
    const {
      stakingContractAddress,
      amount,
      beneficiaryAddress,
      operatorAddress,
      authorizerAddress,
    } = data

    const extraData =
      "0x" +
      Buffer.concat([
        Buffer.from(beneficiaryAddress.substr(2), "hex"),
        Buffer.from(operatorAddress.substr(2), "hex"),
        Buffer.from(authorizerAddress.substr(2), "hex"),
      ]).toString("hex")

    await stakingStrategy.stake(stakingContractAddress, amount, extraData)
  }
}
