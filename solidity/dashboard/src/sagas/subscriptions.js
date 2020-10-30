import { fork, take, call, put } from "redux-saga/effects"
import { createSubcribeToContractEventChannel } from "./web3"
import { getContractsContext, getWeb3Context } from "./utils"
import { add, sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"

export function* subscribeToKeepTokenTransferEvent() {
  yield take("keep-token/balance_request_success")
  yield fork(observeKeepTokenTransfer)
}

function* observeKeepTokenTransfer() {
  const { token } = yield getContractsContext()
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    token,
    "Transfer"
  )

  // Observe and dispatch an action that updates keep token balance.
  while (true) {
    try {
      const {
        returnValues: { from, to, value },
      } = yield take(contractEventCahnnel)
      console.log("data", from, to, value)

      let arithmeticOpration
      if (isSameEthAddress(defaultAccount, from)) {
        arithmeticOpration = sub
      } else if (isSameEthAddress(defaultAccount, to)) {
        arithmeticOpration = add
      }
      if (arithmeticOpration) {
        yield put({
          type: "keep-token/transfered",
          payload: { value, arithmeticOpration },
        })
      }
    } catch (error) {
      console.error(`Failed subscribing to Transfer event`, error)
      contractEventCahnnel.close()
    }
  }
}
