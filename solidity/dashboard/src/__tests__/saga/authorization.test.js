import { expectSaga } from "redux-saga-test-plan"
import { throwError } from "redux-saga-test-plan/providers"
import { call } from "redux-saga/effects"
import { watchFetchKeepRandomBeaconAuthData } from "../../sagas/authorization"
import authorizationReducer from "../../reducers/authorization"
import { beaconAuthorizationService } from "../../services/beacon-authorization.service"
import {
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS,
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE,
  KEEP_RRANDOM_BEACON_AUTHORIZED,
} from "../../actions"

describe("Authorization saga test", () => {
  const mockedAddress = "0x0"
  describe("Keep Random Beacon watchers", () => {
    const contractName = "Keep Random Beacon Operator Contract"
    const operatorAddress = 0x1
    const mockedRandomBeaconAuthData = [
      {
        operatorAddress,
        stakeAmount: 100,
        contracts: [
          {
            contractName,
            operatorContractAddress: "0x2",
            isAuthorized: false,
          },
        ],
      },
    ]

    it("should fetch Keep Random Beacon authorization data correctly", () => {
      return expectSaga(watchFetchKeepRandomBeaconAuthData)
        .withReducer(authorizationReducer)
        .provide([
          [
            call(
              beaconAuthorizationService.fetchRandomBeaconAuthorizationData,
              mockedAddress
            ),
            mockedRandomBeaconAuthData,
          ],
        ])
        .dispatch({
          type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
          payload: { address: mockedAddress },
        })
        .put({ type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START })
        .put({
          type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_SUCCESS,
          payload: mockedRandomBeaconAuthData,
        })
        .hasFinalState({
          authData: mockedRandomBeaconAuthData,
          isFetching: false,
          error: null,
        })
        .run()
    })

    it("should log error if the beacon auth service has failed", () => {
      const error = new Error("Beacon auth error")
      return expectSaga(watchFetchKeepRandomBeaconAuthData)
        .withReducer(authorizationReducer)
        .provide([
          [
            call(
              beaconAuthorizationService.fetchRandomBeaconAuthorizationData,
              mockedAddress
            ),
            throwError(error),
          ],
        ])
        .dispatch({
          type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
          payload: { address: mockedAddress },
        })
        .put({ type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_START })
        .put({
          type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_FAILURE,
          payload: { error: error.message },
        })
        .hasFinalState({
          authData: [],
          isFetching: false,
          error: error.message,
        })
        .run()
    })

    it("should update the authorized contract correctly", () => {
      const expectedAuthData = [...mockedRandomBeaconAuthData]
      expectedAuthData[0].contracts[0].isAuthorized = true
      return expectSaga(watchFetchKeepRandomBeaconAuthData)
        .withReducer(authorizationReducer)
        .provide([
          [
            call(
              beaconAuthorizationService.fetchRandomBeaconAuthorizationData,
              mockedAddress
            ),
            mockedRandomBeaconAuthData,
          ],
        ])
        .dispatch({
          type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
          payload: { address: mockedAddress },
        })
        .dispatch({
          type: KEEP_RRANDOM_BEACON_AUTHORIZED,
          payload: { contractName, operatorAddress },
        })
        .hasFinalState({
          authData: expectedAuthData,
          isFetching: false,
          error: null,
        })
        .run()
    })
  })
})
