export const STAKE_KEEP_TO_T = "threshold/stake_keep_to_t"

export const stakeKeepToT = (operatorAddress, meta) => {
  return {
    type: STAKE_KEEP_TO_T,
    payload: {
      operator: operatorAddress,
    },
    meta,
  }
}
