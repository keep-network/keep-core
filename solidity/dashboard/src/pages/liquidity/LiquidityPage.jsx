import React from "react"
import PageWrapper from "../../components/PageWrapper"
import CardContainer from "../../components/CardContainer"
import LiquidityRewardCard from "../../components/LiquidityRewardCard"
import { LIQUIDITY_REWARD_PAIR } from "../../constants/constants"
import * as Icons from "../../components/Icons"

const LiquidityPage = ({ title }) => {
  return (
    <PageWrapper title={title}>
      <CardContainer className={"flex wrap"}>
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.KEEP_ETH.label}
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIR.KEEP_ETH.viewPoolLink}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.KEEP_TBTC.label}
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.TBTC}
          viewPoolLink={LIQUIDITY_REWARD_PAIR.KEEP_ETH.viewPoolLink}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.TBTC_ETH.label}
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIR.KEEP_ETH.viewPoolLink}
        />
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
