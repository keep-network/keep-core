import React from "react"
import { withTimeline } from "../withTimeline"
import {
  KEEP_TO_T_EXCHANGE_RATE_IN_WEI,
  STAKE_ON_THRESHOLD_TIMELINE_STEPS,
} from "../../../constants/constants"
import { StakeOnThresholdTimeline } from "./components"
import { ModalBody, ModalFooter } from "../Modal"
import TokenAmount from "../../TokenAmount"
import OnlyIf from "../../OnlyIf"
import List from "../../List"
import Button, { SubmitButton } from "../../Button"
import { shortenAddress } from "../../../utils/general.utils"
// import { useDispatch } from "react-redux"
import { ViewInBlockExplorer } from "../../ViewInBlockExplorer"
import * as Icons from "./../../Icons"
import { KEEP, ThresholdToken } from "../../../utils/token.utils"
import BigNumber from "bignumber.js"

const StakeOnThresholdComponent = ({
  bodyTitle,
  keepAmount,
  operator,
  beneficiary,
  authorizer,
  transactionHash = false,
  onClose,
}) => {
  // const dispatch = useDispatch()

  const thresholdTokenAmount = (amount) => {
    const floatingPointDivisor = new BigNumber(10).pow(15)
    const amountInBN = new BigNumber(amount)
    const wrappedRemainder = amountInBN.modulo(floatingPointDivisor)
    const convertibleAmount = amountInBN.minus(wrappedRemainder)

    return convertibleAmount
      .multipliedBy(KEEP_TO_T_EXCHANGE_RATE_IN_WEI)
      .dividedBy(floatingPointDivisor)
      .toString()
  }

  return (
    <>
      <ModalBody>
        <h3 className="mb-1">
          {bodyTitle
            ? bodyTitle
            : transactionHash
            ? "Set up a PRE node now to be eligible for earning monthly rewards."
            : "You are about to stake on threshold:"}
        </h3>
        <OnlyIf condition={transactionHash}>
          <a
            href={"https://google.com"}
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
            iconMeasurements={{ width: 20, height: 20 }}
            withMetricSuffix
          />
          <Icons.ArrowsRight className={"ml-1 mr-1"} />
          <TokenAmount
            token={ThresholdToken}
            amount={thresholdTokenAmount(keepAmount)}
            amountClassName="h3 text-black"
            symbolClassName="h3 text-black"
            withIcon
            iconMeasurements={{ width: 20, height: 20 }}
            iconProps={{ className: "" }}
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
          ) : (
            "This requires two transactions – an authorization and a confirmation."
          )}
        </p>
        <List className="mt-2">
          <List.Content className="text-grey-50">
            <List.Item className="flex row center">
              <span className="mr-a">Operator</span>
              <span>{shortenAddress(operator)}</span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Beneficiary</span>
              <span>{shortenAddress(beneficiary)}</span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Authorizer</span>
              <span>{shortenAddress(authorizer)}</span>
            </List.Item>
            <List.Item className="flex row center">
              <span className="mr-a">Exchange Rate</span>
              <span>
                1 KEEP ={" "}
                {ThresholdToken.displayAmountWithSymbol(
                  thresholdTokenAmount(KEEP.fromTokenUnit(1)),
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
          <SubmitButton className="btn btn-primary btn-lg mr-2" type="submit">
            stake on t
          </SubmitButton>
        </OnlyIf>
        <Button
          className={`btn btn-${
            transactionHash ? "secondary btn-lg" : "unstyled text-link"
          }`}
          onClick={onClose}
        >
          {transactionHash ? "Close" : "Cancel"}
        </Button>
      </ModalFooter>
    </>
  )
}

const StakeOnThreshold = withTimeline({
  title: "Stake on Threshold",
  timelineComponent: StakeOnThresholdTimeline,
  timelineProps: {
    step: STAKE_ON_THRESHOLD_TIMELINE_STEPS.NONE,
  },
})(StakeOnThresholdComponent)

export default StakeOnThreshold
