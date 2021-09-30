import React from "react"
import Tooltip from "./Tooltip"
import { covKEEP, KEEP } from "../utils/token.utils"
import OnlyIf from "./OnlyIf"

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
  decimalsToDisplay = null,
  displayDecimalsInTooltip = false,
  withSymbol = true,
  symbol = "",
}) => {
  if (token === KEEP || covKEEP) {
    displayDecimalsInTooltip = true
  }
  const CurrencyIcon = withIcon ? icon || token.icon : () => <></>

  const _smallestPrecisionUnit =
    smallestPrecisionUnit || token.smallestPrecisionUnit

  const _smallestPrecisionDecimals =
    smallestPrecisionDecimals || token.smallestPrecisionDecimals

  const _symbol = symbol || token.symbol
  const _decimalsToDisplay = decimalsToDisplay | token.decimalsToDisplay

  const _decimalsToDisplayInTooltip = displayDecimalsInTooltip
    ? token.smallestPrecisionDecimals
    : 0

  const formattedAmount = withMetricSuffix
    ? token.displayAmountWithMetricSuffix(amount, _decimalsToDisplay)
    : token.displayAmount(amount, _decimalsToDisplay)

  return (
    <div className={`token-amount ${wrapperClassName}`}>
      <CurrencyIcon
        width={32}
        height={32}
        {...iconProps}
        className={`token-amount__icon ${iconProps.className}`}
      />
      <OnlyIf condition={icon}>&nbsp;</OnlyIf>
      <Tooltip
        simple
        triggerComponent={() => (
          <span className={amountClassName} style={amountStyles}>
            {formattedAmount}
          </span>
        )}
        delay={0}
        className="token-amount__tooltip"
        contentWrapperStyles={
          _decimalsToDisplayInTooltip > 0
            ? {
                wordBreak: "break-word",
              }
            : null
        }
      >
        {`${token.toFormat(
          token.toTokenUnit(amount, _smallestPrecisionDecimals),
          _decimalsToDisplayInTooltip
        )} ${_smallestPrecisionUnit}`}
      </Tooltip>
      {withSymbol && (
        <span className={symbolClassName} style={symbolStyles}>
          &nbsp;{_symbol}
        </span>
      )}
    </div>
  )
}

export default TokenAmount
