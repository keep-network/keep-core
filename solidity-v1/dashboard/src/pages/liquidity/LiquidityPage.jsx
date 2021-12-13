import React, { useEffect } from "react"
import { useSelector, useDispatch } from "react-redux"
import {
  useWeb3Address,
  useWeb3Context,
} from "../../components/WithWeb3Context"
import PageWrapper from "../../components/PageWrapper"
import LiquidityRewardCard from "../../components/LiquidityRewardCard"
import {
  LINK,
  LIQUIDITY_REWARD_PAIRS,
  POOL_TYPE,
} from "../../constants/constants"
import * as Icons from "../../components/Icons"
import {
  addMoreLpTokens,
  withdrawAllLiquidityRewards,
} from "../../actions/web3"
import Banner from "../../components/Banner"
import { useHideComponent } from "../../hooks/useHideComponent"
import { gt } from "../../utils/arithmetics.utils"
import MasonryFlexContainer from "../../components/MasonryFlexContainer"
import KeepOnlyPoolCard from "../../components/liquidity/KeepOnlyPoolCard"

const cards = [
  {
    id: "TBTCV2_SADDLE",
    title: LIQUIDITY_REWARD_PAIRS.TBTCV2_SADDLE.label,
    liquidityPairContractName:
      LIQUIDITY_REWARD_PAIRS.TBTCV2_SADDLE.contractName,
    MainIcon: Icons.TBTC_V2,
    SecondaryIcon: Icons.Saddle,
    viewPoolLink: LIQUIDITY_REWARD_PAIRS.TBTCV2_SADDLE.viewPoolLink,
    pool: LIQUIDITY_REWARD_PAIRS.TBTCV2_SADDLE.pool,
    lpTokens: LIQUIDITY_REWARD_PAIRS.TBTCV2_SADDLE.lpTokens,
    wrapperClassName: "tbtc-v2-saddle",
  },
  {
    id: "TBTC_ETH",
    title: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.label,
    liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.contractName,
    MainIcon: Icons.TBTC,
    SecondaryIcon: Icons.EthToken,
    viewPoolLink: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.viewPoolLink,
    pool: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.pool,
    lpTokens: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.lpTokens,
    wrapperClassName: "tbtc-eth",
    incentivesRemoved: true,
    incentivesRemovedBannerProps: {
      link: LINK.proposals.removeIncentivesForTBTCETHpool,
    },
  },
  {
    id: "KEEP_ETH",
    title: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.label,
    liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.contractName,
    MainIcon: Icons.KeepBlackGreen,
    SecondaryIcon: Icons.EthToken,
    viewPoolLink: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink,
    pool: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.pool,
    lpTokens: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.lpTokens,
    wrapperClassName: "keep-eth",
  },
  {
    id: "TBTC_SADDLE",
    title: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.label,
    liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.contractName,
    MainIcon: Icons.TBTC,
    SecondaryIcon: Icons.Saddle,
    viewPoolLink: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.viewPoolLink,
    pool: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.pool,
    lpTokens: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.lpTokens,
    wrapperClassName: "tbtc-saddle",
    incentivesRemoved: true,
    incentivesRemovedBannerProps: {
      link: LINK.proposals.shiftingIncentivesToCoveragePools,
    },
  },
  {
    id: "KEEP_ONLY",
    title: LIQUIDITY_REWARD_PAIRS.KEEP_ONLY.label,
    pool: LIQUIDITY_REWARD_PAIRS.KEEP_ONLY.pool,
    MainIcon: Icons.KeepBlackGreen,
    SecondaryIcon: Icons.Saddle,
    incentivesRemoved: true,
  },
  {
    id: "TBTCV2_MBTC",
    title: LIQUIDITY_REWARD_PAIRS.TBTCV2_MBTC.label,
    MainIcon: Icons.TBTC_V2,
    SecondaryIcon: Icons.MBTC,
    viewPoolLink: LIQUIDITY_REWARD_PAIRS.TBTCV2_MBTC.viewPoolLink,
    pool: LIQUIDITY_REWARD_PAIRS.TBTCV2_MBTC.pool,
    lpTokens: LIQUIDITY_REWARD_PAIRS.TBTCV2_MBTC.lpTokens,
    wrapperClassName: "tbtc-v2-saddle",
  },
  {
    id: "KEEP_TBTC",
    title: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.label,
    liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.contractName,
    MainIcon: Icons.KeepBlackGreen,
    SecondaryIcon: Icons.TBTC,
    viewPoolLink: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.viewPoolLink,
    pool: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.pool,
    lpTokens: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.lpTokens,
    wrapperClassName: "keep-tbtc",
    incentivesRemoved: true,
    incentivesRemovedBannerProps: {
      link: LINK.proposals.removeIncentivesForKEEPTBTCpool,
    },
  },
]

const LiquidityPage = ({ headerTitle }) => {
  const [isBannerVisible, hideBanner] = useHideComponent(false)
  const { isConnected } = useWeb3Context()
  const dispatch = useDispatch()
  const address = useWeb3Address()
  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)
  const { KEEP_ONLY, ...liquidityPools } = useSelector(
    (state) => state.liquidityRewards
  )

  useEffect(() => {
    if (isConnected) {
      dispatch({
        type: "liquidity_rewards/fetch_data_request",
        payload: { address },
      })
    }
  }, [dispatch, address, isConnected])

  useEffect(() => {
    if (!isConnected)
      dispatch({
        type: "liquidity_rewards/fetch_apy_request",
      })
  }, [dispatch, isConnected])

  useEffect(() => {
    if (isBannerVisible && isConnected && gt(keepTokenBalance.value || 0, 0)) {
      hideBanner()
    }
  }, [isConnected, keepTokenBalance.value, hideBanner, isBannerVisible])

  const addLpTokens = (poolId, wrappedTokenBalance, awaitingPromise) => {
    const liquidityRewardPair = LIQUIDITY_REWARD_PAIRS[poolId]
    dispatch(
      addMoreLpTokens(
        liquidityRewardPair,
        wrappedTokenBalance,
        address,
        awaitingPromise
      )
    )
  }

  const withdrawLiquidityRewards = (poolId, amount, awaitingPromise) => {
    const liquidityRewardPair = LIQUIDITY_REWARD_PAIRS[poolId]
    dispatch(
      withdrawAllLiquidityRewards(liquidityRewardPair, amount, awaitingPromise)
    )
  }

  return (
    <PageWrapper title={headerTitle}>
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
                href={
                  "https://balancer.exchange/#/swap/ether/0x85eee30c52b0b379b046fb0f85f4f3dc3009afec"
                }
                className="text-link"
              >
                Balancer
              </a>
              &nbsp;or&nbsp;
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={
                  "https://app.uniswap.org/#/swap?inputCurrency=ETH&outputCurrency=0x85eee30c52b0b379b046fb0f85f4f3dc3009afec"
                }
                className="text-link"
              >
                Uniswap
              </a>
            </Banner.Description>
          </div>
          <Banner.CloseIcon onClick={hideBanner} />
        </Banner>
      )}
      <MasonryFlexContainer maxHeight={"2060px"}>
        {cards.map(({ id, ...data }) => {
          if (data.pool === POOL_TYPE.MSTABLE) {
            return (
              <LiquidityRewardCard
                key={id}
                poolId={id}
                {...data}
                displaySubtitle
                displayGoToPoolButton
              />
            )
          }
          if (data.pool === POOL_TYPE.TOKEN_GEYSER) {
            return (
              <KeepOnlyPoolCard
                key={id}
                poolId={id}
                {...data}
                lpBalance={KEEP_ONLY.lpBalance}
                rewardBalance={KEEP_ONLY.reward}
              />
            )
          }
          return (
            <LiquidityRewardCard
              key={id}
              poolId={id}
              {...data}
              apy={liquidityPools[id].apy}
              percentageOfTotalPool={liquidityPools[id].shareOfPoolInPercent}
              rewardBalance={liquidityPools[id].reward}
              wrappedTokenBalance={liquidityPools[id].wrappedTokenBalance}
              lpBalance={liquidityPools[id].lpBalance}
              lpTokenBalance={liquidityPools[id].lpTokenBalance}
              rewardMultiplier={liquidityPools[id].rewardMultiplier}
              isFetching={liquidityPools[id].isFetching}
              addLpTokens={addLpTokens}
              withdrawLiquidityRewards={withdrawLiquidityRewards}
              isAPYFetching={liquidityPools[id].isAPYFetching}
              displaySubtitle
              displayMetrics={!data.incentivesRemoved}
              displayLPTokenBalance
              displayRewards
            />
          )
        })}
      </MasonryFlexContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
