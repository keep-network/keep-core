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
import ClosableContainer from "../../components/ClosableContainer";
import { gt } from "../../utils/arithmetics.utils";

const LiquidityPage = ({ headerTitle }) => {
  const { yourAddress, provider } = useWeb3Context()
  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)

  const isActive = !!(yourAddress && provider)

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

  const addLpTokens = (
    wrappedTokenBalance,
    liquidityPairContractName,
    awaitingPromise
  ) => {
    dispatch(
      addMoreLpTokens(
        wrappedTokenBalance,
        address,
        liquidityPairContractName,
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
      <ClosableContainer
        className={"tile liquidity-banner"}
        hide={isActive || gt(keepTokenBalance.value, 0)}
      >
        <div className={"liquidity-banner__keep-logo"}>
          <Icons.KeepGreenOutline />
        </div>
        <div className={"liquidity-banner__content"}>
          <h4 className={"liquidity-banner__title"}>
            Don&apos;t yet have KEEP tokens?
          </h4>
          <span className={"liquidity-banner__info text-small"}>
            What are you waiting for? KEEP can be bought on the open market on&nbsp;
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
          </span>
        </div>
      </ClosableContainer>

      <CardContainer>
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.KEEP_ETH.contractName
          }
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink}
          apy={3.65}
          percentageOfTotalPool={KEEP_ETH.shareOfPoolInPercent}
          rewardBalance={KEEP_ETH.reward}
          wrappedTokenBalance={KEEP_ETH.wrappedTokenBalance}
          lpBalance={KEEP_ETH.lpBalance}
          isFetching={KEEP_ETH.isFetching}
          wrapperClassName="keep-eth"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.contractName
          }
          MainIcon={Icons.KeepBlackGreen}
          SecondaryIcon={Icons.TBTC}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.viewPoolLink}
          apy={6.68}
          percentageOfTotalPool={KEEP_TBTC.shareOfPoolInPercent}
          rewardBalance={KEEP_TBTC.reward}
          wrappedTokenBalance={KEEP_TBTC.wrappedTokenBalance}
          lpBalance={KEEP_TBTC.lpBalance}
          isFetching={KEEP_TBTC.isFetching}
          wrapperClassName="keep-tbtc"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
        />
        <LiquidityRewardCard
          title={LIQUIDITY_REWARD_PAIRS.TBTC_ETH.label}
          liquidityPairContractName={
            LIQUIDITY_REWARD_PAIRS.TBTC_ETH.contractName
          }
          MainIcon={Icons.TBTC}
          SecondaryIcon={Icons.EthToken}
          viewPoolLink={LIQUIDITY_REWARD_PAIRS.TBTC_ETH.viewPoolLink}
          apy={0.67}
          percentageOfTotalPool={TBTC_ETH.shareOfPoolInPercent}
          rewardBalance={TBTC_ETH.reward}
          wrappedTokenBalance={TBTC_ETH.wrappedTokenBalance}
          lpBalance={TBTC_ETH.lpBalance}
          isFetching={TBTC_ETH.isFetching}
          wrapperClassName="tbtc-eth"
          addLpTokens={addLpTokens}
          withdrawLiquidityRewards={withdrawLiquidityRewards}
        />
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
