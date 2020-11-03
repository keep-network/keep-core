import { take, takeEvery, call, fork, put, select } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper, logError } from "./utils"
import { sendTransaction } from "./web3"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { gt, sub } from "../utils/arithmetics.utils"
import { fromTokenUnit } from "../utils/token.utils"
import { tokensPageService } from "../services/tokens-page.service"
import { fetchAvailableTopUps } from "../services/top-ups.service"

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
  yield take("staking/fetch_delegations_success")
  yield take("staking/fetch_top_ups_request")

  yield fork(fetchTopUps)
}

function* fetchTopUps() {
  try {
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
