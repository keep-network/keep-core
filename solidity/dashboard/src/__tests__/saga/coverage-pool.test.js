import { expectSaga } from "redux-saga-test-plan"
import { throwError } from "redux-saga-test-plan/providers"
import { call } from "redux-saga/effects"
import BigNumber from "bignumber.js"
import {
  watchFetchTvl,
  watchFetchAPY,
  watchFetchCovPoolData,
  subscribeToAuctionCreatedEvent,
  subscribeToAuctionClosedEvent,
  subscribeToAssetPoolDepositedEvent,
  subscribeToWithdrawalInitiatedEvent,
  subscribeToWithdrawalCompletedEvent,
} from "../../sagas/coverage-pool"

import coveragePoolReducer, {
  coveragePoolInitialData,
} from "../../reducers/coverage-pool"
import { KEEP, Token } from "../../utils/token.utils"
import {
  fetchTvlRequest,
  fetchTvlStart,
  fetchTvlSuccess,
  COVERAGE_POOL_FETCH_TVL_ERROR,
  fetchAPYRequest,
  fetchAPYStart,
  fetchAPYSuccess,
  COVERAGE_POOL_FETCH_APY_ERROR,
  fetchCovPoolDataRequest,
  fetchCovPoolDataStart,
  fetchCovPoolDataSuccess,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR,
  covTokenUpdated,
  RISK_MANAGER_AUCTION_CREATED_EVENT_EMITTED,
  RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED,
  COVERAGE_POOL_ASSET_POOL_DEPOSITED_EVENT_EMITTED,
  COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
  COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
} from "../../actions/coverage-pool"
import { Keep } from "../../contracts"
import * as matchers from "redux-saga-test-plan/matchers"
import { add, sub } from "../../utils/arithmetics.utils"
import { select } from "redux-saga-test-plan/matchers"
import selectors from "../../sagas/selectors"
import { OPEN_MODAL } from "../../actions/modal"
import { modalComponentType } from "../../components/Modal"

// TODO: Mock globally
// Mock TrezorConnector due to `This version of trezor-connect is not suitable
// to work without browser. Use trezor-connect@extended package instead` error.
jest.mock("../../connectors/trezor", () => ({
  ...jest.requireActual("../../components/Modal"),
  TrezorConnector: Object,
}))

describe("Coverage pool saga test", () => {
  describe("Fetch tvl watcher", () => {
    const tvl = KEEP.fromTokenUnit(100000)
    const keepInUSD = new BigNumber(0.5)
    const totalAllocatedRewards = KEEP.fromTokenUnit(200000)
    const tvlInUSD = keepInUSD.multipliedBy(KEEP.toTokenUnit(tvl)).toFormat(2)
    const totalCoverageClaimed = KEEP.fromTokenUnit(50000)

    it("should fetch tvl data correctly", () => {
      return expectSaga(watchFetchTvl)
        .withReducer(coveragePoolReducer)
        .provide([
          [call(Keep.coveragePoolV1.totalValueLocked), tvl],
          [call(Keep.exchangeService.getKeepTokenPriceInUSD), keepInUSD],
          [
            call(Keep.coveragePoolV1.totalAllocatedRewards),
            totalAllocatedRewards,
          ],
          [
            call(Keep.coveragePoolV1.totalCoverageClaimed),
            totalCoverageClaimed,
          ],
        ])
        .dispatch(fetchTvlRequest())
        .put(fetchTvlStart())
        .put(
          fetchTvlSuccess({
            tvl,
            tvlInUSD,
            totalAllocatedRewards,
            totalCoverageClaimed,
          })
        )
        .hasFinalState({
          ...coveragePoolInitialData,
          totalValueLocked: tvl,
          totalValueLockedInUSD: tvlInUSD,
          totalCoverageClaimed: totalCoverageClaimed,
          isTotalValueLockedFetching: false,
          tvlError: null,
          totalAllocatedRewards,
        })
        .run()
    })

    it("should log error if an any Keep lib function has failed", () => {
      const mockedError = new Error("Fake error")
      return expectSaga(watchFetchTvl)
        .withReducer(coveragePoolReducer)
        .provide([
          [call(Keep.coveragePoolV1.totalValueLocked), throwError(mockedError)],
          [call(Keep.exchangeService.getKeepTokenPriceInUSD), keepInUSD],
          [
            call(Keep.coveragePoolV1.totalAllocatedRewards),
            totalAllocatedRewards,
          ],
        ])
        .dispatch(fetchTvlRequest())
        .put(fetchTvlStart())
        .put({
          type: COVERAGE_POOL_FETCH_TVL_ERROR,
          payload: { error: mockedError.message },
        })
        .hasFinalState({
          ...coveragePoolInitialData,
          tvlError: mockedError.message,
        })
        .run()
    })
  })

  describe("Fetch apy watcher", () => {
    const apy = 0.25

    it("should fetch apy data correctly", () => {
      return expectSaga(watchFetchAPY)
        .withReducer(coveragePoolReducer)
        .provide([[call(Keep.coveragePoolV1.apy), apy]])
        .dispatch(fetchAPYRequest())
        .put(fetchAPYStart())
        .put(fetchAPYSuccess(apy))
        .hasFinalState({
          ...coveragePoolInitialData,
          apy,
        })
        .run()
    })

    it("should log error if function has failed", () => {
      const mockedError = new Error("Fake error")
      return expectSaga(watchFetchAPY)
        .withReducer(coveragePoolReducer)
        .provide([[call(Keep.coveragePoolV1.apy), throwError(mockedError)]])
        .dispatch(fetchAPYRequest())
        .put(fetchAPYStart())
        .put({
          type: COVERAGE_POOL_FETCH_APY_ERROR,
          payload: { error: mockedError.message },
        })
        .hasFinalState({
          ...coveragePoolInitialData,
          apyError: mockedError.message,
        })
        .run()
    })
  })

  describe("Fetch cov pool data watcher", () => {
    const balanceOf = Token.fromTokenUnit("100").toString()
    const totalSupply = Token.fromTokenUnit("1000").toString()
    const shareOfPool = 0.5
    const estimatedKeepBalance = Token.fromTokenUnit("50").toString()
    const estimatedRewards = Token.fromTokenUnit("10").toString()
    const address = "0x086813525A7dC7dafFf015Cdf03896Fd276eab60"
    const withdrawalDelays = {
      withdrawalDelay: 1,
      withdrawalTimeout: 2,
    }
    const pendingWithdrawal = 0
    const withdrawalInitiatedTimestamp = 0
    const hasOpenAuctions = true

    it("should fetch apy data correctly", () => {
      return expectSaga(watchFetchCovPoolData)
        .withReducer(coveragePoolReducer)
        .provide([
          [call(Keep.coveragePoolV1.covBalanceOf, address), balanceOf],
          [call(Keep.coveragePoolV1.covTotalSupply), totalSupply],
          [call(Keep.coveragePoolV1.withdrawalDelays), withdrawalDelays],
          [
            call(Keep.coveragePoolV1.pendingWithdrawal, address),
            pendingWithdrawal,
          ],
          [
            call(Keep.coveragePoolV1.withdrawalInitiatedTimestamp, address),
            withdrawalInitiatedTimestamp,
          ],
          [
            call(Keep.coveragePoolV1.shareOfPool, totalSupply, balanceOf),
            shareOfPool,
          ],
          [
            call(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance,
              shareOfPool
            ),
            estimatedKeepBalance,
          ],
          [
            call(Keep.coveragePoolV1.estimatedRewards, address, shareOfPool),
            estimatedRewards,
          ],
          [
            call(Keep.coveragePoolV1.hasRiskManagerOpenAuctions),
            hasOpenAuctions,
          ],
        ])
        .dispatch(fetchCovPoolDataRequest(address))
        .put(fetchCovPoolDataStart())
        .put(
          fetchCovPoolDataSuccess({
            shareOfPool,
            covBalance: balanceOf,
            covTokensAvailableToWithdraw: balanceOf,
            covTotalSupply: totalSupply,
            estimatedRewards,
            estimatedKeepBalance,
            withdrawalDelay: withdrawalDelays.withdrawalDelay,
            withdrawalTimeout: withdrawalDelays.withdrawalTimeout,
            pendingWithdrawal,
            withdrawalInitiatedTimestamp,
            hasRiskManagerOpenAuctions: hasOpenAuctions,
          })
        )
        .hasFinalState({
          ...coveragePoolInitialData,
          shareOfPool,
          covBalance: balanceOf,
          covTokensAvailableToWithdraw: balanceOf,
          covTotalSupply: totalSupply,
          estimatedRewards,
          estimatedKeepBalance,
          withdrawalDelay: withdrawalDelays.withdrawalDelay,
          withdrawalTimeout: withdrawalDelays.withdrawalTimeout,
          pendingWithdrawal,
          withdrawalInitiatedTimestamp,
          hasRiskManagerOpenAuctions: hasOpenAuctions,
        })
        .run()
    })

    it("should log error if function has failed", () => {
      const mockedError = new Error("Fake error")

      return expectSaga(watchFetchCovPoolData)
        .withReducer(coveragePoolReducer)
        .provide([
          [
            call(Keep.coveragePoolV1.covBalanceOf, address),
            throwError(mockedError),
          ],
        ])
        .dispatch(fetchCovPoolDataRequest(address))
        .put(fetchCovPoolDataStart())
        .put({
          type: COVERAGE_POOL_FETCH_COV_POOL_DATA_ERROR,
          payload: { error: mockedError.message },
        })
        .hasFinalState({
          ...coveragePoolInitialData,
          error: mockedError.message,
        })
        .run()
    })
  })

  describe("Subscribe to Asset Pool's Deposit event", () => {
    const address = "0x086813525A7dC7dafFf015Cdf03896Fd276eab60"
    const initialCovTotalSupply = Token.fromTokenUnit(100).toString()
    const initialCovBalance = Token.fromTokenUnit(30).toString()
    const initialCovTokensAvailableToWithdraw =
      Token.fromTokenUnit(30).toString()

    const depositEventData = {
      underwriter: address,
      covAmount: KEEP.fromTokenUnit("300").toString(),
      amount: KEEP.fromTokenUnit("350").toString(),
    }

    const initialState = {
      ...coveragePoolInitialData,
      covTotalSupply: initialCovTotalSupply,
      covBalance: initialCovBalance,
      covTokensAvailableToWithdraw: initialCovTokensAvailableToWithdraw,
    }
    const updatedShareOfPool = 0.8
    const estimatedKeepBalance = KEEP.fromTokenUnit(350).toString()
    const estimatedRewards = KEEP.fromTokenUnit(35).toString()
    const updatedTvl = KEEP.fromTokenUnit(10000).toString()
    const keepInUSD = new BigNumber(0.25)
    const updatedAPY = 0.5
    const tvlInUSD = new BigNumber(keepInUSD)
      .multipliedBy(KEEP.toTokenUnit(updatedTvl))
      .toFormat(2)

    it("should update data correctly if the `Deposit` event has been emitted by current logged user", () => {
      const mockedEvent = {
        returnValues: { ...depositEventData },
      }

      const updatedCovBalance = add(
        initialCovBalance,
        depositEventData.covAmount
      ).toString()
      const updatedCovTotalSupply = add(
        initialCovTotalSupply,
        depositEventData.covAmount
      ).toString()
      const updatedCovTokensAvailableToWithdraw = add(
        initialCovTokensAvailableToWithdraw,
        depositEventData.covAmount
      ).toString()

      return expectSaga(subscribeToAssetPoolDepositedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [
            matchers.call.fn(Keep.coveragePoolV1.shareOfPool),
            updatedShareOfPool,
          ],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance
            ),
            estimatedKeepBalance,
          ],
          [
            matchers.call.fn(Keep.coveragePoolV1.estimatedRewards),
            estimatedRewards,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.totalValueLocked), updatedTvl],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.exchangeService.getKeepTokenPriceInUSD
            ),
            keepInUSD,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.apy), updatedAPY],
        ])
        .dispatch({
          type: COVERAGE_POOL_ASSET_POOL_DEPOSITED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put(
          covTokenUpdated({
            covBalance: updatedCovBalance,
            covTokensAvailableToWithdraw: updatedCovTokensAvailableToWithdraw,
            covTotalSupply: updatedCovTotalSupply,
            shareOfPool: updatedShareOfPool,
            estimatedKeepBalance,
            estimatedRewards,
            totalValueLocked: updatedTvl,
            totalValueLockedInUSD: tvlInUSD,
            apy: updatedAPY,
          })
        )
        .hasFinalState({
          ...initialState,
          covBalance: updatedCovBalance,
          covTokensAvailableToWithdraw: updatedCovTokensAvailableToWithdraw,
          covTotalSupply: updatedCovTotalSupply,
          shareOfPool: updatedShareOfPool,
          estimatedKeepBalance,
          estimatedRewards,
          totalValueLocked: updatedTvl,
          totalValueLockedInUSD: tvlInUSD,
          apy: updatedAPY,
        })
        .run()
    })

    it("should update data correctly for current user if the `Deposit` event has been emitted by another user", () => {
      const mockedEvent = {
        returnValues: {
          ...depositEventData,
          underwriter: "0x065993c332b02ab8674Ac033CaCDBccBe7bc9047",
        },
      }

      const updatedCovTotalSupply = add(
        initialCovTotalSupply,
        depositEventData.covAmount
      ).toString()

      return expectSaga(subscribeToAssetPoolDepositedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [
            matchers.call.fn(Keep.coveragePoolV1.shareOfPool),
            updatedShareOfPool,
          ],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance
            ),
            estimatedKeepBalance,
          ],
          [
            matchers.call.fn(Keep.coveragePoolV1.estimatedRewards),
            estimatedRewards,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.totalValueLocked), updatedTvl],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.exchangeService.getKeepTokenPriceInUSD
            ),
            keepInUSD,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.apy), updatedAPY],
        ])
        .dispatch({
          type: COVERAGE_POOL_ASSET_POOL_DEPOSITED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put(
          covTokenUpdated({
            covBalance: initialCovBalance,
            covTokensAvailableToWithdraw: initialCovTokensAvailableToWithdraw,
            covTotalSupply: updatedCovTotalSupply,
            shareOfPool: updatedShareOfPool,
            estimatedKeepBalance,
            estimatedRewards,
            totalValueLocked: updatedTvl,
            totalValueLockedInUSD: tvlInUSD,
            apy: updatedAPY,
          })
        )
        .hasFinalState({
          ...initialState,
          covBalance: initialCovBalance,
          covTokensAvailableToWithdraw: initialCovTokensAvailableToWithdraw,
          covTotalSupply: updatedCovTotalSupply,
          shareOfPool: updatedShareOfPool,
          estimatedKeepBalance,
          estimatedRewards,
          totalValueLocked: updatedTvl,
          totalValueLockedInUSD: tvlInUSD,
          apy: updatedAPY,
        })
        .run()
    })

    it("should show correct success modal window, when depositing KEEPs in Coverage Pool", () => {
      const transactionHash = "0x123453525A7454et543tr5Cdf03896Fd276avg45"
      const mockedEvent = {
        returnValues: { ...depositEventData },
        transactionHash,
      }

      const expectedModalType =
        modalComponentType.COV_POOLS.KEEP_DEPOSITED_SUCCESS
      const expectedTitle = "Deposit"

      return expectSaga(subscribeToAssetPoolDepositedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [
            matchers.call.fn(Keep.coveragePoolV1.shareOfPool),
            updatedShareOfPool,
          ],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance
            ),
            estimatedKeepBalance,
          ],
          [
            matchers.call.fn(Keep.coveragePoolV1.estimatedRewards),
            estimatedRewards,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.totalValueLocked), updatedTvl],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.exchangeService.getKeepTokenPriceInUSD
            ),
            keepInUSD,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.apy), updatedAPY],
        ])
        .dispatch({
          type: COVERAGE_POOL_ASSET_POOL_DEPOSITED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put.like({
          action: {
            type: OPEN_MODAL,
            payload: {
              modalComponentType: expectedModalType,
              modalProps: {
                title: expectedTitle,
              },
              componentProps: {
                transactionFinished: true,
                transactionHash: mockedEvent.transactionHash,
                amount: mockedEvent.returnValues.amount,
              },
            },
          },
        })
        .run()
    })
  })

  describe("Subscribe to WithdrawalInitiated event", () => {
    const address = "0x086813525A7dC7dafFf015Cdf03896Fd276eab60"
    const transactionHash = "0x123453525A7454et543tr5Cdf03896Fd276avg45"
    const initialCovTotalSupply = Token.fromTokenUnit(100).toString()
    const initialCovBalance = Token.fromTokenUnit(30).toString()
    const initialCovTokensAvailableToWithdraw =
      Token.fromTokenUnit(30).toString()

    const initialState = {
      ...coveragePoolInitialData,
      covTotalSupply: initialCovTotalSupply,
      covBalance: initialCovBalance,
      covTokensAvailableToWithdraw: initialCovTokensAvailableToWithdraw,
    }

    it("should update data correctly if the `WithdrawalInitiated` event has been emitted by current logged user", () => {
      const underwriterAddress = address
      const withdrawalInitiatedEventData = {
        underwriter: underwriterAddress,
        covAmount: KEEP.fromTokenUnit("20").toString(),
        timestamp: 1631857610,
      }
      const mockedEvent = {
        returnValues: withdrawalInitiatedEventData,
      }

      const mockedModalData = {
        componentProps: {
          pendingWithdrawalBalance: KEEP.fromTokenUnit("200").toString(),
          amount: KEEP.fromTokenUnit("20").toString(),
        },
      }

      const updatedCovTokensAvailableToWithdraw = sub(
        initialCovTokensAvailableToWithdraw,
        withdrawalInitiatedEventData.covAmount
      ).toString()

      return expectSaga(subscribeToWithdrawalInitiatedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put(
          covTokenUpdated({
            pendingWithdrawal: withdrawalInitiatedEventData.covAmount,
            withdrawalInitiatedTimestamp:
              withdrawalInitiatedEventData.timestamp,
            covTokensAvailableToWithdraw: updatedCovTokensAvailableToWithdraw,
          })
        )
        .hasFinalState({
          ...initialState,
          pendingWithdrawal: withdrawalInitiatedEventData.covAmount,
          withdrawalInitiatedTimestamp: withdrawalInitiatedEventData.timestamp,
          covTokensAvailableToWithdraw: updatedCovTokensAvailableToWithdraw,
        })
        .run()
    })

    it("should update data correctly for current user if the `WithdrawalInitiated` event has been emitted by another user", () => {
      const underwriterAddress = "0x065993c332b02ab8674Ac033CaCDBccBe7bc9047"
      const withdrawalInitiatedEventData = {
        underwriter: underwriterAddress,
        covAmount: KEEP.fromTokenUnit("20").toString(),
        timestamp: 1631857610,
      }
      const mockedEvent = {
        returnValues: withdrawalInitiatedEventData,
      }

      const mockedModalData = {
        componentProps: {
          pendingWithdrawalBalance: KEEP.fromTokenUnit("200").toString(),
          amount: KEEP.fromTokenUnit("20").toString(),
        },
      }

      return expectSaga(subscribeToWithdrawalInitiatedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .hasFinalState({
          ...initialState,
        })
        .run()
    })

    it("should show correct success modal window, when initiating a withdrawal", () => {
      const underwriterAddress = address
      const transactionHash = "0x123453525A7454et543tr5Cdf03896Fd276avg45"
      const withdrawalInitiatedEventData = {
        underwriter: underwriterAddress,
        covAmount: KEEP.fromTokenUnit("20").toString(),
        timestamp: 1631857610,
      }
      const mockedEvent = {
        returnValues: withdrawalInitiatedEventData,
        transactionHash,
      }

      const mockedModalData = {
        componentProps: {
          amount: KEEP.fromTokenUnit("20").toString(),
        },
      }

      const expectedModalType = modalComponentType.COV_POOLS.INITIATE_WITHDRAWAL
      const expectedTitle = "Withdraw"

      return expectSaga(subscribeToWithdrawalInitiatedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put.like({
          action: {
            type: OPEN_MODAL,
            payload: {
              modalComponentType: expectedModalType,
              modalProps: {
                title: expectedTitle,
              },
            },
          },
        })
        .run()
    })

    it("should show correct success modal window, when re-initiating a withdrawal", () => {
      const underwriterAddress = address
      const withdrawalInitiatedEventData = {
        underwriter: underwriterAddress,
        covAmount: KEEP.fromTokenUnit("20").toString(),
        timestamp: 1631857610,
      }
      const mockedEvent = {
        returnValues: withdrawalInitiatedEventData,
        transactionHash,
      }

      const mockedModalData = {
        componentProps: {
          pendingWithdrawalBalance: Token.fromTokenUnit(20).toString(),
          amount: KEEP.fromTokenUnit("0").toString(),
        },
      }

      const expectedModalType =
        modalComponentType.COV_POOLS.RE_INITIATE_WITHDRAWAL
      const expectedTitle = "Re-initiate withdrawal"

      return expectSaga(subscribeToWithdrawalInitiatedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put.like({
          action: {
            type: OPEN_MODAL,
            payload: {
              modalComponentType: expectedModalType,
              modalProps: {
                title: expectedTitle,
              },
            },
          },
        })
        .run()
    })

    it("should show correct success modal window, when increasing value of a pending withdrawal", () => {
      const underwriterAddress = address
      const transactionHash = "0x123453525A7454et543tr5Cdf03896Fd276avg45"
      const withdrawalInitiatedEventData = {
        underwriter: underwriterAddress,
        covAmount: KEEP.fromTokenUnit("20").toString(),
        timestamp: 1631857610,
      }
      const mockedEvent = {
        returnValues: withdrawalInitiatedEventData,
        transactionHash,
      }

      const mockedModalData = {
        componentProps: {
          pendingWithdrawalBalance: Token.fromTokenUnit(20).toString(),
          amount: KEEP.fromTokenUnit("10").toString(),
        },
      }

      const expectedModalType = modalComponentType.COV_POOLS.INCREASE_WITHDRAWAL
      const expectedTitle = "Re-initiate withdrawal"

      return expectSaga(subscribeToWithdrawalInitiatedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_INITIATED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put.like({
          action: {
            type: OPEN_MODAL,
            payload: {
              modalComponentType: expectedModalType,
              modalProps: {
                title: expectedTitle,
              },
            },
          },
        })
        .run()
    })
  })

  describe("Subscribe to WithdrawalCompleted event", () => {
    const address = "0x086813525A7dC7dafFf015Cdf03896Fd276eab60"
    const transactionHash = "0x123453525A7454et543tr5Cdf03896Fd276avg45"
    const initialCovTotalSupply = Token.fromTokenUnit("1000").toString()
    const initialCovBalance = Token.fromTokenUnit("400").toString()

    const withdrawalCompletedEventData = {
      underwriter: address,
      amount: KEEP.fromTokenUnit("320").toString(),
      covAmount: KEEP.fromTokenUnit("300").toString(),
    }
    const initialCovTokensAvailableToWithdraw = sub(
      initialCovBalance,
      withdrawalCompletedEventData.covAmount
    ).toString()

    const mockedEvent = {
      returnValues: withdrawalCompletedEventData,
      transactionHash,
    }

    const mockedModalData = {
      componentProps: {
        pendingWithdrawalBalance: KEEP.fromTokenUnit("0").toString(),
        amount: KEEP.fromTokenUnit("0").toString(),
      },
    }

    const initialState = {
      ...coveragePoolInitialData,
      covTotalSupply: initialCovTotalSupply,
      covBalance: initialCovBalance,
      covTokensAvailableToWithdraw: initialCovTokensAvailableToWithdraw,
    }

    const updatedShareOfPool = 0.8
    const estimatedKeepBalance = KEEP.fromTokenUnit(350).toString()
    const estimatedRewards = KEEP.fromTokenUnit(35).toString()
    const updatedTvl = KEEP.fromTokenUnit(10000).toString()
    const keepInUSD = new BigNumber(0.25)
    const updatedAPY = 0.5
    const tvlInUSD = new BigNumber(keepInUSD)
      .multipliedBy(KEEP.toTokenUnit(updatedTvl))
      .toFormat(2)

    it("should update data correctly if the `WithdrawalCompleted` event has been emitted by current logged user", () => {
      const updatedCovBalance = sub(
        initialCovBalance,
        withdrawalCompletedEventData.covAmount
      ).toString()
      const updatedCovTotalSupply = sub(
        initialCovTotalSupply,
        withdrawalCompletedEventData.covAmount
      ).toString()

      return expectSaga(subscribeToWithdrawalCompletedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
          [
            matchers.call.fn(Keep.coveragePoolV1.shareOfPool),
            updatedShareOfPool,
          ],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance
            ),
            estimatedKeepBalance,
          ],
          [
            matchers.call.fn(Keep.coveragePoolV1.estimatedRewards),
            estimatedRewards,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.totalValueLocked), updatedTvl],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.exchangeService.getKeepTokenPriceInUSD
            ),
            keepInUSD,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.apy), updatedAPY],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put(
          covTokenUpdated({
            shareOfPool: updatedShareOfPool,
            covBalance: updatedCovBalance,
            covTotalSupply: updatedCovTotalSupply,
            estimatedRewards,
            estimatedKeepBalance,
            totalValueLockedInUSD: tvlInUSD,
            totalValueLocked: updatedTvl,
            apy: updatedAPY,
            pendingWithdrawal: "0",
            withdrawalInitiatedTimestamp: "0",
          })
        )
        .hasFinalState({
          ...initialState,
          shareOfPool: updatedShareOfPool,
          covBalance: updatedCovBalance,
          covTotalSupply: updatedCovTotalSupply,
          estimatedRewards,
          estimatedKeepBalance,
          totalValueLockedInUSD: tvlInUSD,
          totalValueLocked: updatedTvl,
          apy: updatedAPY,
          pendingWithdrawal: "0",
          withdrawalInitiatedTimestamp: "0",
        })
        .run()
    })

    it("should update data correctly for current user if the `WithdrawalCompleted` event has been emitted by another user", () => {
      const underwriterAddress = "0x065993c332b02ab8674Ac033CaCDBccBe7bc9047"

      const withdrawalCompletedEventData = {
        underwriter: underwriterAddress,
        amount: KEEP.fromTokenUnit("320").toString(),
        covAmount: KEEP.fromTokenUnit("300").toString(),
      }
      const initialCovTokensAvailableToWithdraw = sub(
        initialCovBalance,
        withdrawalCompletedEventData.covAmount
      ).toString()

      const mockedEvent = {
        returnValues: withdrawalCompletedEventData,
      }

      const initialState = {
        ...coveragePoolInitialData,
        covTotalSupply: initialCovTotalSupply,
        covBalance: initialCovBalance,
        covTokensAvailableToWithdraw: initialCovTokensAvailableToWithdraw,
        pendingWithdrawal: "300",
        withdrawalInitiatedTimestamp: "1631857111",
      }

      const updatedCovTotalSupply = sub(
        initialCovTotalSupply,
        withdrawalCompletedEventData.covAmount
      ).toString()

      return expectSaga(subscribeToWithdrawalCompletedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [select(selectors.getModalData), mockedModalData],
          [
            matchers.call.fn(Keep.coveragePoolV1.shareOfPool),
            updatedShareOfPool,
          ],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance
            ),
            estimatedKeepBalance,
          ],
          [
            matchers.call.fn(Keep.coveragePoolV1.estimatedRewards),
            estimatedRewards,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.totalValueLocked), updatedTvl],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.exchangeService.getKeepTokenPriceInUSD
            ),
            keepInUSD,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.apy), updatedAPY],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put(
          covTokenUpdated({
            shareOfPool: updatedShareOfPool,
            covBalance: initialCovBalance,
            covTotalSupply: updatedCovTotalSupply,
            estimatedRewards,
            estimatedKeepBalance,
            totalValueLockedInUSD: tvlInUSD,
            totalValueLocked: updatedTvl,
            apy: updatedAPY,
          })
        )
        .hasFinalState({
          ...initialState,
          shareOfPool: updatedShareOfPool,
          covBalance: initialCovBalance,
          covTotalSupply: updatedCovTotalSupply,
          estimatedRewards,
          estimatedKeepBalance,
          totalValueLockedInUSD: tvlInUSD,
          totalValueLocked: updatedTvl,
          apy: updatedAPY,
        })
        .run()
    })

    it("should show correct success modal window, when claiming tokens from a pending withdrawal", () => {
      const expectedModalType =
        modalComponentType.COV_POOLS.WITHDRAWAL_COMPLETED
      const expectedTitle = "Claim tokens"

      return expectSaga(subscribeToWithdrawalCompletedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [select(selectors.getCoveragePool), initialState],
          [select(selectors.getUserAddress), address],
          [
            matchers.call.fn(Keep.coveragePoolV1.shareOfPool),
            updatedShareOfPool,
          ],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.estimatedCollateralTokenBalance
            ),
            estimatedKeepBalance,
          ],
          [
            matchers.call.fn(Keep.coveragePoolV1.estimatedRewards),
            estimatedRewards,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.totalValueLocked), updatedTvl],
          [
            matchers.call.fn(
              Keep.coveragePoolV1.exchangeService.getKeepTokenPriceInUSD
            ),
            keepInUSD,
          ],
          [matchers.call.fn(Keep.coveragePoolV1.apy), updatedAPY],
        ])
        .dispatch({
          type: COVERAGE_POOL_WITHDRAWAL_COMPLETED_EVENT_EMITTED,
          payload: { event: mockedEvent },
        })
        .put.like({
          action: {
            type: OPEN_MODAL,
            payload: {
              modalComponentType: expectedModalType,
              modalProps: {
                title: expectedTitle,
              },
            },
          },
        })
        .run()
    })
  })

  describe("Subscribe to AuctionCreated event", () => {
    it("should update data correctly if `AuctionCreated` event has been emitted", () => {
      const initialState = {
        ...coveragePoolInitialData,
        hasRiskManagerOpenAuctions: false,
      }

      return expectSaga(subscribeToAuctionCreatedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .dispatch({
          type: RISK_MANAGER_AUCTION_CREATED_EVENT_EMITTED,
        })
        .put(
          covTokenUpdated({
            hasRiskManagerOpenAuctions: true,
          })
        )
        .hasFinalState({
          ...initialState,
          hasRiskManagerOpenAuctions: true,
        })
        .run()
    })

    it("should update data correctly if `AuctionClosed` event has been emitted and there are no more open auctions left", () => {
      const initialState = {
        ...coveragePoolInitialData,
        hasRiskManagerOpenAuctions: false,
      }

      const mockedHasOpenAuctions = false

      return expectSaga(subscribeToAuctionClosedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [
            matchers.call.fn(Keep.coveragePoolV1.hasRiskManagerOpenAuctions),
            mockedHasOpenAuctions,
          ],
        ])
        .dispatch({
          type: RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED,
        })
        .put(
          covTokenUpdated({
            hasRiskManagerOpenAuctions: mockedHasOpenAuctions,
          })
        )
        .hasFinalState({
          ...initialState,
          hasRiskManagerOpenAuctions: mockedHasOpenAuctions,
        })
        .run()
    })

    it("should update data correctly if `AuctionClosed` event has been emitted but open auctions still exist", () => {
      const initialState = {
        ...coveragePoolInitialData,
        hasRiskManagerOpenAuctions: false,
      }

      const mockedHasOpenAuctions = true

      return expectSaga(subscribeToAuctionClosedEvent)
        .withReducer(coveragePoolReducer, initialState.coveragePool)
        .withState(initialState)
        .provide([
          [
            matchers.call.fn(Keep.coveragePoolV1.hasRiskManagerOpenAuctions),
            mockedHasOpenAuctions,
          ],
        ])
        .dispatch({
          type: RISK_MANAGER_AUCTION_CLOSED_EVENT_EMITTED,
        })
        .put(
          covTokenUpdated({
            hasRiskManagerOpenAuctions: mockedHasOpenAuctions,
          })
        )
        .hasFinalState({
          ...initialState,
          hasRiskManagerOpenAuctions: mockedHasOpenAuctions,
        })
        .run()
    })
  })
})
