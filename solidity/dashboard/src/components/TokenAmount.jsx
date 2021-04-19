import React from "react"
import Tooltip from "./Tooltip"
import { KEEP } from "../utils/token.utils"

const TokenAmount = ({
  amount,
  token = KEEP,
  wrapperClassName = "",
  amountClassName = "",
  icon = null,
  iconProps = {},
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
      <Tooltip
        simple
        triggerComponent={() => (
          <span className={`h2 ${amountClassName}`}>{formattedAmount}</span>
        )}
        delay={0}
        className="token-amount-tooltip"
      >
        {`${token.toFormat(
          token.toTokenUnit(amount, _smallestPrecisionDecimals),
          _smallestPrecisionDecimals
        )} ${_smallestPrecisionUnit}`}
      </Tooltip>
      <span className="h3">&nbsp;{token.symbol}</span>
    </div>
  )
}

export default TokenAmount
