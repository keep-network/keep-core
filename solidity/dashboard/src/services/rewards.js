import web3Utils from "web3-utils"
import { ContractsLoaded, CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { ECDSARewardsHelper } from "../utils/rewardsHelper"
import { add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import rewardsData from "../rewards-allocation/rewards.json"

// The merkle root is as key in the `rewards.json` data. First merkle root
// points to a first reward interval, second to a second reward interval and so on.
const merkleRoots = Object.keys(rewardsData)

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

/**
 * @typedef {Object} Claim
 * @property {number} index - Index of the account in the merkle tree.
 * @property {string} amount - The amount of KEEP reward to be claimed.
 * @property {Array<string>} proof - Array of merkle proofs.
 * @property {string} operator - Address of the operator.
 * @property {number} interval - Interval.
 * @property {string} merkleRoot - Merkle root.
 */

/**
 * Gets the available rewards based on the output from the merkle object
 * generator (github.com/keep-network/keep-ecdsa/tree/master/staker-rewards) which
 * is stored in `src/rewards-allocation/rewards.json` file
 * for the given operators between the interval 0 and the
 * current interval. Note that rewards may already be withdrawn so need to take
 * into account `RewardsClaimed` evnets from the `ECDSARewardsDistributor`.
 *
 * @param {Array<string>} operators Array of operators
 * @return {Array<Claim>} Available claims
 */
export const fetchECDSAAvailableRewards = async (operators) => {
  const toWithdrawn = []
  if (isEmptyArray(operators)) {
    return toWithdrawn
  }

  for (
    let interval = 0;
    interval < ECDSARewardsHelper.currentInterval;
    interval++
  ) {
    const merkleRoot = merkleRootOfInterval(interval)
    for (const operator of operators) {
      /**
       * Contains all necessary information to calim rewards from
       * `ECDSARewardsDistributor` contract.
       * @typedef {Object} OperatorClaim
       * @property {number} index - Index of the account in the merkle tree.
       * @property {string} amount - The amount of KEEP reward to be claimed in hex.
       * @property {Array<string>} proof - Array of merkle proofs.
       */

      /**
       * @type {OperatorClaim} operatorClaim
       */
      const operatorClaim = rewardsData[merkleRoot].claims[operator]
      if (operatorClaim) {
        toWithdrawn.push({
          ...operatorClaim,
          operator,
          merkleRoot,
          interval,
          amount: web3Utils.toBN(operatorClaim.amount).toString(),
          rewardsPeriod: ECDSARewardsHelper.periodOf(interval),
        })
      }
    }
  }

  return toWithdrawn
}

const merkleRootOfInterval = (interval) => {
  return !isEmptyArray(merkleRoots) ? merkleRoots[interval] : null
}

export const fetchECDSAClaimedRewards = async (operators) => {
  const { ECDSARewardsDistributorContract } = await ContractsLoaded
  if (isEmptyArray(operators)) {
    return []
  }

  return (
    await ECDSARewardsDistributorContract.getPastEvents("RewardsClaimed", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.ECDSARewardsDistributorContract,
      filter: { operator: operators },
    })
  ).map((event) => {
    const intervalOf = merkleRoots.indexOf(event.returnValues.merkleRoot)

    return {
      ...event.returnValues,
      rewardsPeriod: ECDSARewardsHelper.periodOf(intervalOf),
    }
  })
}
