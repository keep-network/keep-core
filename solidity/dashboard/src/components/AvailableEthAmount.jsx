import React from "react"
import { colors } from "../constants/colors"

const styles = {
  ethAmount: {
    textAlign: "right",
    padding: "0 1rem",
    paddingLeft: "2rem",
    borderRadius: "100px",
    border: `1px solid ${colors.grey20}`,
    backgroundColor: `${colors.grey10}`,
    width: "90px",
    overflowX: "auto",
  },
}

const AvailableEthAmount = React.memo(({ availableETH }) => {
  return (
    <div className="flex row center">
      <span className="text-big text-grey-70" style={styles.ethAmount}>
        {availableETH}
      </span>
      <span className="text-grey-60">&nbsp;ETH</span>
    </div>
  )
})

export default React.memo(AvailableEthAmount)
