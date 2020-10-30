import { combineReducers } from "redux"
import messages from "./messages"
import copyStake from "./copy-stake"
import keepTokenBalance from "./keep-balance"

const reducers = combineReducers({ messages, copyStake, keepTokenBalance })

export default reducers
