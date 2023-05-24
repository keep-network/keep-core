import { all, fork, take, cancel, put } from "redux-saga/effects"
import { createMockTask } from "@redux-saga/testing-utils"
import { expectSaga } from "redux-saga-test-plan"
import rootSaga, { runTasks } from "../../sagas"
import * as messagesSaga from "../../sagas/messages"
import * as delegateStakeSaga from "../../sagas/staking"
import * as tokenGrantSaga from "../../sagas/token-grant"
import {
  watchSendRawTransactionsInSequenceRequest,
  watchSendTransactionRequest,
} from "../../sagas/web3"
import * as copyStakeSaga from "../../sagas/copy-stake"
import * as subscriptions from "../../sagas/subscriptions"
import * as keepTokenBalance from "../../sagas/keep-balance"
import * as operator from "../../sagas/operartor"
import * as authorization from "../../sagas/authorization"
import * as coveragePool from "../../sagas/coverage-pool"

// TODO: Mock globally
// Mock TrezorConnector due to `This version of trezor-connect is not suitable
// to work without browser. Use trezor-connect@extended package instead` error.
jest.mock("../../connectors/trezor", () => ({
  ...jest.requireActual("../../components/Modal"),
  TrezorConnector: Object,
}))

// TODO: Mock globally
// Mock TrezorConnector due to `This version of trezor-connect is not suitable
// to work without browser. Use trezor-connect@extended package instead` error.
jest.mock("../../connectors/trezor", () => ({
  ...jest.requireActual("../../components/Modal"),
  TrezorConnector: Object,
}))

const sagas = [...Object.values(messagesSaga), ...Object.values(coveragePool)]

const loginRequiredSagas = [
  ...Object.values(delegateStakeSaga),
  watchSendTransactionRequest,
  watchSendRawTransactionsInSequenceRequest,
  ...Object.values(tokenGrantSaga),
  ...Object.values(copyStakeSaga),
  ...Object.values(subscriptions),
  ...Object.values(keepTokenBalance),
  ...Object.values(operator),
  ...Object.values(authorization),
  ...Object.values(tbtcMigration),
]

describe("Test root saga", () => {
  it("should start correctly and handle login flow", () => {
    const mockTasks = [fork(() => {}), fork(() => {})]

    return expectSaga(rootSaga)
      .provide([[all(), mockTasks]])
      .dispatch({ type: "app/login", payload: { address: "0x0" } })
      .put({ type: "app/set_account", payload: { address: "0x0" } })
      .dispatch({ type: "app/logout" })
      .put({ type: "app/reset_store" })
      .run()
  })

  it("should handle account switching flow", () => {
    return expectSaga(rootSaga)
      .dispatch({ type: "app/login", payload: { address: "0x0" } })
      .put({ type: "app/set_account", payload: { address: "0x0" } })
      .dispatch({ type: "app/account_changed", payload: { address: "0x1" } })
      .put({ type: "app/reset_store" })
      .put({ type: "app/set_account", payload: { address: "0x1" } })
      .not.put({ type: "app/logout" })
      .run()
  })
})

describe("Test root saga step by step", () => {
  let generator = null
  const mockTask = createMockTask()
  const mockAddress = "0x0"

  beforeAll(() => {
    generator = rootSaga()
  })

  it("should run sagas", () => {
    const expectedYieldAll = all(sagas.map(fork))
    expect(generator.next().value).toStrictEqual(expectedYieldAll)
  })

  it("should wait for start action", () => {
    const expectedYield = take("app/login")
    expect(generator.next().value).toStrictEqual(expectedYield)
  })

  it("should set account", () => {
    const expectedYield = put({
      type: "app/set_account",
      payload: { address: mockAddress },
    })
    const mockedAction = {
      type: "app/login",
      payload: { address: mockAddress },
    }
    expect(generator.next(mockedAction).value).toStrictEqual(expectedYield)
  })

  it("should fork background tasks", () => {
    const expectedYield = fork(runTasks)
    expect(generator.next().value).toStrictEqual(expectedYield)
  })

  it("should wait for stop action", () => {
    expect(generator.next(mockTask).value).toStrictEqual(take("app/logout"))
  })

  it("should cancel background task", () => {
    const mockedAction = { type: "app/logout" }
    expect(generator.next(mockedAction).value).toStrictEqual(cancel(mockTask))
  })

  it("should dispatch action that restets the store", () => {
    expect(generator.next().value).toStrictEqual(
      put({ type: "app/reset_store" })
    )
  })
})

describe("Test account switching saga step by step", () => {
  let generator = null
  const mockTask = createMockTask()
  const mockAddress = "0x0"

  beforeAll(() => {
    generator = runTasks()
  })

  it("should fork all sagas", () => {
    const expectedAllYield = all(loginRequiredSagas.map(fork))
    expect(generator.next().value).toStrictEqual(expectedAllYield)
  })

  it("should wait for stop action", () => {
    const expectedTakeYield = take("app/account_changed")
    expect(generator.next(mockTask).value).toStrictEqual(expectedTakeYield)
  })

  it("should cancel tasks", () => {
    const mockedAction = {
      type: "app/account_changed",
      payload: { address: mockAddress },
    }

    const expectedCancelYield = cancel(mockTask)
    expect(generator.next(mockedAction).value).toStrictEqual(
      expectedCancelYield
    )
  })

  it("should dispatch reset store action", () => {
    const expectedPutYield = put({ type: "app/reset_store" })
    expect(generator.next().value).toStrictEqual(expectedPutYield)
  })

  it("should disptach set account action", () => {
    const expectedPutYield = put({
      type: "app/set_account",
      payload: { address: mockAddress },
    })
    expect(generator.next().value).toStrictEqual(expectedPutYield)
  })
})
