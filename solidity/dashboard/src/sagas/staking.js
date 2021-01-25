import {
  take,
  takeEvery,
  call,
  fork,
  put,
  select,
  delay,
} from "redux-saga/effects"
import moment from "moment"
import { getContractsContext, submitButtonHelper, logError } from "./utils"
import { sendTransaction } from "./web3"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { gt, sub } from "../utils/arithmetics.utils"
import { fromTokenUnit } from "../utils/token.utils"
import { tokensPageService } from "../services/tokens-page.service"
import { fetchAvailableTopUps } from "../services/top-ups.service"
import { isEmptyArray } from "../utils/array.utils"
import { SHOW_MESSAGE, Message } from "../actions/messages"
import { isSameEthAddress } from "../utils/general.utils"
import { messageType } from "../components/Message"

function* delegateStake(action) {
  yield call(submitButtonHelper, resolveStake, action)
}

export function* watchDelegateStakeRequest() {
  yield takeEvery("staking/delegate_request", delegateStake)
}

function* resolveStake(action) {
  const { token, stakingContract } = yield getContractsContext()
  const {
    amount,
    grantId,
    beneficiaryAddress,
    operatorAddress,
    authorizerAddress,
  } = action.payload

  const tokenAmount = fromTokenUnit(amount).toString()
  const stakingContractAddress = stakingContract.options.address
  const delegationData =
    "0x" +
    Buffer.concat([
      Buffer.from(beneficiaryAddress.substr(2), "hex"),
      Buffer.from(operatorAddress.substr(2), "hex"),
      Buffer.from(authorizerAddress.substr(2), "hex"),
    ]).toString("hex")

  const data = {
    ...action.payload,
    delegationData,
    stakingContractAddress,
    amount: tokenAmount,
  }

  if (grantId) {
    yield call(stakeFromGrant, data)
  } else {
    yield call(sendTransaction, {
      payload: {
        contract: token,
        methodName: "approveAndCall",
        args: [stakingContractAddress, tokenAmount, delegationData],
      },
    })
  }
}

function* stakeFromGrant(data) {
  const {
    managedGrantContractInstance,
    isManagedGrant,
    grantId,
    amount,
    delegationData,
    stakingContractAddress,
  } = data
  const { grantContract } = yield getContractsContext()
  const amountLeft = yield call(
    stakeFirstFromEscrow,
    grantId,
    amount,
    delegationData
  )

  const defaultArgs = [stakingContractAddress, amountLeft, delegationData]

  if (gt(amountLeft, 0)) {
    yield call(sendTransaction, {
      payload: {
        contract: isManagedGrant ? managedGrantContractInstance : grantContract,
        methodName: "stake",
        args: isManagedGrant ? defaultArgs : [grantId, ...defaultArgs],
      },
    })
  }
}

function* stakeFirstFromEscrow(grantId, amount, extraData) {
  const { tokenStakingEscrow } = yield getContractsContext()

  const escrowDeposits = yield call(
    [tokenStakingEscrow, tokenStakingEscrow.getPastEvents],
    "Deposited",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
      filter: { grantId },
    }
  )

  let amountLeft = amount

  for (const deposit of escrowDeposits) {
    const {
      returnValues: { operator },
    } = deposit

    const availableAmount = yield call(
      tokenStakingEscrow.methods.availableAmount(operator).call
    )

    if (gt(amountLeft, 0) && gt(availableAmount, 0)) {
      try {
        const amountToRedelegate = gt(amountLeft, availableAmount)
          ? availableAmount
          : amountLeft

        yield call(sendTransaction, {
          payload: {
            contract: tokenStakingEscrow,
            methodName: "redelegate",
            args: [operator, amountToRedelegate, extraData],
          },
        })

        amountLeft = sub(amountLeft, amountToRedelegate)
      } catch (err) {
        continue
      }
    }
  }

  return amountLeft
}

export function* watchFetchDelegationRequest() {
  // Fetch data only once and update data based on evnets.
  yield take("staking/fetch_delegations_request")
  yield fork(fetchDelegations)
}

function* fetchDelegations() {
  try {
    yield put({ type: "staking/fetch_delegations_start" })
    const data = yield call(tokensPageService.fetchTokensPageData)
    yield put({ type: "staking/fetch_delegations_success", payload: data })
  } catch (error) {
    yield* logError("staking/fetch_delegations_failure", error)
  }
}

export function* watchFetchTopUpsRequest() {
  // Fetch data only once and update data based on evnets.
  yield take("staking/fetch_top_ups_request")
  yield fork(fetchTopUps)
}

function* fetchTopUps() {
  const getDelegationsFetchingStatus = (state) =>
    state.staking.delegationsFetchingStatus
  try {
    // We want to fetch top ups based on previously fetched delegations.
    let delegationsFetchingStatus = yield select(getDelegationsFetchingStatus)
    while (delegationsFetchingStatus !== "completed") {
      yield take()
      delegationsFetchingStatus = yield select(getDelegationsFetchingStatus)
    }

    yield put({ type: "staking/fetch_top_ups_start" })
    const { delegations, undelegations } = yield select(
      (state) => state.staking
    )
    const operators = [...undelegations, ...delegations].map(
      ({ operatorAddress }) => operatorAddress
    )
    const topUps = yield call(fetchAvailableTopUps, operators)
    yield put({ type: "staking/fetch_top_ups_success", payload: topUps })
  } catch (error) {
    yield* logError("staking/fetch_top_ups_failure", error)
  }
}

export function* watchTopUpReadyToBeCommitted() {
  // Waiting for top-ups data.
  yield take("staking/fetch_top_ups_success")
  yield call(notifyTopUpReadyToBeCommitted)

  while (true) {
    // Every 5 minutes.
    yield delay(moment.duration(5, "minutes").asMilliseconds())
    yield call(notifyTopUpReadyToBeCommitted)
  }
}

function* notifyTopUpReadyToBeCommitted() {
  const topUps = yield select((state) => state.staking.topUps)
  const initializationPeriod = yield select(
    (state) => state.staking.initializationPeriod
  )

  const topUpsReadyToCommit = topUps.filter(({ createdAt }) =>
    moment
      .unix(createdAt)
      .add(initializationPeriod, "seconds")
      .isBefore(moment())
  )

  if (isEmptyArray(topUpsReadyToCommit)) {
    return
  }

  const { delegations, undelegations } = yield select((state) => state.staking)
  let isFromGrant = false
  let isFromLiquidTokens = false
  const staking = [...delegations, undelegations]
  for (const { operatorAddress } of topUpsReadyToCommit) {
    if (!isFromLiquidTokens) {
      isFromLiquidTokens = staking.some(
        (_) =>
          isSameEthAddress(_.operatorAddress, operatorAddress) && !_.isFromGrant
      )
    } else if (!isFromGrant) {
      isFromGrant = staking.some(
        (_) =>
          isSameEthAddress(_.operatorAddress, operatorAddress) && _.isFromGrant
      )
    }
  }

  if (isFromGrant) {
    yield put({
      type: SHOW_MESSAGE,
      payload: Message.create({
        // TODO: add the new type of message which handles the `top up ready to
        // be committed` notification. The PR which will add support for
        // displaying the specific notification via `type` field is still in
        // prgreess[1]. We should sync with that changes as soon as the PR is merged
        // to the `master` branch.
        // [1]: https://github.com/keep-network/keep-core/pull/2272
        type: messageType.WALLET,
        context: "grant",
        title: "[Grant] Top Up grant Ready To Be Committed",
        sticky: true,
      }),
    })
  }
  if (isFromLiquidTokens) {
    yield put({
      type: SHOW_MESSAGE,
      payload: Message.create({
        type: messageType.WALLET,
        context: "wallet",
        title: "[Wallet] Top Up Ready To Be Committed",
        sticky: true,
      }),
    })
  }
}
