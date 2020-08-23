import React, { useCallback, useState } from "react"
import LiquidationEventTimelineElement from "../components/LiquidationEventTimelineElement"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { TBTC_SYSTEM_CONTRACT_NAME } from "../constants/constants"
import PageWrapper from "../components/PageWrapper"
import { liquidationService } from "../services/tbtc-liquidation.service"
import { useFetchData } from "../hooks/useFetchData"
import { LoadingOverlay } from "../components/Loadable"
import DelegatedTokensSkeleton from "../components/skeletons/DelegatedTokensSkeleton"
import {
  VerticalTimeline,
  VerticalTimelineElement,
} from "react-vertical-timeline-component"
import "react-vertical-timeline-component/style.min.css"
import * as Icons from "../components/Icons"

const LiquidationsPage = (props) => {
  const [redemptionRequestedEvents, setRedemptionReqEvData, refreshRedemptionReqEvData] = useFetchData(
    liquidationService.getPastRedemptionRequestedEvents,
    {}
  )
  let { isFetching: isRedemptionReqEvFetching, data: redemptionRequestEvData } = redemptionRequestedEvents

  const [redemptionSignatureEvents, setRedemptionSigEvData, refreshRedemptionSigEvData] = useFetchData(
    liquidationService.getPastRedemptionSignatureEvents,
    {}
  )
  let { isFetching: isRedemptionSigEvFetching, data: redemptionSignatureEvData } = redemptionSignatureEvents

  // eslint-disable-next-line
  const [courtesyCalledEvents, setCourtesyCallEvData, refreshCourtesyCalledEvData] = useFetchData(
    liquidationService.getPastCourtesyCalledEvents,
    {}
  )
  let { isFetching: isCourtesyCallEvFetching, data: courtesyCalledEvData } = courtesyCalledEvents
  let allEventData = []
  let anyPastEventCallFetching = isRedemptionReqEvFetching || isRedemptionSigEvFetching || isCourtesyCallEvFetching

  if (!anyPastEventCallFetching) {
    allEventData = [...redemptionRequestEvData.filter((ev) => ev.isAlignedToDeposit),
      ...redemptionSignatureEvData.filter((ev) => ev.isAlignedToDeposit),
      ...courtesyCalledEvData.filter((ev) => ev.isAlignedToDeposit)]
      allEventData = allEventData.sort((a, b) => b.timestamp - a.timestamp)
  }

  // control number of events displayed for "load older" feature. Initially 10
  const [eventCountHorizon, setEventCountHorizon] = useState(10)

  const increaseEventHorizon = () => {
    setEventCountHorizon(eventCountHorizon + 10)
  }

  const redemptionRequestsUpdated = useCallback(
    (latestEvent) => {
      setRedemptionReqEvData([latestEvent, ...redemptionRequestedEvents.data])
    },
    [redemptionRequestedEvents, setRedemptionReqEvData]
  )
  const subscribeToRedemptionRequestedCallback = (event) => {
    redemptionRequestsUpdated(event)
  }
  useSubscribeToContractEvent(
    TBTC_SYSTEM_CONTRACT_NAME,
    "RedemptionRequested",
    subscribeToRedemptionRequestedCallback
  )

  const redemptionSignaturesUpdated = useCallback(
    (latestEvent) => {
      setRedemptionSigEvData([latestEvent, ...redemptionSignatureEvents.data])
    },
    [redemptionSignatureEvents, setRedemptionSigEvData]
  )
  const subscribeToRedemptionSignatureCallback = (event) => {
    redemptionSignaturesUpdated(event)
  }
  useSubscribeToContractEvent(
    TBTC_SYSTEM_CONTRACT_NAME,
    "GotRedemptionSignature",
    subscribeToRedemptionSignatureCallback
  )

  // TODO: Include useSubscribeToContractEvent for CourtesyCalled. 

  const refreshAllPastEvData = () => {
    refreshRedemptionReqEvData()
    refreshRedemptionSigEvData()
    refreshCourtesyCalledEvData()
  }

  const getLatestEvents = () => {
    const evList = allEventData.sort((a, b) => b.timestamp - a.timestamp)
    const componentList = []
    for (
      let index = 0;
      index < evList.length && index < eventCountHorizon + 1;
      index++
    ) {
      const element = (
        <LiquidationEventTimelineElement
          isLoading={anyPastEventCallFetching}
          key={evList[index].transactionHash}
          event={evList[index]}
        />
      )
      componentList.push(element)
    }
    return componentList
  }

  return (
    <PageWrapper title="Liquidations">
      <LoadingOverlay
        isFetching={anyPastEventCallFetching}
        skeletonComponent={<DelegatedTokensSkeleton />}
      ></LoadingOverlay>
      <VerticalTimeline>
        <VerticalTimelineElement
          iconOnClick={() => {
            refreshAllPastEvData()
          }}
          iconClassName="vertical-timeline-element-icon--button"
          icon={<Icons.Load />}
          iconStyle={{ background: "rgb(33, 150, 243)", color: "#fff" }}
          date="load newer"
        />
        {anyPastEventCallFetching ? null : getLatestEvents()}
        <VerticalTimelineElement
          iconOnClick={() => {
            increaseEventHorizon()
          }}
          iconClassName="vertical-timeline-element-icon--button"
          icon={<Icons.Load />}
          date="load older"
        />
      </VerticalTimeline>
    </PageWrapper>
  )
}

export default LiquidationsPage
