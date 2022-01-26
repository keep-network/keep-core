import React from "react"
import Timeline from "../../../Timeline"
import { colors } from "../../../../constants/colors"
import Chip from "../../../Chip"
import OnlyIf from "../../../OnlyIf"
import Banner from "../../../Banner"
import { COV_POOL_TIMELINE_STEPS } from "../../../../constants/constants"

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

export const CovPoolTimeline = ({ step, withDescription = false }) => {
  return (
    <div className="modal__timeline">
      <h4>Timeline</h4>
      <Timeline>
        <Timeline.Element>
          <Timeline.Breakpoint>
            <Timeline.BreakpointDot>1</Timeline.BreakpointDot>
            <Timeline.BreakpointLine
              active={step > COV_POOL_TIMELINE_STEPS.DEPOSITED_TOKENS}
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Timeline.ElementDefaultCard style={styles.defaultCard.wrapper}>
              <h4 style={styles.defaultCard.title}>Deposit your tokens</h4>
              <OnlyIf
                condition={
                  step === COV_POOL_TIMELINE_STEPS.DEPOSITED_TOKENS &&
                  withDescription
                }
              >
                <div className="text-grey-60" style={styles.defaultCard.desc}>
                  No min KEEP amount for your deposit.
                </div>
              </OnlyIf>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>

        <Timeline.Element>
          <Timeline.Breakpoint>
            <Timeline.BreakpointDot
              active={step > COV_POOL_TIMELINE_STEPS.DEPOSITED_TOKENS}
            >
              2
            </Timeline.BreakpointDot>
            <Timeline.BreakpointLine
              active={step > COV_POOL_TIMELINE_STEPS.WITHDRAW_DEPOSIT}
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Timeline.ElementDefaultCard
              active={step > COV_POOL_TIMELINE_STEPS.DEPOSITED_TOKENS}
              style={styles.defaultCard.wrapper}
            >
              <h4 style={styles.defaultCard.title}>Withdraw deposit</h4>
              <OnlyIf
                condition={
                  step === COV_POOL_TIMELINE_STEPS.WITHDRAW_DEPOSIT &&
                  withDescription
                }
              >
                <div className="text-grey-60" style={styles.defaultCard.desc}>
                  After cooldown, you will need to claim your tokens.
                </div>
              </OnlyIf>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>

        <Timeline.Element>
          <Timeline.Breakpoint>
            <Timeline.BreakpointDot
              lineBreaker
              lineBreakerColor={
                step > COV_POOL_TIMELINE_STEPS.WITHDRAW_DEPOSIT
                  ? `violet-80`
                  : `grey-40`
              }
              style={{ backgroundColor: colors.brandViolet10 }}
            />
            <Timeline.BreakpointLine
              active={step > COV_POOL_TIMELINE_STEPS.COOLDOWN}
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Chip
              text="21 day cooldown"
              color={
                step > COV_POOL_TIMELINE_STEPS.WITHDRAW_DEPOSIT
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
              active={step > COV_POOL_TIMELINE_STEPS.COOLDOWN}
            >
              3
            </Timeline.BreakpointDot>
            <Timeline.BreakpointLine
              active={step > COV_POOL_TIMELINE_STEPS.CLAIM_TOKENS}
            />
          </Timeline.Breakpoint>
          <Timeline.Content>
            <Timeline.ElementDefaultCard
              active={step > COV_POOL_TIMELINE_STEPS.COOLDOWN}
              style={styles.defaultCard.wrapper}
            >
              <h4 style={styles.defaultCard.title}>Claim tokens</h4>
              <OnlyIf
                condition={
                  step === COV_POOL_TIMELINE_STEPS.CLAIM_TOKENS &&
                  withDescription
                }
              >
                <div className="text-grey-60" style={styles.defaultCard.desc}>
                  There’s a 2 day claim window to claim your tokens and rewards.
                </div>
              </OnlyIf>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>
      </Timeline>
      <Banner className="mt-2">
        <Banner.Description
          style={styles.cooldownBanner.desc}
          className="text-black"
        >
          During cooldown, your funds will accumulate rewards but are also
          subject to an insurance event.
        </Banner.Description>
      </Banner>
    </div>
  )
}
