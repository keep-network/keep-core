import React from "react"
import * as Icons from "./Icons"
import BeaconRewardsHelper, { ECDSARewardsHelper } from "../utils/rewardsHelper"
import { KEEP } from "../utils/token.utils"
import { formatDate } from "../utils/general.utils"

const MinKeepInInterval = ({ label, minKeep }) => {
  return (
    <li className="flex row center mb-1">
      <Icons.Beacon width={12} height={12} />
      <span className="text-small ml-1">{label}</span>
      <span className="text-small text-grey-60 ml-a">{minKeep}</span>
    </li>
  )
}

const NextInterval = ({ nextIntervalStart, intervalsLeft }) => {
  return (
    <li>
      <div className="flex row center">
        <Icons.Time width={12} height={12} className="time-icon--black" />
        <span className="text-small ml-1">Next Interval</span>
        <span className="text-small text-grey-60 ml-a">
          {nextIntervalStart}
        </span>
      </div>
      <div
        style={{ marginLeft: "calc(12px + 1rem)" }}
        className="text-small text-grey-50"
      >
        {intervalsLeft}&nbsp;intervals left
      </div>
    </li>
  )
}

const BeaconRewardsDetails = ({ pastRewards }) => {
  const intervalsLeft =
    BeaconRewardsHelper.keepAllocationsInInterval.length -
    BeaconRewardsHelper.currentInterval -
    1

  const nextIntervalStart = formatDate(
    BeaconRewardsHelper.intervalStartOf(BeaconRewardsHelper.currentInterval + 1)
  )

  return (
    <>
      <h4 className="text-grey-70 mb-2">Rewards Details</h4>
      <ul>
        <MinKeepInInterval
          label="Min. Groups"
          minKeep={BeaconRewardsHelper.minimumKeepsPerInterval}
        />
        {/* <li className="flex row center mb-1">
          <Icons.KeepOutline
            width={12}
            height={12}
            className="keep-outline--black"
          />
          <span className="text-small ml-1">Past Rewards</span>
          <span className="text-small text-grey-60 ml-a">
            {displayAmount(pastRewards)}
          </span>
        </li> */}
        <NextInterval
          nextIntervalStart={nextIntervalStart}
          intervalsLeft={intervalsLeft}
        />
      </ul>
    </>
  )
}

const ECDSARewardsDetails = ({ pastRewards }) => {
  const intervalsLeft =
    ECDSARewardsHelper.intervals - ECDSARewardsHelper.currentInterval - 1

  const nextIntervalStart = formatDate(
    ECDSARewardsHelper.intervalStartOf(ECDSARewardsHelper.currentInterval + 1)
  )

  return (
    <>
      <h4 className="text-grey-70 mb-2">Rewards Details</h4>
      <ul>
        <MinKeepInInterval
          label="Min. Deposits"
          minKeep={ECDSARewardsHelper.minimumKeepsPerInterval}
        />
        <li className="flex row center mb-1">
          <Icons.KeepOutline
            width={12}
            height={12}
            className="keep-outline--black"
          />
          <span className="text-small ml-1">Past Rewards</span>
          <span className="text-small text-grey-60 ml-a">
            {`${KEEP.displayAmountWithSymbol(pastRewards)}`}
          </span>
        </li>
        <NextInterval
          nextIntervalStart={nextIntervalStart}
          intervalsLeft={intervalsLeft}
        />
      </ul>
    </>
  )
}

export { BeaconRewardsDetails, ECDSARewardsDetails }
