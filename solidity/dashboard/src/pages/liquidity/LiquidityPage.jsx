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
          title={LIQUIDITY_REWARD_PAIR.KEEP_ETH}
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.EthToken}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.KEEP_TBTC}
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.TBTC}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.TBTC_ETH}
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.EthToken}
        />
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
