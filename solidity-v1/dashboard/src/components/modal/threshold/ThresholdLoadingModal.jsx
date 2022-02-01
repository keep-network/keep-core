import React from "react"
import { ModalBody } from "../Modal"
import { StakeOnThresholdTimeline } from "./components"
import { STAKE_ON_THRESHOLD_TIMELINE_STEPS } from "../../../constants/constants"
import { KeepLoadingIndicator } from "../../Loadable"
import { withTimeline } from "../withTimeline"

const ThresholdLoadingModal = ({ text }) => {
  return (
    <>
      <ModalBody className={"threshold-loading-modal-body"}>
        <KeepLoadingIndicator />
        <p>{text}</p>
      </ModalBody>
    </>
  )
}

export const ThresholdAuthorizationLoadingModal = withTimeline({
  title: "Sign Authorization (1/2)",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.AUTHORIZE_CONTRACT,
  },
})(ThresholdLoadingModal)

export const ThresholdStakeConfirmationLoadingModal = withTimeline({
  title: "Confirm Stake (2/2)",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.CONFIRM_STAKE,
  },
})(ThresholdLoadingModal)
