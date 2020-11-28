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
  const [redemptionRequestedEvents, setData, refreshData] = useFetchData(
    liquidationService.getPastRedemptionRequestedEvents,
    {}
  )
  let { isFetching, data } = redemptionRequestedEvents

  if (data && data.length >= 0) {
    data = data.filter((ev) => ev.isAlignedToDeposit)
    data = data.sort((a, b) => b.timestamp - a.timestamp)
  }

  // control number of events displayed for "load older" feature. Initially 10
  const [eventCountHorizon, setEventCountHorizon] = useState(10)

  const increaseEventHorizon = () => {
    setEventCountHorizon(eventCountHorizon + 10)
  }

  const redemptionRequestsUpdated = useCallback(
    (latestEvent) => {
      setData([latestEvent, ...redemptionRequestedEvents.data])
    },
    [redemptionRequestedEvents, setData]
  )

  const subscribeToRedemptionRequestedCallback = (event) => {
    redemptionRequestsUpdated(event)
  }

  useSubscribeToContractEvent(
    TBTC_SYSTEM_CONTRACT_NAME,
    "RedemptionRequested",
    subscribeToRedemptionRequestedCallback
  )

  const getLatestEvents = () => {
    const evList = data.sort((a, b) => b.timestamp - a.timestamp)
    const componentList = []
    for (
      let index = 0;
      index < evList.length && index < eventCountHorizon + 1;
      index++
    ) {
      const element = (
        <LiquidationEventTimelineElement
          isLoading={isFetching}
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
        isFetching={isFetching}
        skeletonComponent={<DelegatedTokensSkeleton />}
      ></LoadingOverlay>
      <VerticalTimeline>
        <VerticalTimelineElement
          iconOnClick={() => {
            refreshData()
          }}
          iconClassName="vertical-timeline-element-icon--button"
          icon={<Icons.Load />}
          iconStyle={{ background: "rgb(33, 150, 243)", color: "#fff" }}
          date="load newer"
        />
        {isFetching ? null : getLatestEvents()}
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
