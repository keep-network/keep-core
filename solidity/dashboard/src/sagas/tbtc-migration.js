import { put, call, take, actionChannel, select } from "redux-saga/effects"
import { takeOnlyOnce } from "./effects"
import selectors from "./selectors"
import { Keep } from "../contracts"
import { tbtcV2Migration } from "../actions"
import { identifyTaskByAddress, logErrorAndThrow } from "./utils"
import { isSameEthAddress } from "../utils/general.utils"

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

export function* watchFetchTvl() {
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
  }
}
