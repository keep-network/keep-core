import React, { useEffect, useMemo, useState } from "react"
import BigNumber from "bignumber.js"
import { Keep } from "../../contracts"
import moment from "moment"
import * as Icons from "../Icons"
import Tooltip from "../Tooltip"

const BaseExchangeRate = ({
  covToken,
  collateralToken,
  amount,
  htmlTag = "div",
  className = "",
  style = {},
}) => {
  const Tag = htmlTag
  return (
    <Tag className={className} style={style}>
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
}) => {
  const [
    exchangeRateCalculatedUnixTimestamp,
    setExchangeRateCalculatedUnixTimestamp,
  ] = useState(0)

  const calculateAmount = useMemo(() => {
    const amount = Keep.coveragePoolV1.estimatedBalanceFor(
      collateralToken.fromTokenUnit(1).toString(),
      covTotalSupply,
      totalValueLocked
    )
    return amount
  }, [collateralToken, covTotalSupply, totalValueLocked])

  useEffect(() => {
    setExchangeRateCalculatedUnixTimestamp(moment().unix())
  }, [calculateAmount])

  const exchangeRateCalculatedMoment = useMemo(() => {
    return moment.unix(exchangeRateCalculatedUnixTimestamp)
  }, [exchangeRateCalculatedUnixTimestamp])

  return (
    <>
      <BaseExchangeRate
        {...restProps}
        covToken={covToken}
        collateralToken={collateralToken}
        amount={calculateAmount}
        style={{ marginRight: "0.3rem" }}
      />
      <Tooltip simple delay={0} triggerComponent={Icons.MoreInfo}>
        {exchangeRateCalculatedUnixTimestamp === 0
          ? "Loading..."
          : `Rate calculated on ${exchangeRateCalculatedMoment.format(
              "MM/DD/YYYY"
            )} at ${exchangeRateCalculatedMoment.format("HH:mm:ss")}`}
      </Tooltip>
    </>
  )
}
