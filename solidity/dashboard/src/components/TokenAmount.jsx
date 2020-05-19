import React from "react"
import * as Icons from "./Icons"
import { displayAmount, getNumberWithMetricSuffix } from "../utils/token.utils"

const TokenAmount = ({
  wrapperClassName,
  currencyIcon,
  currencyIconProps,
  amount,
  amountClassName,
  suffixClassName,
  withMetricSuffix,
}) => {
  const { value, suffix } = getNumberWithMetricSuffix(
    displayAmount(amount, false)
  )
  const CurrencyIcon = currencyIcon

  return (
    <div className={`token-amount flex row center ${wrapperClassName || ""}`}>
      <CurrencyIcon {...currencyIconProps} />
      <span className={amountClassName} style={{ marginLeft: "10px" }}>
        {withMetricSuffix ? value : displayAmount(amount)}
        {withMetricSuffix && (
          <span
            className={suffixClassName}
            style={{ marginLeft: "3px", alignSelf: "flex-end" }}
          >
            {suffix}
          </span>
        )}
      </span>
    </div>
  )
}

TokenAmount.defaultProps = {
  currencyIcon: Icons.KeepOutline,
  currencyIconProps: {},
  amountClassName: "h1 text-primary",
  suffixClassName: "h3",
  withMetricSuffix: true,
  wrapperClassName: "",
}

export default TokenAmount
