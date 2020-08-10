import { takeEvery, call } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper } from "./utils"
import { sendTransaction } from "./web3"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { gt, sub } from "../utils/arithmetics.utils"

function* delegateStake(action) {
  yield call(submitButtonHelper, resolveStake, action)
}

export function* watchDelegateStakeRequest() {
  yield takeEvery("staking/delegate_request", delegateStake)
}

function* resolveStake(data) {
  const { token } = yield getContractsContext()
  const { stakingContractAddress, delegationData, amount, grantId } = data
  if (grantId) {
    yield call(stakeFromGrant, data)
  } else {
    yield call(sendTransaction, {
      payload: {
        contract: token,
        method: "approveAndCall",
        args: [stakingContractAddress, amount, delegationData],
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
            args: [operator, availableAmount, extraData],
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
