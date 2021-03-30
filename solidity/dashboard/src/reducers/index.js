import { combineReducers } from "redux"
import messages from "./messages"
import copyStake from "./copy-stake"
import keepTokenBalance from "./keep-balance"
import staking from "./staking"
import tokenGrants from "./token-grant"
import rewards from "./rewards"
import liquidityRewards from "./liquidity-rewards"
import notificationsData from "./notifications-data"
import operator from "./operator"

const app = (state = { address: null }, action) => {
  switch (action.type) {
    case "app/set_account":
      return { ...state, address: action.payload.address }

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
  notificationsData,
  operator,
})

const rootReducer = (state, action) => {
  if (action.type === "app/reset_store") {
    state = undefined
  }

  return reducers(state, action)
}

export default rootReducer
