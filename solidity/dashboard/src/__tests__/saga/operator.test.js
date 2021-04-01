import { expectSaga } from "redux-saga-test-plan"
import { throwError } from "redux-saga-test-plan/providers"
import { call } from "redux-saga/effects"
import {
  watchFetchOperatorDelegationRequest,
  watchFetchOperatorSlashedTokensRequest,
} from "../../sagas/operartor"
import { operatorService } from "../../services/token-staking.service"
import {
  FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
  FETCH_OPERATOR_SLASHED_TOKENS_FAILURE,
  FETCH_OPERATOR_SLASHED_TOKENS_START,
  FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS,
  FETCH_OPERATOR_DELEGATIONS_START,
  FETCH_OPERATOR_DELEGATIONS_SUCCESS,
  FETCH_OPERATOR_DELEGATIONS_RERQUEST,
  FETCH_OPERATOR_DELEGATIONS_FAILURE,
} from "../../actions"
import operatorReducer from "../../reducers/operator"
import { slashedTokensService } from "../../services/slashed-tokens.service"
import { ZERO_ADDRESS } from "../../utils/ethereum.utils"

describe("Test operator saga", () => {
  const mockedAddress = "0x0"
  const mockedDelegationData = {
    stakedBalance: "100",
    ownerAddress: "0x1",
    beneficiaryAddress: "0x2",
    authorizerAddress: "0x3",
  }

  const mockedSlshedTokensData = [
    { groupIndex: 1, type: "SLASHED", amount: "100" },
    { groupIndex: 2, type: "SEIZED", amount: "300" },
  ]

  const initialState = {
    stakedBalance: "0",
    ownerAddress: ZERO_ADDRESS,
    beneficiaryAddress: ZERO_ADDRESS,
    authorizerAddress: ZERO_ADDRESS,
    isFetching: false,
    error: null,

    areSlashedTokensFetching: false,
    slashedTokens: [],
    slashedTokensError: null,
  }

  describe("Delegation data fetching test", () => {
    it("should fetch operator's delegation", () => {
      return expectSaga(watchFetchOperatorDelegationRequest)
        .withReducer(operatorReducer)
        .provide([
          [
            call(operatorService.fetchDelegatedTokensData, mockedAddress),
            mockedDelegationData,
          ],
        ])
        .dispatch({
          type: FETCH_OPERATOR_DELEGATIONS_RERQUEST,
          payload: { address: mockedAddress },
        })
        .put({ type: FETCH_OPERATOR_DELEGATIONS_START })
        .put({
          type: FETCH_OPERATOR_DELEGATIONS_SUCCESS,
          payload: mockedDelegationData,
        })
        .hasFinalState({
          ...initialState,
          ...mockedDelegationData,
        })
        .run()
    })

    it("should log error if a service has failed", () => {
      const error = new Error("error")
      return expectSaga(watchFetchOperatorDelegationRequest)
        .withReducer(operatorReducer)
        .provide([
          [
            call(operatorService.fetchDelegatedTokensData, mockedAddress),
            throwError(error),
          ],
        ])
        .dispatch({
          type: FETCH_OPERATOR_DELEGATIONS_RERQUEST,
          payload: { address: mockedAddress },
        })
        .put({ type: FETCH_OPERATOR_DELEGATIONS_START })
        .put({
          type: FETCH_OPERATOR_DELEGATIONS_FAILURE,
          payload: { error: error.message },
        })
        .hasFinalState({
          ...initialState,
          error: error.message,
        })
        .run()
    })
  })

  describe("Slashed tokens data fetching test", () => {
    it("should fetch slashed tokens data correctly", () => {
      return expectSaga(watchFetchOperatorSlashedTokensRequest)
        .withReducer(operatorReducer)
        .provide([
          [
            call(slashedTokensService.fetchSlashedTokens, mockedAddress),
            mockedSlshedTokensData,
          ],
        ])
        .dispatch({
          type: FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
          payload: { address: mockedAddress },
        })
        .put({ type: FETCH_OPERATOR_SLASHED_TOKENS_START })
        .put({
          type: FETCH_OPERATOR_SLASHED_TOKENS_SUCCESS,
          payload: mockedSlshedTokensData,
        })
        .hasFinalState({
          ...initialState,
          slashedTokens: mockedSlshedTokensData,
        })
        .run()
    })

    it("should log error if a slashed service has failed", () => {
      const error = new Error("Slashed tokens error")
      return expectSaga(watchFetchOperatorSlashedTokensRequest)
        .withReducer(operatorReducer)
        .provide([
          [
            call(slashedTokensService.fetchSlashedTokens, mockedAddress),
            throwError(error),
          ],
        ])
        .dispatch({
          type: FETCH_OPERATOR_SLASHED_TOKENS_RERQUEST,
          payload: { address: mockedAddress },
        })
        .put({ type: FETCH_OPERATOR_SLASHED_TOKENS_START })
        .put({
          type: FETCH_OPERATOR_SLASHED_TOKENS_FAILURE,
          payload: { error: error.message },
        })
        .hasFinalState({
          ...initialState,
          slashedTokensError: error.message,
        })
        .run()
    })
  })
})
