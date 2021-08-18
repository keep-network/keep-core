import React from "react"
import Timeline from "../Timeline"
import { colors } from "../../constants/colors"
import Chip from "../Chip"
import OnlyIf from "../OnlyIf"

export const MODAL_WITH_TIMELINE_STEPS = {
  DEPOSITED_TOKENS: 1,
  WITHDRAW_DEPOSIT: 2,
  COOLDOWN: 3,
  CLAIM_TOKENS: 4,
}

const ModalWithTimeline = ({
  children,
  className = "",
  step = MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS,
  withDescription = false,
}) => {
  return (
    <div className={`modal-with-timeline__content-container ${className}`}>
      <div className={"modal-with-timeline-modal__info"}>{children}</div>
      <div className={"modal-with-timeline__timeline-container"}>
        <h4>Overview</h4>
        <Timeline className={"modal-with-timeline__timeline"}>
          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot>1</Timeline.BreakpointDot>
              <Timeline.BreakpointLine
                active={step > MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS}
              />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Timeline.ElementDefaultCard>
                <h4>Deposit your tokens</h4>
                <OnlyIf
                  condition={
                    step === MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS &&
                    withDescription
                  }
                >
                  <span className="text-grey-60">
                    There is no minimum KEEP amount for your deposit and no
                    minimum time lock.
                  </span>
                </OnlyIf>
              </Timeline.ElementDefaultCard>
            </Timeline.Content>
          </Timeline.Element>

          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot
                active={step > MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS}
              >
                2
              </Timeline.BreakpointDot>
              <Timeline.BreakpointLine
                active={step > MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT}
              />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Timeline.ElementDefaultCard
                active={step > MODAL_WITH_TIMELINE_STEPS.DEPOSITED_TOKENS}
              >
                <h4>Withdraw deposit</h4>
                <OnlyIf
                  condition={
                    step === MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT &&
                    withDescription
                  }
                >
                  <span className="text-grey-60">
                    Withdrawing requires two steps. First, there is a 21 day
                    cooldown. Second, after 21 days your tokens will be
                    available to claim in the dashboard.
                  </span>
                </OnlyIf>
              </Timeline.ElementDefaultCard>
            </Timeline.Content>
          </Timeline.Element>
          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot
                lineBreaker
                lineBreakerColor={
                  step > MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT
                    ? `violet-80`
                    : `grey-40`
                }
                style={{ backgroundColor: colors.brandViolet10 }}
              />
              <Timeline.BreakpointLine
                active={step > MODAL_WITH_TIMELINE_STEPS.COOLDOWN}
              />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Chip
                text="21 day cooldown"
                color={
                  step > MODAL_WITH_TIMELINE_STEPS.WITHDRAW_DEPOSIT
                    ? "strong"
                    : "subtle"
                }
                size="big"
              />
            </Timeline.Content>
          </Timeline.Element>

          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot
                active={step > MODAL_WITH_TIMELINE_STEPS.COOLDOWN}
              >
                3
              </Timeline.BreakpointDot>
              <Timeline.BreakpointLine
                active={step > MODAL_WITH_TIMELINE_STEPS.CLAIM_TOKENS}
              />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Timeline.ElementDefaultCard
                active={step > MODAL_WITH_TIMELINE_STEPS.COOLDOWN}
              >
                <h4>Claim tokens</h4>
                <OnlyIf
                  condition={
                    step === MODAL_WITH_TIMELINE_STEPS.CLAIM_TOKENS &&
                    withDescription
                  }
                >
                  <span className="text-grey-60">
                    There’s a 2 day claim window to claim your tokens and
                    rewards.
                  </span>
                </OnlyIf>
              </Timeline.ElementDefaultCard>
            </Timeline.Content>
          </Timeline.Element>
        </Timeline>
      </div>
    </div>
  )
}

export default ModalWithTimeline
