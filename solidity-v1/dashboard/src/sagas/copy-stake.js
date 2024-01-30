import { takeLatest, call, put, delay } from "redux-saga/effects"
import {
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS,
  FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE,
  INCREMENT_STEP,
} from "../actions"
import { fetchOldDelegations } from "../services/staking-port-backer.service"
import { getContractsContext } from "./utils"
import { sendTransaction } from "./web3"
import { getContractDeploymentBlockNumber } from "../contracts"
import { showMessage, showCreatedMessage } from "../actions/messages"
import { isEmptyArray } from "../utils/array.utils"
import { messageType } from "../components/Message"
import { STAKING_PORT_BACKER_CONTRACT_NAME } from "../constants/constants"

function* fetchOldStakingDelegations() {
  try {
    yield delay(500)
    const { undelegationPeriod, intializationPeriod, delegations } = yield call(
      fetchOldDelegations
    )
    yield put({
      type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_SUCCESS,
      payload: { delegations, undelegationPeriod, intializationPeriod },
    })
  } catch (error) {
    yield put({
      type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_FAILURE,
      payload: error,
    })
  }

  // yield fork(observeEvents)
}

export function* watchFetchOldStakingContract() {
  yield takeLatest(
    FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST,
    fetchOldStakingDelegations
  )
}

function* copyStake(action) {
  const { operatorAddress, isUndelegating } = action.payload

  try {
    // At first call `copyStake` from `StakingPortBacker` contract.
    yield call(safeCopyStake, operatorAddress)
    // Next call undelegation from old staking contract.
    if (!isUndelegating) {
      yield call(undelegateFromOldContract, action)
    }
    yield put({ type: "copy-stake/copy-stake_success" })
    yield put({ type: INCREMENT_STEP })
  } catch (error) {
    yield put(
      showMessage({
        messageType: messageType.ERROR,
        messageProps: {
          content: error.message,
          sticky: true,
        },
      })
    )
    yield put({ type: "copy-stake/copy-stake_failure", payload: error })
  }
}

function* safeCopyStake(operator) {
  const { stakingPortBackerContract } = yield getContractsContext()

  const events = yield call(
    [stakingPortBackerContract, stakingPortBackerContract.getPastEvents],
    "StakeCopied",
    {
      fromBlock: yield call(
        getContractDeploymentBlockNumber,
        STAKING_PORT_BACKER_CONTRACT_NAME
      ),
      filter: { operator },
    }
  )

  if (isEmptyArray(events)) {
    yield call(sendTransaction, {
      payload: {
        contract: stakingPortBackerContract,
        methodName: "copyStake",
        args: [operator],
      },
    })
  } else {
    const txHash = events[0].transactionHash
    yield put(
      showCreatedMessage({
        id: txHash,
        messageType: messageType.DELEGATION_ALREADY_COPIED,
        messageProps: {
          content: txHash,
          withTransactionHash: true,
          sticky: true,
        },
      })
    )
  }
}

export function* watchCopyStakeRequest() {
  yield takeLatest("copy-stake/copy-stake_request", copyStake)
}

function* undelegateFromOldContract(action) {
  const delegation = action.payload
  const {
    operatorAddress,
    isFromGrant,
    isManagedGrant,
    managedGrantContractInstance,
  } = delegation

  const { oldTokenStakingContract, grantContract } = yield getContractsContext()
  let contractInstance

  if (isManagedGrant) {
    contractInstance = managedGrantContractInstance
  } else if (isFromGrant) {
    contractInstance = grantContract
  } else {
    contractInstance = oldTokenStakingContract
  }

  yield call(sendTransaction, {
    payload: {
      contract: contractInstance,
      methodName: "undelegate",
      args: [operatorAddress],
    },
  })
}

function* undelegateFromOldContractWorker(action) {
  try {
    yield call(undelegateFromOldContract, action)
    yield put({ type: "copy-stake/undelegation_success" })
    yield put({ type: INCREMENT_STEP })
  } catch (error) {
    yield put({ type: "copy-stake/undelegation_failure", payload: error })
  }
}

function* recoverFromOldStakingContract(action) {
  const delegation = action.payload
  const {
    operatorAddress,
    isFromGrant,
    isManagedGrant,
    managedGrantContractInstance,
  } = delegation

  const { oldTokenStakingContract, grantContract } = yield getContractsContext()
  let contractInstance

  if (isManagedGrant) {
    contractInstance = managedGrantContractInstance
  } else if (isFromGrant) {
    contractInstance = grantContract
  } else {
    contractInstance = oldTokenStakingContract
  }

  try {
    yield call(sendTransaction, {
      payload: {
        contract: contractInstance,
        methodName: "recoverStake",
        args: [operatorAddress],
      },
    })
    yield put({ type: "copy-stake/recover_success" })
    yield put({ type: INCREMENT_STEP })
  } catch (error) {
    yield put({ type: "copy-stake/recover_failure", payload: error })
  }
}

export function* watchUndelegateOldStakeRequest() {
  yield takeLatest(
    "copy-stake/undelegate_request",
    undelegateFromOldContractWorker
  )
}

export function* watchRecovereOldStakeRequest() {
  yield takeLatest("copy-stake/recover_request", recoverFromOldStakingContract)
}

// function* observeEvents() {
//   const { oldTokenStakingContract, stakingPortBackerContract } = yield call(
//     getContractsContext
//   )

//   yield fork(removeOldDelegationWatcher, oldTokenStakingContract, "Undelegated")
//   yield fork(
//     removeOldDelegationWatcher,
//     stakingPortBackerContract,
//     "StakeCopied"
//   )
//   yield fork(removeOldDelegationWatcher, oldTokenStakingContract, "Recovered")
// }

// function* removeOldDelegationWatcher(contract, eventName) {
//   // Create subscription channel.
//   const contractEventCahnnel = yield call(
//     createSubcribeToContractEventChannel,
//     contract,
//     eventName
//   )

//   // Observe and dispatch an action that updates copy-stake reducer.
//   while (true) {
//     try {
//       const {
//         returnValues: { operator },
//       } = yield take(contractEventCahnnel)
//       yield put({ type: "copy-stake/remove_old_delegation", payload: operator })
//     } catch (error) {
//       console.error(`Failed subscribing to ${eventName} event`, error)
//       contractEventCahnnel.close()
//     }
//   }
// }
