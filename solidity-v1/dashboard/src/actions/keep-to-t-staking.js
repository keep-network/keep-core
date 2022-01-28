export const THRESHOLD_STAKE_KEEP_EVENT_EMITTED =
  "threshold/stake_keep_event_emitted"
export const STAKE_KEEP_TO_T = "threshold/stake_keep_to_t"

export const thresholdStakeKeepEventEmitted = (event) => {
  return {
    type: THRESHOLD_STAKE_KEEP_EVENT_EMITTED,
    payload: { event },
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
