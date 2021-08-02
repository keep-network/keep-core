import React from "react"
import * as Icons from "../Icons"
import TokenAmount from "../TokenAmount"
import { KEEP } from "../../utils/token.utils"
import Divider from "../Divider"
import Button from "../Button"
import OnlyIf from "../OnlyIf"
import Timeline from "../Timeline";
import {colors} from "../../constants/colors";
import Chip from "../Chip";

const ClaimTokensModal = ({
  amount,
  submitBtnText,
  onBtnClick,
  onCancel,
  transactionFinished = false,
  transactionHash = "",
}) => {
  return (
    <div className={"claim-tokens-modal__content-container"}>
      <div className={"claim-tokens-modal__info"}>
        <OnlyIf condition={!transactionFinished}>
          <h3 className={"mb-1"}>You are about to claim:</h3>
        </OnlyIf>
        <OnlyIf condition={transactionFinished}>
          <h3>Success!</h3>
          <h4 className={"text-gray-70 mb-1"}>View your transaction here.</h4>
        </OnlyIf>
        <div className={"claim-tokens-modal__data"}>
          <TokenAmount
            amount={"20000"}
            wrapperClassName={"claim-tokens-modal__token-amount"}
            token={KEEP}
            withIcon
          />
          <div className={"claim-tokens-modal__data-row"}>
            <h4 className={"text-grey-50"}>Initial Withdrawal &nbsp;</h4>
            <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
          <div className={"claim-tokens-modal__data-row"}>
            <h4 className={"text-grey-50"}>Rewards earned &nbsp;</h4>
            <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
          <div className={"claim-tokens-modal__data-row"}>
            <h4 className={"text-grey-50"}>Wallet &nbsp;</h4>
            <h4 className={"claim-tokens-modal__data__value text-grey-70"}>
              1,000 KEEP
            </h4>
          </div>
        </div>
        <Divider style={{ margin: "0.5rem 0" }} />
        <div className="flex row center mt-2">
          <OnlyIf condition={!transactionFinished}>
            <Button
              className="btn btn-lg btn-primary"
              type="submit"
              disabled={false}
              onClick={onBtnClick}
            >
              {submitBtnText}
            </Button>
            <span onClick={onCancel} className="ml-1 text-link">
              Cancel
            </span>
          </OnlyIf>
          <OnlyIf condition={transactionFinished}>
            <Button
              className="btn btn-lg btn-secondary"
              disabled={false}
              onClick={onCancel}
            >
              Close
            </Button>
          </OnlyIf>
        </div>
      </div>
      <div className={"claim-tokens-modal__timeline-container"}>
        <h4>Overview</h4>
        <Timeline className={"claim-tokens-modal__timeline"}>
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

export default ClaimTokensModal
