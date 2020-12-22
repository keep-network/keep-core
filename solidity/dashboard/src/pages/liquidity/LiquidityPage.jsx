import React from "react"
import PageWrapper from "../../components/PageWrapper"
import CardContainer from "../../components/CardContainer"
import LiquidityRewardCard from "../../components/LiquidityRewardCard"
import { LIQUIDITY_REWARD_PAIRS } from "../../constants/constants"
import * as Icons from "../../components/Icons"

const LiquidityPage = ({ title }) => {
  return (
    <PageWrapper title={title}>
      <CardContainer className={"flex wrap"}>
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.KEEP_ETH.contractName
          }
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.contractName
          }
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.TBTC}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.TBTC_ETH.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.TBTC_ETH.contractName
          }
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink}
        />
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
