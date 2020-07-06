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
  displayAmountFunction,
}) => {
  const { value, suffix } = withMetricSuffix
    ? getNumberWithMetricSuffix(displayAmountFunction(amount, false))
    : { value: "0", suffix: "" }
  const CurrencyIcon = currencyIcon

  return (
    <div className={`token-amount flex row center ${wrapperClassName || ""}`}>
      <CurrencyIcon {...currencyIconProps} />
      <span className={amountClassName} style={{ marginLeft: "10px" }}>
        {withMetricSuffix ? value : displayAmountFunction(amount)}
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
  displayAmountFunction: displayAmount,
}

export default TokenAmount
