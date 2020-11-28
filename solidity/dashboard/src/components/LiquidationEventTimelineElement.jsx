import React, { useContext, useCallback } from "react"
import { Link } from "react-router-dom"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { Web3Context } from "./WithWeb3Context"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { liquidationService } from "../services/tbtc-liquidation.service"
import { VerticalTimelineElement } from "react-vertical-timeline-component"
import { satsToTBtcViaWeitoshi } from "../utils/token.utils"
import { useShowMessage, messageType } from "./Message"
import { DEPOSIT_STATES } from "../constants/constants"
import moment from "moment"
import * as Icons from "./Icons"

export const LiquidationEventTimelineElement = ({ event, isLoading }) => {
  const web3Context = useContext(Web3Context)
  const showMessage = useShowMessage()

  const { returnValues, depositStatusObj } = event

  const getLiquidationDueTimestamp = (eventTimestamp, depositStatusId) => {
    const originalEventTimestamp = moment(eventTimestamp)
    const dso = DEPOSIT_STATES.find(
      (dso) => dso.depositStatusId === depositStatusId
    )
    switch (dso.name) {
      case "AWAITING_WITHDRAWAL_SIGNATURE":
        return originalEventTimestamp.add(7200, "seconds")
      case "AWAITING_WITHDRAWAL_PROOF":
        return originalEventTimestamp.add(21600, "seconds")
      default:
        return originalEventTimestamp
    }
  }

  const isTimeToProceed = moment().isAfter(
    getLiquidationDueTimestamp(
      event.timestamp,
      depositStatusObj.depositStatusId
    )
  )

  const actionableStates = DEPOSIT_STATES.filter((dst) =>
    [5, 6, 8].includes(dst.depositStatusId)
  )
  const isActionable = actionableStates.includes(depositStatusObj)

  const badgeStatusPerDepositStatus = () => {
    const completedStates = DEPOSIT_STATES.filter((dst) =>
      [7, 11].includes(dst.depositStatusId)
    )
    if (isActionable) {
      if (isTimeToProceed) return BADGE_STATUS.PENDING
      return BADGE_STATUS.ACTIVE
    } else if (completedStates.includes(depositStatusObj)) {
      return BADGE_STATUS.COMPLETE
    } else {
      return BADGE_STATUS.DISABLED
    }
  }

  const onLiquidateClick = useCallback(async () => {
    const getLiquidationStateFn = (depositStatus) => {
      const dso = DEPOSIT_STATES.find(
        (dso) => dso.depositStatusId === depositStatus
      )
      switch (dso.name) {
        case "AWAITING_WITHDRAWAL_SIGNATURE":
          return liquidationService.depositNotifySignatureTimeout
        case "AWAITING_WITHDRAWAL_PROOF":
          return liquidationService.depositNotifyRedemptionProofTimedOut
        default:
          return null
      }
    }
    try {
      await getLiquidationStateFn(depositStatusObj.depositStatusId)(
        web3Context,
        returnValues._depositContractAddress
      )
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Liquidation transaction sent to your web wallet",
      })
    } catch (error) {
      showMessage({
        type: messageType.ERROR,
        title: "Liquidation action has failed",
        content: error.message,
      })
      throw error
    }
  }, [returnValues, depositStatusObj, showMessage, web3Context])

  const DepositAddress = React.memo(({ address, label }) => (
    <h6 className="text-grey-50" style={{ marginTop: "0.5rem" }}>
      {label}&nbsp;
      <AddressShortcut
        address={address}
        classNames="h6 text-normal text-grey-50"
      />
    </h6>
  ))

  const element = isLoading ? null : (
    <VerticalTimelineElement
      className="vertical-timeline-element--work"
      contentStyle={{ background: "#FFF", color: "#4C4C4C" }}
      contentArrowStyle={{ borderRight: "7px solid  rgb(33, 150, 243)" }}
      date={event.timestamp.toString()}
      iconStyle={{ background: "rgb(33, 150, 243)", color: "#fff" }}
      icon={<Icons.Glossary />}
    >
      <h2>Event: {event.event}</h2>
      <StatusBadge
        status={badgeStatusPerDepositStatus()}
        className="self-start"
        text={event.depositStatusObj.name}
        onlyIcon={badgeStatusPerDepositStatus() === BADGE_STATUS.COMPLETE}
      />
      <DepositAddress
        address={returnValues._depositContractAddress}
        label={"Deposit Address:"}
      />
      <p>
        Deposit size:{" "}
        {satsToTBtcViaWeitoshi(event.returnValues._utxoValue).toString()} tBTC
      </p>
      <span>
        Liquidate Not Before:{" "}
        {getLiquidationDueTimestamp(
          event.timestamp,
          event.depositStatusObj.depositStatusId
        ).toString()}
      </span>
      <br />
      <SubmitButton
        onSubmitAction={onLiquidateClick}
        className="btn btn-primary btn-sm"
        style={{ marginLeft: "auto" }}
        disabled={!isTimeToProceed || !isActionable}
      >
        Liquidate
      </SubmitButton>
      {/* <Link to={`/liquidations/${returnValues._depositContractAddress}`} className="btn btn-secondary mt-2"> */}
      <Link to={`/liquidations/${returnValues._depositContractAddress}`} className="btn btn-secondary btn-sm">
        View Liquidation
      </Link>
    </VerticalTimelineElement>
  )

  return element
}

export default React.memo(LiquidationEventTimelineElement)
