import React from "react"
import TokenAmount from "./TokenAmount"
import { ETH } from "../utils/token.utils"
import { colors } from "../constants/colors"

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

const AvailableEthAmount = React.memo(({ availableETHInWei }) => {
  return (
    <TokenAmount
      amount={availableETHInWei}
      token={ETH}
      amountClassName=""
      amountStyles={styles.ethAmount}
      symbolClassName="text-grey-60"
    />
  )
})

export default AvailableEthAmount
