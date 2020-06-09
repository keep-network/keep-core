import React from "react"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import { Link } from "react-router-dom"
import { colors } from "../constants/colors"
import { add, div, isZero } from "../utils/arithmetics.utils"
import web3Utils from "web3-utils"

const boxWrapperStyle = {
  border: `1px solid ${colors.grey20}`,
  padding: "2rem",
  minWidth: "325px",
  minHeight: "345px",
  margin: "2rem",
}

const TokenOverview = ({
  totalKeepTokenBalance,
  totalGrantedTokenBalance,
  totalGrantedStakedBalance,
  totalOwnedStakedBalance,
}) => {
  const totalKeep = add(totalKeepTokenBalance, totalGrantedTokenBalance)

  const totalStakedInPercentage = add(
    totalGrantedStakedBalance,
    totalOwnedStakedBalance
  )
    .div(isZero(totalKeep) ? web3Utils.toBN("1") : totalKeep)
    .mul(web3Utils.toBN(100))

  const ownedStakedInPercentage = div(
    totalOwnedStakedBalance,
    isZero(totalKeepTokenBalance) ? "1" : totalKeepTokenBalance
  ).mul(web3Utils.toBN(100))

  const grantedStakedInPercentage = div(
    totalGrantedStakedBalance,
    isZero(totalGrantedTokenBalance) ? "1" : totalGrantedTokenBalance
  ).mul(web3Utils.toBN(100))

  return (
    <section className="tile" id="token-overview-balance">
      <div style={{ marginRight: "auto" }}>
        <h2 className="text-grey-70">Total Balance</h2>
        <TokenAmount amount={totalKeep} withMetricSuffix />
        <h3 className="text-grey-30">{`${totalStakedInPercentage}% Staked`}</h3>
      </div>
      <div
        style={boxWrapperStyle}
        className="flex column center space-between mt-1"
      >
        <Icons.MoneyWalletOpen />
        <h4 className="text-grey-70">Granted Tokens</h4>
        <TokenAmount
          currencyIconProps={{ width: 18, heigh: 18 }}
          amountClassName="h4 text-primary"
          suffixClassName="text-small text-primary"
          amount={totalGrantedStakedBalance}
          withMetricSuffix
        />
        <p className="text-small">{`${grantedStakedInPercentage}% Staked`}</p>
        <Link to="/tokens/delegate" className="btn btn-primary mt-2">
          manage
        </Link>
      </div>
      <div
        style={boxWrapperStyle}
        className="flex column center space-between mt-1"
      >
        <Icons.GrantContextIcon />
        <h4 className="text-grey-70">Owned Tokens</h4>
        <TokenAmount
          currencyIconProps={{ width: 18, heigh: 18 }}
          amountClassName="h4 text-primary"
          suffixClassName="text-small text-primary"
          icons
          amount={totalOwnedStakedBalance}
          withMetricSuffix
        />
        <p className="text-small">{`${ownedStakedInPercentage}% Staked`}</p>
        <Link to="/tokens/delegate" className="btn btn-primary mt-2">
          manage
        </Link>
      </div>
    </section>
  )
}

export default React.memo(TokenOverview)
