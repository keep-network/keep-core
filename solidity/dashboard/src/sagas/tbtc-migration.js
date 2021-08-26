import {
  put,
  call,
  take,
  actionChannel,
  select,
  takeEvery,
} from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import selectors from "./selectors"
import { Keep } from "../contracts"
import { tbtcV2Migration } from "../actions"
import {
  identifyTaskByAddress,
  logErrorAndThrow,
  submitButtonHelper,
} from "./utils"
import { isSameEthAddress } from "../utils/general.utils"
import { sendTransaction, approveAndTransferToken } from "./web3"
import { add } from "../utils/arithmetics.utils"

function* fetchData(action) {
  try {
    const { address } = action.payload
    yield put(tbtcV2Migration.fetchDataStart())

    const tbtcV1Balance = yield call(
      Keep.tBTCV2Migration.tbtcV1BalanceOf,
      address
    )
    const tbtcV2Balance = yield call(
      Keep.tBTCV2Migration.tbtcV2BalanceOf,
      address
    )
    const unmintFee = yield call(Keep.tBTCV2Migration.unmintFee)

    yield put(
      tbtcV2Migration.fetchDataSuccess({
        tbtcV1Balance,
        tbtcV2Balance,
        unmintFee,
      })
    )
  } catch (error) {
    yield* logErrorAndThrow(
      tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_ERROR,
      error
    )
  }
}

export function* watchFetchData() {
  yield takeOnlyOnce(
    tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_REQUEST,
    identifyTaskByAddress,
    fetchData
  )
}

export function* subscribeToTBTCV2MintedEvent() {
  const requestChan = yield actionChannel(
    tbtcV2Migration.TBTCV2_TOKEN_MINTED_EVENT_EMITTED
  )

  while (true) {
    const {
      payload: {
        event: {
          returnValues: { recipient, amount },
        },
      },
    } = yield take(requestChan)
    const address = yield select(selectors.getUserAddress)

    if (!isSameEthAddress(recipient, address)) {
      continue
    }

    yield put({
      type: tbtcV2Migration.TBTCV2_TOKEN_MINTED,
      payload: {
        amount,
      },
    })
    // TODO: Display `MigrationCompletedModal`.
  }
}

export function* subscribeToTBTCV2UnmintedEvent() {
  const requestChan = yield actionChannel(
    tbtcV2Migration.TBTCV2_TOKEN_UNMINTED_EVENT_EMITTED
  )

  while (true) {
    const {
      payload: {
        event: {
          returnValues: { recipient, amount, fee },
        },
      },
    } = yield take(requestChan)
    const address = yield select(selectors.getUserAddress)

    if (!isSameEthAddress(recipient, address)) {
      continue
    }

    yield put({
      type: tbtcV2Migration.TBTCV2_TOKEN_UNMINTED,
      payload: {
        amount,
        fee,
      },
    })

    // TODO: Display `MigrationCompletedModal`.
  }
}

function* mint(action) {
  const { amount } = action.payload
  const vendingMachineAddress = Keep.tBTCV2Migration.vendingMachine.address

  yield call(sendTransaction, {
    payload: {
      contract: Keep.tBTCV2Migration.tbtcV1.instance,
      methodName: "approveAndCall",
      args: [vendingMachineAddress, amount, []],
    },
  })
}

function* mintWorker(action) {
  yield call(submitButtonHelper, mint, action)
}

export function* watchMintTBTCV2() {
  yield takeEvery(tbtcV2Migration.TBTCV2_MINT, mintWorker)
}

function* unmint(action) {
  const { amount } = action.payload
  const address = yield select(selectors.getUserAddress)
  const { unmintFee } = yield select(selectors.getTBTCV2Migration)

  const vendingMachineAddress = Keep.tBTCV2Migration.vendingMachine.address
  const unmintFeeFor = yield call(
    Keep.tBTCV2Migration.unmintFeeFor,
    amount,
    unmintFee
  )
  const amountToApprove = add(amount, unmintFeeFor)

  yield* approveAndTransferToken(
    address,
    vendingMachineAddress,
    amountToApprove,
    Keep.tBTCV2Migration.tbtcV2,
    Keep.tBTCV2Migration.vendingMachine,
    "unmint",
    [amount]
  )
}

function* unmintWorker(action) {
  yield call(submitButtonHelper, unmint, action)
}

export function* watchUnmintTBTCV2() {
  yield takeEvery(tbtcV2Migration.TBTCV2_UNMINT, unmintWorker)
}
