import React from "react"
import BigNumber from "bignumber.js"
import { Keep } from "../../contracts"

const BaseExchangeRate = ({
  covToken,
  collateralToken,
  amount,
  htmlTag = "div",
  className = "",
}) => {
  const Tag = htmlTag
  return (
    <Tag className={className}>
      {`1 ${covToken.symbol}`} =&nbsp;
      {collateralToken.displayAmountWithSymbol(amount, 3, (amount) =>
        new BigNumber(amount).toFormat(3, BigNumber.ROUND_DOWN)
      )}
    </Tag>
  )
}

export const CoveragePoolV1ExchangeRate = ({
  covToken,
  collateralToken,
  covTotalSupply,
  totalValueLocked,
  ...restProps
}) => (
  <BaseExchangeRate
    {...restProps}
    covToken={covToken}
    collateralToken={collateralToken}
    amount={Keep.coveragePoolV1.estimatedBalanceFor(
      collateralToken.fromTokenUnit(1).toString(),
      covTotalSupply,
      totalValueLocked
    )}
  />
)
