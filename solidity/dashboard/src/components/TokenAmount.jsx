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
  displayWithMetricSuffix,
  currencySymbol,
  displayAmountFunction,
  withTooltip,
  tooltipText,
}) => {
  const { value, suffix } = displayWithMetricSuffix
    ? getNumberWithMetricSuffix(displayAmountFunction(amount, false))
    : { value: "0", suffix: "" }
  const CurrencyIcon = currencyIcon

  return (
    <div
      className={`token-amount tooltip flex row center ${
        wrapperClassName || ""
      }`}
    >
      {withTooltip && <span className="tooltip-text top">{tooltipText}</span>}
      <CurrencyIcon {...currencyIconProps} />
      <span className={amountClassName} style={{ marginLeft: "10px" }}>
        {displayWithMetricSuffix ? value : displayAmountFunction(amount)}
        {displayWithMetricSuffix && (
          <span
            className={suffixClassName}
            style={{ marginLeft: "3px", alignSelf: "flex-end" }}
          >
            {suffix}
          </span>
        )}
        {currencySymbol && <span>&nbsp;{currencySymbol}</span>}
      </span>
    </div>
  )
}

TokenAmount.defaultProps = {
  currencyIcon: Icons.KeepOutline,
  currencyIconProps: {},
  amountClassName: "h1 text-primary",
  suffixClassName: "h3",
  displayWithMetricSuffix: true,
  wrapperClassName: "",
  currencySymbol: null,
  displayAmountFunction: displayAmount,
}

export default TokenAmount
