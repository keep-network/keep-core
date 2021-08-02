import React from "react";
import Timeline from "../Timeline";
import {colors} from "../../constants/colors";
import Chip from "../Chip";

const ModalWithTimeline = ({children, className = ""}) => {
  return (
    <div className={`modal-with-timeline__content-container ${className}`}>
      <div className={"modal-with-timeline-modal__info"}>
        {children}
      </div>
      <div className={"modal-with-timeline__timeline-container"}>
        <h4>Overview</h4>
        <Timeline className={"modal-with-timeline__timeline"}>
          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot>1</Timeline.BreakpointDot>
              <Timeline.BreakpointLine active />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Timeline.ElementDefaultCard>
                <h4 className="text-violet-80">Deposit your tokens</h4>
              </Timeline.ElementDefaultCard>
            </Timeline.Content>
          </Timeline.Element>

          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot>2</Timeline.BreakpointDot>
              <Timeline.BreakpointLine active />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Timeline.ElementDefaultCard>
                <h4 className="text-violet-80">Withdraw deposit</h4>
              </Timeline.ElementDefaultCard>
            </Timeline.Content>
          </Timeline.Element>
          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot
                lineBreaker
                lineBreakerColor="violet-80"
                style={{ backgroundColor: colors.brandViolet10 }}
              />
              <Timeline.BreakpointLine active />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Chip text="21 day cooldown" color="strong" size="big" />
            </Timeline.Content>
          </Timeline.Element>

          <Timeline.Element>
            <Timeline.Breakpoint>
              <Timeline.BreakpointDot>3</Timeline.BreakpointDot>
              <Timeline.BreakpointLine active />
            </Timeline.Breakpoint>
            <Timeline.Content>
              <Timeline.ElementDefaultCard>
                <h4 className="text-violet-80">Claim tokens</h4>
                <span className="text-grey-60">
                  There’s a 2 day claim window to claim your tokens and rewards.
                </span>
              </Timeline.ElementDefaultCard>
            </Timeline.Content>
          </Timeline.Element>
        </Timeline>
      </div>
    </div>
  )
}

export default ModalWithTimeline
