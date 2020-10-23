import React, { useMemo } from "react"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import { Link } from "react-router-dom"
import { colors } from "../constants/colors"
import { add, percentageOf } from "../utils/arithmetics.utils"
import Divider from "./Divider"
import ProgressBar from "./ProgressBar"
import Chip from "./Chip"
import { SpeechBubbleTooltip } from "./SpeechBubbleTooltip"

const TokenOverview = ({
  totalKeepTokenBalance,
  totalGrantedTokenBalance,
  totalGrantedStakedBalance,
  totalOwnedStakedBalance,
}) => {
  // const totalKeep = useMemo(() => {
  //   return add(totalKeepTokenBalance, totalGrantedTokenBalance)
  // }, [totalKeepTokenBalance, totalGrantedTokenBalance])

  // const totalStakedInPercentage = useMemo(() => {
  //   return percentageOf(
  //     add(totalGrantedStakedBalance, totalOwnedStakedBalance),
  //     totalKeep
  //   )
  // }, [totalGrantedStakedBalance, totalOwnedStakedBalance, totalKeep])

  // const ownedStakedInPercentage = useMemo(() => {
  //   return percentageOf(totalOwnedStakedBalance, totalKeepTokenBalance)
  // }, [totalOwnedStakedBalance, totalKeepTokenBalance])

  // const grantedStakedInPercentage = useMemo(() => {
  //   return percentageOf(totalGrantedStakedBalance, totalGrantedTokenBalance)
  // }, [totalGrantedStakedBalance, totalGrantedTokenBalance])

  return (
    <div className="balances-layout">
      <TotalKeepBalance />
      <TokenBalance type="granted" icon={Icons.Grant} />
      <TokenBalance type="wallet" />
    </div>
  )
}

const TotalKeepBalance = () => {
  return (
    <section className="balance__overview">
      <h3>Total KEEP Balance</h3>
      <TokenAmount wrapperClassName="mb-3" />
      <div className="balance__overview__granted">
        <h4>Granted Tokens</h4>
        <h4 className="ml-a">No data to display</h4>
      </div>
      <Divider className="balance__overview__divider" />
      <div className="balance__overview__wallet">
        <h4>Wallet Tokens</h4>
        <h4 className="ml-a">No data to display</h4>
      </div>
    </section>
  )
}

const TokenBalance = ({
  tooltipProps,
  balance,
  total,
  type = "wallet",
  icon: IconComponent = Icons.Wallet,
}) => {
  return (
    <section className={`balance__${type}`}>
      <header className="flex row center mb-2">
        <Chip
          icon={<IconComponent />}
          color="disabled"
          className={`balance__${type}__chip`}
        />
        <div>
          <h3 className="text-grey-70 flex row center">
            {type}&nbsp;
            {tooltipProps && <SpeechBubbleTooltip {...tooltipProps} />}
          </h3>
          <h4 className="text-grey-40">{balance || "No data"}</h4>
        </div>
      </header>
      <ProgressBar items={[]} total={total} styles={{ margin: 0 }} />
      <span className="text-small text-grey-40">-% Staked</span>
      <Link
        to={{ pathname: "/delegate", hash: type }}
        className="btn btn-secondary btn-lg mt-2"
        style={{ width: "100%" }}
      >
        stake
      </Link>
    </section>
  )
}

export default React.memo(TokenOverview)
