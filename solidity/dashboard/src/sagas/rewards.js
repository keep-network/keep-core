import {
  take,
  takeLatest,
  call,
  put,
  fork,
  select,
  takeEvery,
} from "redux-saga/effects"
import { logError, submitButtonHelper, getContractsContext } from "./utils"
import {
  fetchtTotalDistributedRewards,
  fetchECDSAAvailableRewards,
  fetchECDSAClaimedRewards,
} from "../services/rewards"
import { sendTransaction } from "./web3"
import { isSameEthAddress } from "../utils/general.utils"
import { add } from "../utils/arithmetics.utils"

function* fetchBeaconDistributedRewards(address) {
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
  const { payload } = yield take(
    "rewards/beacon_fetch_distributed_rewards_request"
  )
  yield fork(fetchBeaconDistributedRewards, payload)
}

function* withdrawECDSARewards(action) {
  const { payload: availableRewards } = action
  const { ECDSARewardsContract } = yield getContractsContext()
  const unapproved = []

  for (const { operator, interval, withdrawable } of availableRewards) {
    try {
      yield call(sendTransaction, {
        payload: {
          contract: ECDSARewardsContract,
          methodName: "withdrawRewards",
          args: [interval, operator],
        },
      })
    } catch (error) {
      unapproved.push({ operator, interval, withdrawable })
      continue
    }
  }

  yield put({
    type: "rewards/ecdsa_update_available_rewards",
    payload: unapproved,
  })
}

export function* watchWithdrawECDSARewards() {
  yield takeLatest("rewards/ecdsa_withdraw", function* (action) {
    yield call(submitButtonHelper, withdrawECDSARewards, action)
  })
}

export function* fetchECDSARewardsData() {
  try {
    yield put({ type: "rewards/ecdsa_fetch_rewards_data_request" })
    // Fetching delegations for a current loggeed account. Operators are
    // required to fetch ecdsa available rewards for the given owner.
    yield put({ type: "staking/fetch_delegations_request" })
    const operators = yield call(getOperatorsFromStore)

    let availableRewards = yield call(fetchECDSAAvailableRewards, operators)
    const claimedRewards = yield call(fetchECDSAClaimedRewards, operators)

    availableRewards = availableRewards.filter(
      ({ operator, merkleRoot }) =>
        !claimedRewards.find(
          (lookup) =>
            isSameEthAddress(operator, lookup) &&
            merkleRoot === lookup.merkleRoot
        )
    )

    const totalAvailableAmount = availableRewards.reduce(
      (reducer, _) => add(reducer, _.amount),
      0
    )
    yield put({
      type: "rewards/ecdsa_fetch_rewards_data_success",
      payload: { totalAvailableAmount, availableRewards },
    })
  } catch (error) {
    yield* logError("rewards/ecdsa_fetch_rewards_data_failure", error)
  }
}

function* getOperatorsFromStore() {
  const getDelegationsFetchingStatus = (state) =>
    state.staking.delegationsFetchingStatus

  // Waiting for delegation reqest
  let delegationsFetchingStatus = yield select(getDelegationsFetchingStatus)
  while (delegationsFetchingStatus !== "completed") {
    yield take()
    delegationsFetchingStatus = yield select(getDelegationsFetchingStatus)
  }
  const { delegations } = yield select((state) => state.staking)
  return [...delegations].map(({ operatorAddress }) => operatorAddress)
}

export function* watchFetchECDSARewards() {
  yield takeEvery(
    "rewards/ecdsa_fetch_rewards_data_request",
    fetchECDSARewardsData
  )
}
