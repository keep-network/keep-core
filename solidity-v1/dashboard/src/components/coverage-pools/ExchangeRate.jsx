import React from "react"
import { Keep } from "../../contracts"
import TokenAmount from "../TokenAmount"

const BaseExchangeRate = ({
  covToken,
  amount,
  htmlTag = "div",
  className = "",
}) => {
  const Tag = htmlTag
  return (
    <Tag className={className}>
      {`1 ${covToken.symbol}`} =&nbsp;
      <TokenAmount
        amount={amount}
        decimalsToDisplay={3}
        wrapperClassName="flex-inline"
        amountClassName=""
        symbolClassName=""
      />
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
    amount={Keep.coveragePoolV1.estimatedBalanceFor(
      collateralToken.fromTokenUnit(1).toString(),
      covTotalSupply,
      totalValueLocked
    )}
  />
)
