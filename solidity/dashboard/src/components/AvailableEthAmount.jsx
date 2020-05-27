import React from "react"
import { colors } from "../constants/colors"

const AvailableEthAmount = React.memo(({ availableETH }) => {
  return (
    <>
      <span
        className="text-big text-grey-70"
        style={{
          textAlign: "right",
          padding: "0.25rem 1rem",
          paddingLeft: "2rem",
          borderRadius: "100px",
          border: `1px solid ${colors.grey20}`,
          backgroundColor: `${colors.grey10}`,
        }}
      >
        {availableETH}
      </span>
      <span style={{ color: `${colors.grey60}` }}>&nbsp;ETH</span>
    </>
  )
})

export default React.memo(AvailableEthAmount)
