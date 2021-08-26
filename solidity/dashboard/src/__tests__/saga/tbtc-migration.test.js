import { expectSaga } from "redux-saga-test-plan"
import { throwError } from "redux-saga-test-plan/providers"
import { call, select } from "redux-saga/effects"
import { Keep } from "../../contracts"
import tbtcV2MigrationReducer, {
  initialState,
} from "../../reducers/tbtc-migration"
import {
  subscribeToTBTCV2MintedEvent,
  subscribeToTBTCV2UnmintedEvent,
  watchFetchData,
} from "../../sagas/tbtc-migration"
import { TBTC } from "../../utils/token.utils"
import { tbtcV2Migration } from "../../actions"
import selectors from "../../sagas/selectors"
import { add, sub } from "../../utils/arithmetics.utils"

describe("TBTC migration saga test", () => {
  const mockedAddress = "0xa9d41A6a24312505866E38C217e39A447d6b66B4"
  const tbtcV1Balance = TBTC.fromTokenUnit("1").toString()
  const tbtcV2Balance = TBTC.fromTokenUnit("2").toString()
  const fee = TBTC.fromTokenUnit("0.001").toString()

  test("should fetch tbtc migration data correctly", () => {
    return expectSaga(watchFetchData)
      .withReducer(tbtcV2MigrationReducer)
      .provide([
        [
          call(Keep.tBTCV2Migration.tbtcV1BalanceOf, mockedAddress),
          tbtcV1Balance,
        ],
        [
          call(Keep.tBTCV2Migration.tbtcV2BalanceOf, mockedAddress),
          tbtcV2Balance,
        ],
        [call(Keep.tBTCV2Migration.unmintFee), fee],
      ])
      .dispatch(tbtcV2Migration.fetchDataRequest(mockedAddress))
      .put(
        tbtcV2Migration.fetchDataSuccess({
          tbtcV1Balance,
          tbtcV2Balance,
          unmintFee: fee,
        })
      )
      .hasFinalState({
        ...initialState,
        tbtcV1Balance,
        tbtcV2Balance,
        unmintFee: fee,
      })
      .silentRun()
  })

  test("should log error if an any Keep.tbtcV2Migration function has failed", () => {
    const mockedError = new Error("Fake error")

    return expectSaga(watchFetchData)
      .withReducer(tbtcV2MigrationReducer)
      .provide([
        [
          call(Keep.tBTCV2Migration.tbtcV1BalanceOf, mockedAddress),
          throwError(mockedError),
        ],
        [
          call(Keep.tBTCV2Migration.tbtcV2BalanceOf, mockedAddress),
          tbtcV2Balance,
        ],
        [call(Keep.tBTCV2Migration.unmintFee), fee],
      ])
      .dispatch(tbtcV2Migration.fetchDataRequest(mockedAddress))
      .put({
        type: tbtcV2Migration.TBTCV2_MIGRATION_FETCH_DATA_ERROR,
        payload: { error: mockedError.message },
      })
      .hasFinalState({
        ...initialState,
        error: mockedError.message,
      })
      .silentRun()
  })

  test("should update data when tbtc v2 token has been minted", () => {
    const mockedEventData = {
      recipient: mockedAddress,
      amount: TBTC.fromTokenUnit(2).toString(),
    }
    return (
      expectSaga(subscribeToTBTCV2MintedEvent)
        .withReducer(tbtcV2MigrationReducer)
        .withState({ ...initialState, tbtcV1Balance: mockedEventData.amount })
        .provide([[select(selectors.getUserAddress), mockedAddress]])
        .actionChannel(tbtcV2Migration.TBTCV2_TOKEN_MINTED_EVENT_EMITTED)
        .dispatch({
          type: tbtcV2Migration.TBTCV2_TOKEN_MINTED_EVENT_EMITTED,
          payload: { event: { returnValues: mockedEventData } },
        })
        // The second event has been emitted bu a recipient is different than
        // current logged account- we should not update data.
        .dispatch({
          type: tbtcV2Migration.TBTCV2_TOKEN_MINTED_EVENT_EMITTED,
          payload: {
            event: {
              returnValues: {
                ...mockedEventData,
                recipient: "0xFf24F5AF38bab289F13d45155e0dB89b9435163f",
              },
            },
          },
        })
        .put({
          type: tbtcV2Migration.TBTCV2_TOKEN_MINTED,
          payload: { amount: mockedEventData.amount },
        })
        .hasFinalState({
          ...initialState,
          tbtcV1Balance: "0",
          tbtcV2Balance: mockedEventData.amount,
        })
        .silentRun()
    )
  })

  test("should update data when tbtc v2 token has been unminted", () => {
    const mockedEventData = {
      recipient: mockedAddress,
      amount: TBTC.fromTokenUnit(2).toString(),
      fee: TBTC.fromTokenUnit(0.002).toString(),
    }
    const initialTBTCV2Balance = add(
      mockedEventData.amount,
      mockedEventData.fee
    ).toString()
    return (
      expectSaga(subscribeToTBTCV2UnmintedEvent)
        .withReducer(tbtcV2MigrationReducer)
        .withState({ ...initialState, tbtcV2Balance: initialTBTCV2Balance })
        .provide([[select(selectors.getUserAddress), mockedAddress]])
        .actionChannel(tbtcV2Migration.TBTCV2_TOKEN_UNMINTED_EVENT_EMITTED)
        .dispatch({
          type: tbtcV2Migration.TBTCV2_TOKEN_UNMINTED_EVENT_EMITTED,
          payload: { event: { returnValues: mockedEventData } },
        })
        // The second event has been emitted bu a recipient is different than
        // current logged account- we should not update data.
        .dispatch({
          type: tbtcV2Migration.TBTCV2_TOKEN_UNMINTED_EVENT_EMITTED,
          payload: {
            event: {
              returnValues: {
                ...mockedEventData,
                recipient: "0xFf24F5AF38bab289F13d45155e0dB89b9435163f",
              },
            },
          },
        })
        .put({
          type: tbtcV2Migration.TBTCV2_TOKEN_UNMINTED,
          payload: { amount: mockedEventData.amount, fee: mockedEventData.fee },
        })
        .hasFinalState({
          ...initialState,
          tbtcV1Balance: mockedEventData.amount,
          tbtcV2Balance: sub(
            initialTBTCV2Balance,
            add(mockedEventData.amount, mockedEventData.fee)
          ).toString(),
        })
        .silentRun()
    )
  })
})
