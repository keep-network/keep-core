import { takeEvery, call } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper } from "./utils"
import { sendTransaction } from "./web3"
import { gt } from "../utils/arithmetics.utils"

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
