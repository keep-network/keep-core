import { combineReducers } from "redux"
import messages from "./messages"
import copyStake from "./copy-stake"

const reducers = combineReducers({ messages, copyStake })

export default reducers
