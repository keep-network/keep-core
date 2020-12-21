import React, { useEffect } from "react"
import { useSelector, useDispatch } from "react-redux"
import { useWeb3Address } from "../../components/WithWeb3Context"
import PageWrapper from "../../components/PageWrapper"
import CardContainer from "../../components/CardContainer"
import LiquidityRewardCard from "../../components/LiquidityRewardCard"
import { LIQUIDITY_REWARD_PAIR } from "../../constants/constants"
import * as Icons from "../../components/Icons"

const LiquidityPage = ({ title }) => {
  const { KEEP_ETH, TBTC_ETH, KEEP_TBTC } = useSelector(
    (state) => state.liquidityRewards
  )
  const dispatch = useDispatch()
  const address = useWeb3Address()

  useEffect(() => {
    dispatch({
      type: "liquidity_rewards/fetch_data_request",
      payload: { address },
    })
  }, [dispatch, address])

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
