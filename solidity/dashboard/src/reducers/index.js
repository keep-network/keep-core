import { combineReducers } from "redux"
import messages from "./messages"
import copyStake from "./copy-stake"
import keepTokenBalance from "./keep-balance"
import staking from "./staking"
import tokenGrants from "./token-grant"
import rewards from "./rewards"
import liquidityRewards from "./liquidity-rewards"

const reducers = combineReducers({
  messages,
  copyStake,
  keepTokenBalance,
  staking,
  tokenGrants,
  rewards,
  liquidityRewards,
})

const rootReducer = (state, action) => {
  if (action.type === "RESET_APP") {
    state = undefined
  }

  return reducers(state, action)
}

export default rootReducer
