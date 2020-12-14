import {
  take,
  takeLatest,
  call,
  put,
  fork,
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
import { getOperatorsOfBeneficiary } from "../services/token-staking.service"

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

function* fetchECDSARewardsData(beneficiary) {
  try {
    yield put({ type: "rewards/ecdsa_fetch_rewards_data_request" })

    const opeerators = yield call(getOperatorsOfBeneficiary, beneficiary)
    let availableRewards = yield call(fetchECDSAAvailableRewards, opeerators)
    const claimedRewards = yield call(fetchECDSAClaimedRewards, opeerators)

    // Available rewards are fetched from merkle generator's output file. This
    // file doesn't take into account a rewards alredy claimed. So we need to
    // filter out claimed rewards.
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
      payload: {
        totalAvailableAmount,
        availableRewards,
        claimedRewards,
      },
    })
  } catch (error) {
    yield* logError("rewards/ecdsa_fetch_rewards_data_failure", error)
  }
}

export function* watchFetchECDSARewards() {
  yield takeEvery(
    "rewards/ecdsa_fetch_rewards_data_request",
    fetchECDSARewardsData
  )
}
