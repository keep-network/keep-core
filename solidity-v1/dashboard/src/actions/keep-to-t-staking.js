import { THRESHOLD_AUTHORIZED, THRESHOLD_STAKED_TO_T } from "./index"

export const THRESHOLD_STAKE_KEEP_EVENT_EMITTED =
  "threshold/stake_keep_event_emitted"
export const STAKE_KEEP_TO_T = "threshold/stake_keep_to_t"

export const thresholdStakeKeepEventEmitted = (event) => {
  return {
    type: THRESHOLD_STAKE_KEEP_EVENT_EMITTED,
    payload: { event },
  }
}

export const thresholdContractAuthorized = (operatorAddress) => {
  return {
    type: THRESHOLD_AUTHORIZED,
    payload: { operatorAddress },
  }
}

export const stakedToT = (operatorAddress) => {
  return {
    type: THRESHOLD_STAKED_TO_T,
    payload: { operatorAddress },
  }
}

export const stakeKeepToT = (data, meta) => {
  return {
    type: STAKE_KEEP_TO_T,
    payload: {
      operator: data.operatorAddress,
      isAuthorized: data.isAuthorized,
    },
    meta,
  }
}
