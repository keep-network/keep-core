import { take, takeEvery, fork, call, put } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper, logError } from "./utils"
import { sendTransaction } from "./web3"
import { gt } from "../utils/arithmetics.utils"
import { tokenGrantsService } from "../services/token-grants.service"

function* releaseTokensWorker(action) {
  yield call(submitButtonHelper, releaseTokens, action)
}

function* releaseTokens(action) {
  const { grantContract, tokenStakingEscrow } = yield getContractsContext()
  const {
    isManagedGrant,
    managedGrantContractInstance,
    withdrawableAmountGrantOnly,
    escrowOperatorsToWithdraw,
    id,
  } = action.payload

  if (gt(withdrawableAmountGrantOnly, 0)) {
    if (isManagedGrant) {
      yield call(sendTransaction, {
        payload: {
          contract: managedGrantContractInstance,
          methodName: "withdraw",
          args: [],
        },
      })
    } else {
      yield call(sendTransaction, {
        payload: {
          contract: grantContract,
          methodName: "withdraw",
          args: [id],
        },
      })
    }
  }

  // Withdraw from escrow
  const methodName = isManagedGrant ? "withdrawToManagedGrantee" : "withdraw"
  for (const operator of escrowOperatorsToWithdraw) {
    yield call(sendTransaction, {
      payload: {
        contract: tokenStakingEscrow,
        methodName,
        args: [operator],
      },
    })
  }
}

export function* watchReleaseTokens() {
  yield takeEvery("token-grant/release_tokens", releaseTokensWorker)
}

export function* watchFetchGrants() {
  yield take("token-grant/fetch_grants_request")
  yield fork(fetchGrants)
}

function* fetchGrants() {
  try {
    yield put({
      type: "token-grant/fetch_grants_start",
    })
    const tokenGrants = yield call([
      tokenGrantsService,
      tokenGrantsService.fetchGrants,
    ])
    yield put({
      type: "token-grant/fetch_grants_success",
      payload: tokenGrants,
    })
  } catch (err) {
    yield* logError("token-grant/fetch_grants_failure", err)
  }
}
