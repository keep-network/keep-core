import { takeLatest, call, put } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import web3Utils from "web3-utils"
import {
  logError,
  submitButtonHelper,
  getContractsContext,
  identifyTaskByAddress,
} from "./utils"
import {
  fetchtTotalDistributedRewards,
  fetchECDSAAvailableRewards,
  fetchECDSAClaimedRewards,
} from "../services/rewards"
import { sendTransaction } from "./web3"
import { isSameEthAddress } from "../utils/general.utils"
import { add } from "../utils/arithmetics.utils"
import {
  getOperatorsOfBeneficiary,
  getOperatorsOfOwner,
  getOperatorsOfGrantee,
  getOperatorsOfManagedGrantee,
  getOperatorsOfCopiedDelegations,
} from "../services/token-staking.service"
import { REWARD_STATUS } from "../constants/constants"

function* fetchBeaconDistributedRewards(action) {
  const {
    payload: { address },
  } = action
  try {
    yield put({ type: "rewards/beacon_fetch_distributed_rewards_start" })
    const balance = yield call(
      fetchtTotalDistributedRewards,
      address,
      "beaconRewardsContract"
    )
    yield put({
      type: "rewards/beacon_fetch_distributed_rewards_success",
      payload: balance,
    })
  } catch (error) {
    yield* logError("rewards/beacon_fetch_distributed_rewards_failure", error)
  }
}

export function* watchFetchBeaconDistributedRewards() {
  yield takeOnlyOnce(
    "rewards/beacon_fetch_distributed_rewards_request",
    identifyTaskByAddress,
    fetchBeaconDistributedRewards
  )
}

function* withdrawECDSARewards(action) {
  const { payload: availableRewards } = action
  const { ECDSARewardsDistributorContract } = yield getContractsContext()

  for (const {
    merkleRoot,
    index,
    operator,
    amount,
    proof,
  } of availableRewards) {
    try {
      yield call(sendTransaction, {
        payload: {
          contract: ECDSARewardsDistributorContract,
          methodName: "claim",
          args: [merkleRoot, index, operator, amount, proof],
        },
      })
    } catch (error) {
      yield* logError("rewards/ecdsa_withdraw_failure", error)
      throw error
    }
  }
}

export function* watchWithdrawECDSARewards() {
  yield takeLatest("rewards/ecdsa_withdraw", function* (action) {
    yield call(submitButtonHelper, withdrawECDSARewards, action)
  })
}

function* fetchECDSARewardsData(action) {
  const { address } = action.payload
  try {
    yield put({ type: "rewards/ecdsa_fetch_rewards_data_start" })

    // Beneficiary operators.
    const beneficiaryOperators = yield call(getOperatorsOfBeneficiary, address)
    // Owner operators.
    const ownerOperators = yield call(getOperatorsOfOwner, address)
    // Grantee operators.
    const { allOperators: grenteeOperators } = yield call(
      getOperatorsOfGrantee,
      address
    )
    // Managed grantee operators.
    const { allOperators: mangedGranteeOperators } = yield call(
      getOperatorsOfManagedGrantee,
      address
    )
    // Get operators of copied delegations where an owner of the old delegations
    // is `address`.
    const copiedOperators = yield call(getOperatorsOfCopiedDelegations, address)

    const operators = Array.from(
      // The same address can be used as beneficiary and owner. So addresses in
      // array may be repeated. Let's convert to a Set to make sure the given
      // address is in array only once.
      new Set(
        beneficiaryOperators
          .concat(ownerOperators)
          .concat(grenteeOperators)
          .concat(mangedGranteeOperators)
          .concat(copiedOperators)
          .map((address) => web3Utils.toChecksumAddress(address))
      ).add(web3Utils.toChecksumAddress(address)) // Operator can also view own rewards.
    )

    let availableRewards = yield call(fetchECDSAAvailableRewards, operators)
    const claimedRewards = yield call(fetchECDSAClaimedRewards, operators)

    // Available rewards are fetched from merkle generator's output file. This
    // file doesn't take into account a rewards alredy claimed. So we need to
    // filter out claimed rewards.
    availableRewards = availableRewards.filter(
      ({ operator, merkleRoot }) =>
        !claimedRewards.find(
          (lookup) =>
            isSameEthAddress(operator, lookup.operator) &&
            merkleRoot === lookup.merkleRoot
        )
    )

    const totalAvailableAmount = availableRewards.reduce(
      (reducer, _) => add(reducer, _.amount),
      0
    )

    const rewardsHistory = availableRewards
      .map((reward) => ({
        ...reward,
        status: REWARD_STATUS.AVAILABLE,
        id: `${reward.operator}-${reward.merkleRoot}`,
      }))
      .concat(
        claimedRewards.map((reward) => ({
          ...reward,
          status: REWARD_STATUS.WITHDRAWN,
          id: `${reward.operator}-${reward.merkleRoot}`,
        }))
      )

    yield put({
      type: "rewards/ecdsa_fetch_rewards_data_success",
      payload: {
        totalAvailableAmount,
        availableRewards,
        rewardsHistory,
      },
    })
  } catch (error) {
    yield* logError("rewards/ecdsa_fetch_rewards_data_failure", error)
  }
}

export function* watchFetchECDSARewards() {
  yield takeOnlyOnce(
    "rewards/ecdsa_fetch_rewards_data_request",
    identifyTaskByAddress,
    fetchECDSARewardsData
  )
}
