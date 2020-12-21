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

export default reducers
