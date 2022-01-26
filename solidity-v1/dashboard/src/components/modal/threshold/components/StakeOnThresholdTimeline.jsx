import React from "react"
import Timeline, { TIMELINE_ELEMENT_STATUS } from "../../../Timeline"
import Banner from "../../../Banner"
import { STAKE_ON_THRESHOLD_TIMELINE_STEPS } from "../../../../constants/constants"

const styles = {
  defaultCard: {
    wrapper: { padding: "1rem" },
    title: {
      margin: "0",
    },
    desc: {
      marginTop: "0.4rem",
    },
  },
  cooldownBanner: {
    desc: {
      marginBottom: 0,
    },
  },
}

export const StakeOnThresholdTimeline = ({ step }) => {
  return (
    <div className="modal__timeline">
      <h4>Overview</h4>
      <Timeline>
        <Timeline.Element>
          <Timeline.Breakpoint>
            <Timeline.BreakpointDot
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.NONE
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            >
              1
            </Timeline.BreakpointDot>
            <Timeline.BreakpointLine
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.AUTHORIZE_CONTRACT
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Timeline.ElementDefaultCard
              style={styles.defaultCard.wrapper}
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.NONE
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            >
              <h4 style={styles.defaultCard.title}>Authorize Contract</h4>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>

        <Timeline.Element>
          <Timeline.Breakpoint>
            <Timeline.BreakpointDot
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.AUTHORIZE_CONTRACT
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            >
              2
            </Timeline.BreakpointDot>
            <Timeline.BreakpointLine
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.CONFIRM_STAKE
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Timeline.ElementDefaultCard
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.AUTHORIZE_CONTRACT
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
              style={styles.defaultCard.wrapper}
            >
              <h4 style={styles.defaultCard.title}>Confirm Stake</h4>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>

        <Timeline.Element>
          <Timeline.Breakpoint>
            <Timeline.BreakpointDot
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.CONFIRM_STAKE
                  ? TIMELINE_ELEMENT_STATUS.SEMI_ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            >
              3
            </Timeline.BreakpointDot>
            <Timeline.BreakpointLine
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.SET_UP_PRE
                  ? TIMELINE_ELEMENT_STATUS.ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Timeline.ElementDefaultCard
              status={
                step > STAKE_ON_THRESHOLD_TIMELINE_STEPS.CONFIRM_STAKE
                  ? TIMELINE_ELEMENT_STATUS.SEMI_ACTIVE
                  : TIMELINE_ELEMENT_STATUS.INACTIVE
              }
              style={styles.defaultCard.wrapper}
            >
              <h4 style={styles.defaultCard.title}>Set up PRE</h4>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>
      </Timeline>
      <Banner className="mt-2">
        <Banner.Description
          style={styles.cooldownBanner.desc}
          className="text-black"
        >
          You will need to set up a PRE node in order to be eligible to earn
          rewards..
        </Banner.Description>
      </Banner>
    </div>
  )
}
