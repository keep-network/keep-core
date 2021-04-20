import React from "react"
import Tooltip from "./Tooltip"
import { KEEP } from "../utils/token.utils"

const TokenAmount = ({
  amount,
  token = KEEP,
  wrapperClassName = "",
  amountClassName = "h2 text-mint-100",
  amountStyles = {},
  symbolClassName = "h3 text-mint-100",
  symbolStyles = {},
  icon = null,
  iconProps = { className: "keep-outline keep-outline--mint-80" },
  withIcon = false,
  withMetricSuffix = false,
  smallestPrecisionUnit = null,
  smallestPrecisionDecimals = null,
}) => {
  const CurrencyIcon = withIcon ? icon || token.icon : () => <></>

  const _smallestPrecisionUnit =
    smallestPrecisionUnit || token.smallestPrecisionUnit

  const _smallestPrecisionDecimals =
    smallestPrecisionDecimals || token.smallestPrecisionDecimals

  const formattedAmount = withMetricSuffix
    ? token.displayAmountWithMetricSuffix(amount)
    : token.displayAmount(amount)

  return (
    <div className={`flex row center ${wrapperClassName}`}>
      <CurrencyIcon width={32} height={32} {...iconProps} />
      &nbsp;
      <Tooltip
        simple
        triggerComponent={() => (
          <span className={amountClassName} style={amountStyles}>
            {formattedAmount}
          </span>
        )}
        delay={0}
        className="token-amount-tooltip"
      >
        {`${token.toFormat(
          token.toTokenUnit(amount, _smallestPrecisionDecimals),
          _smallestPrecisionDecimals
        )} ${_smallestPrecisionUnit}`}
      </Tooltip>
      <span className={symbolClassName} style={symbolStyles}>
        &nbsp;{token.symbol}
      </span>
    </div>
  )
}

export default TokenAmount
