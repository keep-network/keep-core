import { combineReducers } from "redux"
import messages from "./messages"
import copyStake from "./copy-stake"
import keepTokenBalance from "./keep-balance"
import staking from "./staking"
import tokenGrants from "./token-grant"
import rewards from "./rewards"
import liquidityRewards from "./liquidity-rewards"

const app = (state = { isReady: true }, action) => {
  switch (action.type) {
    case "app/reinitialization":
      return { ...state, isReady: false }

    default:
      return state
  }
}

const reducers = combineReducers({
  messages,
  copyStake,
  keepTokenBalance,
  staking,
  tokenGrants,
  rewards,
  liquidityRewards,
  app,
})

const rootReducer = (state, action) => {
  if (action.type === "app/reset_store") {
    state = undefined
  }

  return reducers(state, action)
}

export default rootReducer
