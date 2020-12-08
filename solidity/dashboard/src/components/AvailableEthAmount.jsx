import React from "react"
import Tooltip from "./Tooltip"
import { colors } from "../constants/colors"
import { displayEthAmount } from "../utils/ethereum.utils"

const styles = {
  ethAmount: {
    textAlign: "right",
    padding: "0 0.5rem",
    paddingLeft: "1.5rem",
    borderRadius: "100px",
    border: `1px solid ${colors.grey20}`,
    backgroundColor: `${colors.grey10}`,
  },
}

const AvailableEthAmount = React.memo(({ availableETHInWei, availableETH }) => {
  return (
    <>
      <Tooltip
        simple
        delay={0}
        triggerComponent={() => {
          return (
            <span className="text-big text-grey-70" style={styles.ethAmount}>
              {displayEthAmount(availableETHInWei)}
            </span>
          )
        }}
      >
        {availableETH}&nbsp;ETH
      </Tooltip>
      <span className="text-grey-60">&nbsp;ETH</span>
    </>
  )
})

export default AvailableEthAmount
