import React, { useEffect } from "react"
import { useSelector, useDispatch } from "react-redux"
import {
  useWeb3Address,
  useWeb3Context,
} from "../../components/WithWeb3Context"
import PageWrapper from "../../components/PageWrapper"
import CardContainer from "../../components/CardContainer"
import LiquidityRewardCard from "../../components/LiquidityRewardCard"
import { LIQUIDITY_REWARD_PAIRS } from "../../constants/constants"
import * as Icons from "../../components/Icons"
import {
  addMoreLpTokens,
  withdrawAllLiquidityRewards,
} from "../../actions/web3"
import Banner from "../../components/Banner"
import { useHideComponent } from "../../hooks/useHideComponent"
import { gt } from "../../utils/arithmetics.utils"

const LiquidityPage = ({ headerTitle }) => {
  const [isBannerVisible, hideBanner] = useHideComponent(false)
  const { isConnected } = useWeb3Context()
  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)

  const { TBTC_SADDLE, KEEP_ETH, TBTC_ETH, KEEP_TBTC } = useSelector(
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

  useEffect(() => {
    dispatch({
      type: "liquidity_rewards/fetch_apy_request",
    })
  }, [dispatch])

  useEffect(() => {
    if (isBannerVisible && isConnected && gt(keepTokenBalance.value || 0, 0)) {
      hideBanner()
    }
  }, [isConnected, keepTokenBalance.value, hideBanner, isBannerVisible])

  const addLpTokens = (
    wrappedTokenBalance,
    liquidityPairContractName,
    pool,
    awaitingPromise
  ) => {
    dispatch(
      addMoreLpTokens(
        wrappedTokenBalance,
        address,
        liquidityPairContractName,
        pool,
        awaitingPromise
      )
    )
  }

  const withdrawLiquidityRewards = (
    liquidityPairContractName,
    awaitingPromise
  ) => {
    dispatch(
      withdrawAllLiquidityRewards(liquidityPairContractName, awaitingPromise)
    )
  }

  return (
    <PageWrapper title={headerTitle} newPage={true}>
      {isBannerVisible && (
        <Banner className="liquidity-banner">
          <Banner.Icon
            icon={Icons.KeepGreenOutline}
            className={"liquidity-banner__keep-logo"}
          />
          <div className={"liquidity-banner__content"}>
            <Banner.Title className={"liquidity-banner__title"}>
              Don’t yet have KEEP tokens?
            </Banner.Title>
            <Banner.Description className="text-secondary liquidity-banner__info">
              What are you waiting for? KEEP can be bought on the open market
              on&nbsp;
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={"https://balancer.exchange/#/swap"}
                className="text-link"
              >
                Balancer
              </a>
              &nbsp;or&nbsp;
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={"https://app.uniswap.org/#/swap"}
                className="text-link"
              >
                Uniswap
              </a>
            </Banner.Description>
          </div>
          <Banner.CloseIcon onClick={hideBanner} />
        </Banner>
      )}

      <CardContainer>
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.contractName
          }
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.Saddle}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.viewPoolLink}
          apy={TBTC_SADDLE.apy}
          percentageOfTotalPool={TBTC_SADDLE.shareOfPoolInPercent}
          rewardBalance={TBTC_SADDLE.reward}
          wrappedTokenBalance={TBTC_SADDLE.wrappedTokenBalance}
          lpBalance={TBTC_SADDLE.lpBalance}
          lpTokenBalance={TBTC_SADDLE.lpTokenBalance}
          isFetching={TBTC_SADDLE.isFetching}
          wrapperClassName="tbtc-saddle"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
          isAPYFetching={TBTC_SADDLE.isAPYFetching}
          pool={LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.pool}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.KEEP_ETH.contractName
          }
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink}
          apy={KEEP_ETH.apy}
          percentageOfTotalPool={KEEP_ETH.shareOfPoolInPercent}
          rewardBalance={KEEP_ETH.reward}
          wrappedTokenBalance={KEEP_ETH.wrappedTokenBalance}
          lpBalance={KEEP_ETH.lpBalance}
          lpTokenBalance={KEEP_ETH.lpTokenBalance}
          isFetching={KEEP_ETH.isFetching}
          wrapperClassName="keep-eth"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
          isAPYFetching={KEEP_ETH.isAPYFetching}
          pool={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.pool}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.contractName
          }
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.TBTC}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.viewPoolLink}
          apy={KEEP_TBTC.apy}
          percentageOfTotalPool={KEEP_TBTC.shareOfPoolInPercent}
          rewardBalance={KEEP_TBTC.reward}
          wrappedTokenBalance={KEEP_TBTC.wrappedTokenBalance}
          lpBalance={KEEP_TBTC.lpBalance}
          lpTokenBalance={KEEP_TBTC.lpTokenBalance}
          isFetching={KEEP_TBTC.isFetching}
          wrapperClassName="keep-tbtc"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
          isAPYFetching={KEEP_TBTC.isAPYFetching}
          pool={LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.pool}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.TBTC_ETH.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.TBTC_ETH.contractName
          }
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.TBTC_ETH.viewPoolLink}
          apy={TBTC_ETH.apy}
          percentageOfTotalPool={TBTC_ETH.shareOfPoolInPercent}
          rewardBalance={TBTC_ETH.reward}
          wrappedTokenBalance={TBTC_ETH.wrappedTokenBalance}
          lpBalance={TBTC_ETH.lpBalance}
          lpTokenBalance={TBTC_ETH.lpTokenBalance}
          isFetching={TBTC_ETH.isFetching}
          wrapperClassName="tbtc-eth"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
          isAPYFetching={TBTC_ETH.isAPYFetching}
          pool={LIQUIDITY_REWARD_PAIRS.TBTC_ETH.pool}
        />
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
