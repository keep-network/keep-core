import React from "react"
import { withTimeline } from "../withTimeline"
import {
  LINK,
  STAKE_ON_THRESHOLD_TIMELINE_STEPS,
} from "../../../constants/constants"
import { StakeOnThresholdTimeline } from "./components"
import { ModalBody, ModalFooter } from "../Modal"
import TokenAmount from "../../TokenAmount"
import OnlyIf from "../../OnlyIf"
import List from "../../List"
import Button from "../../Button"
import { shortenAddress } from "../../../utils/general.utils"
import { ViewInBlockExplorer } from "../../ViewInBlockExplorer"
import * as Icons from "./../../Icons"
import { KEEP, ThresholdToken } from "../../../utils/token.utils"
import BigNumber from "bignumber.js"
import { stakeKeepToT } from "../../../actions/keep-to-t-staking"
import { useDispatch } from "react-redux"
import { Keep } from "../../../contracts"

const StakeOnThresholdComponent = ({
  bodyTitle,
  keepAmount,
  operator,
  beneficiary,
  authorizer,
  isAuthorized,
  transactionHash = false,
  onClose,
}) => {
  const dispatch = useDispatch()

  return (
    <>
      <ModalBody>
        <h3 className="stake-on-threshold-modal__body-title mb-1">
          {bodyTitle
            ? bodyTitle
            : transactionHash
            ? "Set up a PRE node now to be eligible for earning monthly rewards."
            : "You're about to stake on Threshold:"}
        </h3>
        <OnlyIf condition={transactionHash}>
          <a
            href={LINK.setUpPRE}
            rel="noopener noreferrer"
            target="_blank"
            className={`btn btn-primary btn-semi-md mb-1`}
          >
            {"set up PRE"} ↗
          </a>
        </OnlyIf>
        <div className={"flex row center"}>
          <TokenAmount
            amount={keepAmount}
            amountClassName="h3 text-mint-100"
            symbolClassName="h3 text-mint-100"
            withIcon
            iconProps={{
              className: "keep-outline keep-outline--mint-80",
              width: 20,
              height: 20,
            }}
            withMetricSuffix
          />
          <Icons.ArrowsRight className={"ml-1 mr-1"} />
          <TokenAmount
            token={ThresholdToken}
            amount={Keep.keepToTStaking.toThresholdTokenAmount(keepAmount)}
            amountClassName="h3 text-black"
            symbolClassName="h3 text-black"
            withIcon
            iconProps={{
              className: "",
              width: 20,
              height: 20,
            }}
            withMetricSuffix
          />
        </div>
        <p className="mt-1 text-grey-70">
          {transactionHash ? (
            <>
              Your stake is confirmed! View your transaction&nbsp;
              <ViewInBlockExplorer
                type="tx"
                className="text-grey-70"
                id={transactionHash}
                text="here"
              />
              .
            </>
          ) : isAuthorized ? (
            "The contract for this stake is already authorized so only one transaction is needed - a confirmation"
          ) : (
            "This requires two transactions – an authorization and a confirmation."
          )}
        </p>
        <List className="mt-2">
          <List.Content className="text-grey-50">
            <List.Item className="flex row center">
              <span className="mr-a text-grey-50">Operator</span>
              <span className={"text-grey-70"}>{shortenAddress(operator)}</span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a text-grey-50">Beneficiary</span>
              <span className={"text-grey-70"}>
                {shortenAddress(beneficiary)}
              </span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a text-grey-50">Authorizer</span>
              <span className={"text-grey-70"}>
                {shortenAddress(authorizer)}
              </span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a text-grey-50">Exchange Rate</span>
              <span className={"text-grey-70"}>
                1 KEEP ={" "}
                {ThresholdToken.displayAmountWithSymbol(
                  Keep.keepToTStaking.toThresholdTokenAmount(
                    KEEP.fromTokenUnit(1)
                  ),
                  3,
                  (amount) =>
                    new BigNumber(amount).toFormat(3, BigNumber.ROUND_DOWN)
                )}
              </span>
            </List.Item>
          </List.Content>
        </List>
      </ModalBody>
      <ModalFooter>
        <OnlyIf condition={!transactionHash}>
          <Button
            className="btn btn-primary btn-lg mr-2"
            type="submit"
            onClick={() => {
              dispatch(
                stakeKeepToT({
                  operatorAddress: operator,
                  isAuthorized: isAuthorized,
                })
              )
            }}
          >
            stake on t
          </Button>
        </OnlyIf>
        <Button
          className={`btn btn-${
            transactionHash ? "secondary btn-lg" : "unstyled text-link"
          }`}
          onClick={onClose}
        >
          {transactionHash ? "close" : "cancel"}
        </Button>
      </ModalFooter>
    </>
  )
}

export const AuthorizeAndStakeOnThreshold = withTimeline({
  title: "Stake on Threshold",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.NONE,
  },
})(StakeOnThresholdComponent)

export const StakeOnThresholdWithoutAuthorization = withTimeline({
  title: "Stake on Threshold",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.AUTHORIZE_CONTRACT,
  },
})(StakeOnThresholdComponent)

export const StakeOnThresholdConfirmed = withTimeline({
  title: "Almost there...",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.SET_UP_PRE,
  },
})(StakeOnThresholdComponent)
