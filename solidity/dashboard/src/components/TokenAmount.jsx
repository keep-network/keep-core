import React from "react"
import * as Icons from "./Icons"
import { displayAmount, getNumberWithMetricSuffix } from "../utils/token.utils"

const TokenAmount = ({
  currencyIcon,
  amount,
  suffixClassName,
  amountClassName,
}) => {
  const { value, suffix } = getNumberWithMetricSuffix(
    displayAmount(amount, false)
  )

  return (
    <div className="token-amount flex row center">
      {currencyIcon}
      <span className={amountClassName} style={{ marginLeft: "10px" }}>
        {value}
        <span
          className={suffixClassName}
          style={{ marginLeft: "3px", alignSelf: "flex-end" }}
        >
          {suffix}
        </span>
      </span>
    </div>
  )
}

TokenAmount.defaultProps = {
  currencyIcon: <Icons.KeepOutline />,
  amountClassName: "h1 text-primary",
  suffixClassName: "h3",
}

export default TokenAmount
