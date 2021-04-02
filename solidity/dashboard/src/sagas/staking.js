import { take, takeEvery, call, put, select, delay } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import {
  getContractsContext,
  submitButtonHelper,
  logErrorAndThrow,
  identifyTaskByAddress,
} from "./utils"
import moment from "moment"
import { sendTransaction } from "./web3"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { gt, sub } from "../utils/arithmetics.utils"
import { fromTokenUnit } from "../utils/token.utils"
import { tokensPageService } from "../services/tokens-page.service"
import {
  fetchAvailableTopUps,
  isTopUpReadyToBeCommitted,
} from "../services/top-ups.service"
import { isEmptyArray } from "../utils/array.utils"
import { showMessage } from "../actions/messages"
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
  yield takeOnlyOnce(
    "staking/fetch_delegations_request",
    (action) => action.payload.address,
    fetchDelegations
  )
}

function* fetchDelegations() {
  try {
    yield put({ type: "staking/fetch_delegations_start" })
    const data = yield call(tokensPageService.fetchTokensPageData)
    yield put({ type: "staking/fetch_delegations_success", payload: data })
  } catch (error) {
    yield* logErrorAndThrow("staking/fetch_delegations_failure", error)
  }
}

export function* watchFetchTopUpsRequest() {
  // Fetch data only once and update data based on evnets.
  yield takeOnlyOnce(
    "staking/fetch_top_ups_request",
    identifyTaskByAddress,
    fetchTopUps
  )
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
    const { delegations, undelegations, initializationPeriod } = yield select(
      (state) => state.staking
    )
    const operators = [...undelegations, ...delegations].map(
      ({ operatorAddress }) => operatorAddress
    )
    const topUps = yield call(fetchAvailableTopUps, operators)
    yield put({
      type: "staking/fetch_top_ups_success",
      payload: topUps.map((_) => ({
        ..._,
        readyToBeCommitted: isTopUpReadyToBeCommitted(_, initializationPeriod),
      })),
    })
  } catch (error) {
    yield* logErrorAndThrow("staking/fetch_top_ups_failure", error)
  }
}

export function* watchTopUpReadyToBeCommitted() {
  // Waiting for top-ups data.
  yield take("staking/fetch_top_ups_success")
  yield call(notifyTopUpReadyToBeCommitted)

  while (true) {
    yield delay(moment.duration(5, "minutes").asMilliseconds())
    yield call(notifyTopUpReadyToBeCommitted)
  }
}

function* notifyTopUpReadyToBeCommitted() {
  const topUps = yield select((state) => state.staking.topUps)
  const initializationPeriod = yield select(
    (state) => state.staking.initializationPeriod
  )

  const topUpsReadyToCommit = topUps.filter((topUp) =>
    isTopUpReadyToBeCommitted(topUp, initializationPeriod)
  )

  if (isEmptyArray(topUpsReadyToCommit)) {
    return
  }

  yield put({
    type: "staking/top_ups_ready_to_be_committed",
    payload: topUpsReadyToCommit,
  })

  const { delegations } = yield select((state) => state.staking)
  const displayedMessages = yield select((state) => state.messages)
  const liquidTopUpNotificationAlreadyDisplayed = displayedMessages.some(
    (message) =>
      message.messageType === messageType.TOP_UP_READY_TO_BE_COMMITTED &&
      !message.messageProps.grantId
  )

  // We only want to display a single notification if in a grant are multiple
  // top-ups, so we store grant ids that have already been notified. The top-ups
  // are grouped by grant in a data table.
  const notifiedGrants = new Set()
  let isFromLiquidTokens = false
  for (const { operatorAddress } of topUpsReadyToCommit) {
    if (!isFromLiquidTokens) {
      isFromLiquidTokens = delegations.some(
        (_) =>
          isSameEthAddress(_.operatorAddress, operatorAddress) && !_.isFromGrant
      )
    }
    const stake = delegations.find(
      (_) =>
        isSameEthAddress(_.operatorAddress, operatorAddress) && _.isFromGrant
    )

    if (
      stake &&
      !notifiedGrants.has(stake.grantId) &&
      !displayedMessages.some(
        (message) =>
          message.messageType === messageType.TOP_UP_READY_TO_BE_COMMITTED &&
          message.messageProps.grantId === stake.grantId
      )
    ) {
      notifiedGrants.add(stake.grantId)
      yield put(
        showMessage({
          messageType: messageType.TOP_UP_READY_TO_BE_COMMITTED,
          messageProps: {
            sticky: true,
            grantId: stake.grantId,
          },
        })
      )
    }
  }

  if (isFromLiquidTokens && !liquidTopUpNotificationAlreadyDisplayed) {
    yield put(
      showMessage({
        messageType: messageType.TOP_UP_READY_TO_BE_COMMITTED,
        messageProps: {
          sticky: true,
        },
      })
    )
  }
}
