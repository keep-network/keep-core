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
import Divider from "../../components/Divider"
import { SubmitButton } from "../../components/Button"
import Tooltip from "../../components/Tooltip"
import { Skeleton } from "../../components/skeletons"

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

      <section className="keep-only-pool">
        <section className="tile keep-only-pool__overview">
          <section>
            <h2 className="h2--alt text-grey-70">Total KEEP Locked</h2>
            <h1 className="text-mint-100 mt-2">
              0&nbsp;<span className="h2">KEEP</span>
            </h1>
            <div className="flex row space-between text-grey-40 mt-1">
              <h4>Deposited KEEP tokens</h4>
              <h4 className="self-end">0 KEEP</h4>
            </div>
            <Divider style={{ margin: "0.5rem 0" }} />
            <div className="flex row space-between text-grey-40">
              <h4>Rewarded KEEP tokens</h4>
              <h4 className="self-end">0 KEEP</h4>
            </div>

            <div className="flex row space-between mt-2">
              <SubmitButton className="btn btn-primary btn-lg">
                deposit keep
              </SubmitButton>
              <SubmitButton className="btn btn-secondary btn-lg">
                withdraw all
              </SubmitButton>
            </div>
          </section>
          <section className="keep-only-pool__overview__info-tiles">
            <div className="liquidity__info-tile bg-mint-10 mb-1">
              <Tooltip
                simple
                delay={0}
                triggerComponent={Icons.MoreInfo}
                className={"liquidity__info-tile__tooltip"}
              >
                Pool APY is calculated using the&nbsp;
                <a
                  target="_blank"
                  rel="noopener noreferrer"
                  href="https://thegraph.com/explorer/subgraph/uniswap/uniswap-v2"
                  className="text-white text-link"
                >
                  Uniswap subgraph API
                </a>
                &nbsp;to fetch the the total pool value and KEEP token in USD.
              </Tooltip>
              <h2 className="liquidity__info-tile__title text-mint-100">10%</h2>
              <h6>Estimate of pool apy</h6>
            </div>
            <div className="liquidity__info-tile bg-mint-10">
              {false ? (
                <Skeleton tag="h2" shining color="color-grey-60" />
              ) : (
                <h2 className={"liquidity__info-tile__title text-mint-100"}>
                  20%
                </h2>
              )}
              <h6>% of total pool</h6>
            </div>
          </section>
        </section>
        <section className="keep-only-pool__icon">
          <Icons.KeepOnlyPool />
        </section>
      </section>
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
