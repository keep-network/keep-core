import { eventChannel, END, buffers } from "redux-saga"
import { take, takeEvery, put, call } from "redux-saga/effects"
import { getContractsContext, submitButtonHelper } from "./utils"
import {
  Message,
  showCreatedMessage,
  showMessage,
  closeMessage,
} from "../actions/messages"
import { messageType } from "../components/Message"

function createTransactionEventChannel(contract, method, args = [], options) {
  const infoMessage = Message.create({
    title: "Waiting for the transaction confirmation...",
    type: messageType.INFO,
    sticky: true,
  })

  const emitter = contract.methods[method](...args).send(options)

  let txHashCache

  return eventChannel((emit) => {
    emit(showCreatedMessage(infoMessage))
    emitter
      .once("transactionHash", (txHash) => {
        emit(closeMessage(infoMessage.id))
        txHashCache = txHash
        emit(
          showCreatedMessage({
            id: txHash,
            content: txHash,
            sticky: true,
            type: messageType.PENDING_ACTION,
            withTransactionHash: true,
          })
        )
      })
      .once("receipt", (receipt) => {
        let id
        if (receipt && receipt.transactionHash) {
          id = receipt.transactionHash
        } else {
          id = txHashCache
        }
        emit(closeMessage(infoMessage.id))
        emit(closeMessage(id))
        emit(
          showMessage({
            title: "Success!",
            content: id,
            sticky: true,
            type: messageType.SUCCESS,
            withTransactionHash: true,
          })
        )
        emit(END)
      })
      .once("error", (error, receipt) => {
        emit(closeMessage(infoMessage.id))
        emit(closeMessage(txHashCache))
        emit(
          showMessage({
            title: "Error",
            content: error.message,
            type: messageType.ERROR,
            sticky: true,
          })
        )
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
