import { take, call, fork, cancel, cancelled } from "redux-saga/effects"

export function takeOnlyOnce(pattern, taskIdentifier, fn, ...args) {
  return fork(function* () {
    const tasks = {}

    while (true) {
      const action = yield take(pattern)
      const id = taskIdentifier(action)

      if (tasks[id]) {
        continue
      }

      tasks[id] = yield fork(function* () {
        try {
          yield call(fn, ...args.concat(action))
        } catch (error) {
          yield cancel()
        } finally {
          if (yield cancelled()) {
            delete tasks[id]
          }
        }
      })
    }
  })
}
