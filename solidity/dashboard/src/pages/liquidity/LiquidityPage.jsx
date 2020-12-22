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
      <CardContainer>
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.KEEP_ETH.label}
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIR.KEEP_ETH.viewPoolLink}
          percentageOfTotalPool={KEEP_ETH.shareOfPoolInPercent}
          rewardBalance={KEEP_ETH.reward}
          wrappedTokenBalance={KEEP_ETH.wrappedTokenBalance}
          lpBalance={KEEP_ETH.lpBalance}
          isFetching={KEEP_ETH.isFetching}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.KEEP_TBTC.label}
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.TBTC}
          viewPoolLink={LIQUIDITY_REWARD_PAIR.KEEP_ETH.viewPoolLink}
          percentageOfTotalPool={KEEP_TBTC.shareOfPoolInPercent}
          rewardBalance={KEEP_TBTC.reward}
          wrappedTokenBalance={KEEP_TBTC.wrappedTokenBalance}
          lpBalance={KEEP_TBTC.lpBalance}
          isFetching={KEEP_TBTC.isFetching}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIR.TBTC_ETH.label}
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIR.KEEP_ETH.viewPoolLink}
          percentageOfTotalPool={TBTC_ETH.shareOfPoolInPercent}
          rewardBalance={TBTC_ETH.reward}
          wrappedTokenBalance={TBTC_ETH.wrappedTokenBalance}
          lpBalance={TBTC_ETH.lpBalance}
          isFetching={TBTC_ETH.isFetching}
        />
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
