import { all, fork, take, cancel, put } from "redux-saga/effects"
import { createMockTask } from "@redux-saga/testing-utils"
import { expectSaga } from "redux-saga-test-plan"
import rootSaga from "../../sagas"
import * as messagesSaga from "../../sagas/messages"
import * as delegateStakeSaga from "../../sagas/staking"
import * as tokenGrantSaga from "../../sagas/token-grant"
import { watchSendTransactionRequest } from "../../sagas/web3"
import * as copyStakeSaga from "../../sagas/copy-stake"
import * as subscriptions from "../../sagas/subscriptions"
import * as keepTokenBalance from "../../sagas/keep-balance"
import * as rewards from "../../sagas/rewards"
import * as liquidityRewards from "../../sagas/liquidity-rewards"

describe("Test root saga", () => {
  it("should start correctly and handle dispatched actions", () => {
    const mockTasks = [fork(() => {}), fork(() => {})]

    return expectSaga(rootSaga)
      .provide([[all(), mockTasks]])
      .dispatch({ type: "app/set_account" })
      .dispatch({ type: "app/account_changed", payload: { address: "0x1" } })
      .put({ type: "app/reset_store" })
      .put({ type: "app/set_account", payload: { address: "0x1" } })
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

  afterAll(() => {
    generator.cancel()
  })

  it("should wait for start action", () => {
    const expectedYield = take("app/set_account")
    expect(generator.next().value).toStrictEqual(expectedYield)
  })

  it("should fork all sagas", () => {
    const expectedAllYield = all(
      [
        ...Object.values(messagesSaga),
        ...Object.values(delegateStakeSaga),
        watchSendTransactionRequest,
        ...Object.values(tokenGrantSaga),
        ...Object.values(copyStakeSaga),
        ...Object.values(subscriptions),
        ...Object.values(keepTokenBalance),
        ...Object.values(rewards),
        ...Object.values(liquidityRewards),
      ].map(fork)
    )
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
