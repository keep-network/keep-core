import { createStore, applyMiddleware } from "redux"
import createSagaMiddleware from "redux-saga"
import { takeOnlyOnce } from "../../sagas/effects"

describe("Test takeOnlyOnce- custom redux-saga effect", () => {
  const taskIdentifier = jest.fn((action) => action.payload.id)
  let callCounter = []
  let worker = jest.fn((arg1, arg2, action) =>
    callCounter.push([arg1, arg2, action])
  )
  let middleware
  let store
  let mainTask

  beforeEach(() => {
    function* root() {
      yield takeOnlyOnce("test/pattern", taskIdentifier, worker, "a1", "a2")
    }

    callCounter = []
    middleware = createSagaMiddleware()
    store = createStore(() => {}, applyMiddleware(middleware))
    mainTask = middleware.run(root)
  })

  it("should only call worker once if id is the same", async () => {
    const dispatcher = Promise.resolve()
      .then(() => {
        store.dispatch({
          type: "test/pattern",
          payload: { id: 10 },
        })
      })
      .then(() =>
        store.dispatch({
          type: "test/pattern",
          payload: { id: 10 },
        })
      )
      .then(() => {
        mainTask.cancel()
      })

    await Promise.all([mainTask.toPromise(), dispatcher]).then(() => {
      expect(taskIdentifier).toHaveBeenCalledWith({
        type: "test/pattern",
        payload: { id: 10 },
      })

      expect(taskIdentifier).toHaveBeenCalledWith({
        type: "test/pattern",
        payload: { id: 10 },
      })

      // Validate passed arguments to the worker.
      expect(callCounter.length).toBe(1)
      expect(callCounter[0][0]).toBe("a1")
      expect(callCounter[0][1]).toBe("a2")
      expect(callCounter[0][2]).toStrictEqual({
        type: "test/pattern",
        payload: { id: 10 },
      })
    })
  })

  it("should call worker if ids are different", async () => {
    const dispatcher = Promise.resolve()
      .then(() => {
        store.dispatch({
          type: "test/pattern",
          payload: { id: 10 },
        })
      })
      .then(() =>
        store.dispatch({
          type: "test/pattern",
          payload: { id: 20 },
        })
      )
      .then(() => {
        mainTask.cancel()
      })

    await Promise.all([mainTask.toPromise(), dispatcher]).then(() => {
      expect(taskIdentifier).toHaveBeenCalledWith({
        type: "test/pattern",
        payload: { id: 10 },
      })

      expect(taskIdentifier).toHaveBeenCalledWith({
        type: "test/pattern",
        payload: { id: 20 },
      })

      expect(callCounter.length).toBe(2)

      expect(callCounter[0][0]).toBe("a1")
      expect(callCounter[0][1]).toBe("a2")
      expect(callCounter[0][2]).toStrictEqual({
        type: "test/pattern",
        payload: { id: 10 },
      })

      expect(callCounter[1][0]).toBe("a1")
      expect(callCounter[1][1]).toBe("a2")
      expect(callCounter[1][2]).toStrictEqual({
        type: "test/pattern",
        payload: { id: 20 },
      })
    })
  })

  it("should handle request if the previous has failed", async () => {
    worker = jest
      .fn()
      .mockImplementationOnce(() => {
        throw new Error("Fake error")
      })
      .mockImplementationOnce(() => callCounter.push([arg1, arg2, action]))

    const dispatcher = Promise.resolve()
      .then(() => {
        store.dispatch({
          type: "test/pattern",
          payload: { id: 10 },
        })
      })
      .then(() =>
        store.dispatch({
          type: "test/pattern",
          payload: { id: 10 },
        })
      )
      .then(() => {
        mainTask.cancel()
      })

    await Promise.all([mainTask.toPromise(), dispatcher]).then(() => {
      expect(taskIdentifier).toHaveBeenCalledWith({
        type: "test/pattern",
        payload: { id: 10 },
      })

      expect(taskIdentifier).toHaveBeenCalledWith({
        type: "test/pattern",
        payload: { id: 10 },
      })

      expect(callCounter.length).toBe(1)
      expect(callCounter[0][0]).toBe("a1")
      expect(callCounter[0][1]).toBe("a2")
      expect(callCounter[0][2]).toStrictEqual({
        type: "test/pattern",
        payload: { id: 10 },
      })
    })
  })
})
