import React from "react"
import * as Icons from "./Icons"

const BeaconRewardsDetails = () => {
  return (
    <>
      <h4 className="text-grey-70 mb-2">Beacon Rewards Details</h4>
      <ul>
        <li className="flex row center mb-1">
          <Icons.Beacon width={12} height={12} />
          <span className="text-small ml-1">Min. Keep Group</span>
          <span className="text-small text-grey-60 ml-a">2</span>
        </li>
        <li className="flex row center mb-1">
          <Icons.Rewards width={12} height={12} />
          <span className="text-small ml-1">Est. Allocation</span>
          <span className="text-small text-grey-60 ml-a">10 KEEP</span>
        </li>
        <li className="flex row center mb-1">
          <Icons.KeepOutline
            width={12}
            height={12}
            className="keep-outline--black"
          />
          <span className="text-small ml-1">Past Rewards</span>
          <span className="text-small text-grey-60 ml-a">200 KEEP</span>
        </li>
        <li className="flex row center">
          <Icons.Time width={12} height={12} className="time-icon--black" />
          <span className="text-small ml-1">Next Interval</span>
          <span className="text-small text-grey-60 ml-a">11/15/2020</span>
        </li>
      </ul>
    </>
  )
}

export default BeaconRewardsDetails
