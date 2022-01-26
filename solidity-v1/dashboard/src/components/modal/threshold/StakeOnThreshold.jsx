import React from "react"
import { withTimeline } from "../withTimeline"
import { STAKE_ON_THRESHOLD_TIMELINE_STEPS } from "../../../constants/constants"
import { StakeOnThresholdTimeline } from "./components"

const StakeOnThresholdComponent = () => {
  return <div>Stake on Threshold</div>
}

const StakeOnThreshold = withTimeline({
  title: "Stake on Threshold",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.NONE,
  },
})(StakeOnThresholdComponent)

export default StakeOnThreshold
