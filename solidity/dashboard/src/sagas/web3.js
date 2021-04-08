import { eventChannel, END, buffers } from "redux-saga"
import { take, takeEvery, put, call } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper } from "./utils"
import {
  showCreatedMessage,
  showMessage,
  closeMessage,
} from "../actions/messages"
import { messageType } from "../components/Message"

function createTransactionEventChannel(
  contract,
  method,
  args = [],
  options = {}
) {
  const showWalletMessage = showMessage({
    messageType: messageType.WALLET,
    messageProps: {
      sticky: true,
    },
  })

  const emitter = contract.methods[method](...args).send(options)

  let txHashCache

  let showPendingActionMessage
  let showSuccessMessage
  let showErrorMessage

  return eventChannel((emit) => {
    emit(showWalletMessage)
    emitter
      .once("transactionHash", (txHash) => {
        emit(closeMessage(showWalletMessage.payload.id))
        txHashCache = txHash
        showPendingActionMessage = showCreatedMessage({
          id: txHash,
          messageType: messageType.PENDING_ACTION,
          messageProps: {
            txHash: txHash,
            withTransactionHash: true,
            sticky: true,
          },
        })
        emit(showPendingActionMessage)
      })
      .once("receipt", (receipt) => {
        let id
        if (receipt && receipt.transactionHash) {
          id = receipt.transactionHash
        } else {
          id = txHashCache
        }
        emit(closeMessage(showWalletMessage.payload.id))
        emit(closeMessage(id))
        showSuccessMessage = showMessage({
          messageType: messageType.SUCCESS,
          messageProps: {
            txHash: id,
            withTransactionHash: true,
            sticky: true,
          },
        })
        emit(showSuccessMessage)
        emit(END)
      })
      .once("error", (error, receipt) => {
        emit(closeMessage(showWalletMessage.payload.id))
        emit(closeMessage(txHashCache))
        if (error.name === "ExplorerModeSubproviderError") return
        showErrorMessage = showMessage({
          messageType: messageType.ERROR,
          messageProps: {
            content: error.message,
            sticky: true,
          },
        })
        emit(showErrorMessage)
        emit(new Error())
      })

    return () => {}
  }, buffers.expanding())
}

export function createSubcribeToContractEventChannel(contract, eventName) {
  const contractHasEvent = contract.options.jsonInterface.find(
    (entry) => entry.type === "event" && entry.name === eventName
  )
  if (!contractHasEvent) {
    return eventChannel((emit) => {
      emit(END)

      return () => {}
    }, buffers.expanding())
  }

  const eventEmitter = contract.events[eventName]()
  let eventTxCache = null

  return eventChannel((emit) => {
    eventEmitter
      .on("data", (event) => {
        if (eventTxCache !== event.transactionHash) {
          eventTxCache = event.transactionHash
          emit(event)
        }
      })
      .on("error", () => {
        emit(new Error())
        emit(END)
      })

    return () => {
      eventEmitter.unsubscribe()
    }
  })
}

export function* sendTransaction(action) {
  const { contract, methodName, args, options } = action.payload

  const transactionEventChannel = createTransactionEventChannel(
    contract,
    methodName,
    args,
    options
  )

  try {
    while (true) {
      const event = yield take(transactionEventChannel)
      yield put(event)
    }
  } catch (error) {
    throw error
  } finally {
    transactionEventChannel.close()
  }
}

export function* watchSendTransactionRequest() {
  yield takeEvery("web3/send_transaction", function* (action) {
    const { contractName, methodName, args, options } = action.payload
    const contracts = yield getContractsContext()

    const sendTransactionPayload = {
      contract: contracts[contractName],
      methodName,
      args,
      options,
    }

    yield call(submitButtonHelper, sendTransaction, {
      payload: sendTransactionPayload,
      meta: action.meta,
    })
  })
}
