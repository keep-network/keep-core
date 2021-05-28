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

const AvailableBondingTokenAmount = React.memo(({ availableTokensInWei, availableTokens }) => {
  return (
    <>
      <Tooltip
        simple
        delay={0}
        triggerComponent={() => {
          return (
            <span className="text-big text-grey-70" style={styles.ethAmount}>
              {displayEthAmount(availableTokensInWei)}
            </span>
          )
        }}
      >
        {availableTokens}&nbsp;ERC20
      </Tooltip>
      <span className="text-grey-60">&nbsp;ERC20</span>
    </>
  )
})

export default AvailableBondingTokenAmount
