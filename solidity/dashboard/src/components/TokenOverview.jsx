import React, { useMemo } from "react"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import { Link } from "react-router-dom"
import { colors } from "../constants/colors"
import { add, percentageOf } from "../utils/arithmetics.utils"

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
  const totalKeep = useMemo(() => {
    return add(totalKeepTokenBalance, totalGrantedTokenBalance)
  }, [totalKeepTokenBalance, totalGrantedTokenBalance])

  const totalStakedInPercentage = useMemo(() => {
    return percentageOf(
      add(totalGrantedStakedBalance, totalOwnedStakedBalance),
      totalKeep
    )
  }, [totalGrantedStakedBalance, totalOwnedStakedBalance, totalKeep])

  const ownedStakedInPercentage = useMemo(() => {
    return percentageOf(totalOwnedStakedBalance, totalKeepTokenBalance)
  }, [totalOwnedStakedBalance, totalKeepTokenBalance])

  const grantedStakedInPercentage = useMemo(() => {
    return percentageOf(totalGrantedStakedBalance, totalGrantedTokenBalance)
  }, [totalGrantedStakedBalance, totalGrantedTokenBalance])

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
          amount={totalGrantedTokenBalance}
          withMetricSuffix
        />
        <p className="text-small">{`${grantedStakedInPercentage}% Staked`}</p>
        <Link
          to={{ pathname: "/tokens/delegate", hash: "#granted" }}
          className="btn btn-primary mt-2"
        >
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
          amount={totalKeepTokenBalance}
          withMetricSuffix
        />
        <p className="text-small">{`${ownedStakedInPercentage}% Staked`}</p>
        <Link
          to={{ pathname: "/tokens/delegate", hash: "#owned" }}
          className="btn btn-primary mt-2"
        >
          manage
        </Link>
      </div>
    </section>
  )
}

export default React.memo(TokenOverview)
