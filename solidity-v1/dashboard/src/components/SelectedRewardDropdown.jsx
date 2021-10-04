import React from "react"
import { shortenAddress } from "../utils/general.utils"

const groupPublicKeyStyle = {
  marginRight: "20px",
  borderRadius: "100px",
  padding: "0.1rem",
}

const SelectedRewardDropdown = ({ groupReward }) => {
  return (
    <div className="flex row flex-1 wrap space-between center">
      <div className="text-smaller text-black">
        {`${groupReward.reward} ETH`}
      </div>
      <div className="bg-grey-10 text-grey-50" style={groupPublicKeyStyle}>
        <span className="text-label text-smaller">group key:</span>
        <span className="text-smaller">
          &nbsp;{shortenAddress(groupReward.groupPublicKey)}
        </span>
      </div>
    </div>
  )
}

export default React.memo(SelectedRewardDropdown)
