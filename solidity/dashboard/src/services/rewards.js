import { ContractsLoaded, CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { getOperatorsOfBeneficiary } from "./token-staking.service"
import { ECDSARewardsHelper } from "../utils/rewardsHelper"
import { add, gt } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"

export const fetchtTotalDistributedRewards = async (
  beneficiary,
  contractName
) => {
  const contracts = await ContractsLoaded
  const rewardsContract = contracts[contractName]
  const tokenContract = contracts.token

  // Filter `Transfer` events at `KeepToken` contract by fields:
  // * `to`- as a beneficiary address,
  // * `from`- as a Rewards contract(eg. `BeaconRewards`, `ECDSARewards`) contract address.

  // In ^ that case we are sure that transferred KEEPs were from rewards contract.
  return (
    await tokenContract.getPastEvents("Transfer", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.token,
      filter: { from: rewardsContract.options.address, to: beneficiary },
    })
  ).reduce((reducer, event) => add(reducer, event.returnValues.value), 0)
}

export const fetchECDSAAvailableRewards = async (beneficiary) => {
  const { ECDSARewardsContract } = await ContractsLoaded

  const operators = await getOperatorsOfBeneficiary(beneficiary)

  let sum = 0
  const toWithdrawn = []
  if (!isEmptyArray(operators)) {
    // If beneficiary has multiple operators, call `getWithdrawableRewards`
    // function one time since `ECDSARewards` contract stores rewards per
    // beneficiary not per operator.
    const operator = operators[0]
    for (
      let interval = 0;
      interval <= ECDSARewardsHelper.currentInterval;
      interval++
    ) {
      const withdrawable = await ECDSARewardsContract.methods
        .getWithdrawableRewards(interval, operator)
        .call()

      if (gt(withdrawable, 0)) {
        sum = add(sum, withdrawable)
        toWithdrawn.push({ operator, interval, withdrawable })
      }
    }
  }

  return { totalAvailableRewards: sum, toWithdrawn }
}
