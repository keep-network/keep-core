import React from "react"
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  ModalCloseButton,
} from "../Modal"
import Timeline from "../../Timeline"
import { colors } from "../../../constants/colors"
import Chip from "../../Chip"
import OnlyIf from "../../OnlyIf"

export const withTimeline =
  ({ title, step, withDescription }) =>
  (WrappedModalContent) => {
    return (props) => {
      return (
        <Modal isOpen onClose={props.onClose} size="xl">
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>{title}</ModalHeader>
            <ModalCloseButton />
            <div className="modal-with-timeline__content-wrapper">
              <CovPoolTimeline step={step} withDescription={withDescription} />
              <WrappedModalContent {...props} />
            </div>
          </ModalContent>
        </Modal>
      )
    }
  }

export const MODAL_WITH_TIMELINE_STEPS = {
  DEPOSITED_TOKENS: 1,
  WITHDRAW_DEPOSIT: 2,
  COOLDOWN: 3,
  CLAIM_TOKENS: 4,
}

const CovPoolTimeline = ({ step, withDescription = false }) => {
  return (
    <div className="modal__timeline">
      <h4>Timeline</h4>
      <Timeline>
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
                  No min KEEP amount for your deposit.
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
                  Withdrawing requires two steps. First, there is a&nbsp;
                  <strong>21 day cooldown</strong>. Second, after 21 days your
                  tokens will be available to claim in the dashboard.
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
                  There’s a 2 day claim window to claim your tokens and rewards.
                </span>
              </OnlyIf>
            </Timeline.ElementDefaultCard>
          </Timeline.Content>
        </Timeline.Element>
      </Timeline>
    </div>
  )
}
