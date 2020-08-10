import { createStore, applyMiddleware } from "redux"
import createSagaMiddleware from "redux-saga"
import rootReducer from "./reducers"
import rootSaga from "./sagas"
import { Web3Loaded, ContractsLoaded } from "./contracts"

const sagaMiddleware = createSagaMiddleware({
  context: {
    web3: Web3Loaded,
    contracts: ContractsLoaded,
  },
})

const store = createStore(rootReducer, applyMiddleware(sagaMiddleware))

sagaMiddleware.run(rootSaga)

export default store
