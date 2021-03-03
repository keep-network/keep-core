import React, { useMemo } from "react"
import CountUp from "react-countup"
import DoubleIcon from "./DoubleIcon"
import * as Icons from "./Icons"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { Skeleton } from "./skeletons"
import Tooltip from "./Tooltip"
import Banner from "./Banner"
import { toTokenUnit } from "../utils/token.utils"
import { gt } from "../utils/arithmetics.utils"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"
import { APY, LPTokenBalance, ShareOfPool } from "./liquidity"
import MetricsTile from "./MetricsTile"

const LiquidityRewardCard = ({
  title,
  liquidityPairContractName,
  MainIcon,
  SecondaryIcon,
  viewPoolLink,
  apy,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance,
  // Balance of the wrapped token.
  wrappedTokenBalance,
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance,
  lpTokenBalance,
  lpTokens,
  isFetching,
  wrapperClassName = "",
  addLpTokens,
  withdrawLiquidityRewards,
  isAPYFetching,
  pool,
}) => {
  const hasWrappedTokens = useMemo(() => gt(wrappedTokenBalance, 0), [
    wrappedTokenBalance,
  ])

  const hasDepositedWrappedTokens = useMemo(() => gt(lpBalance, 0), [lpBalance])

  const renderUserInfoBanner = () => {
    return (
      !hasWrappedTokens && (
        <Banner className="liquidity__new-user-info">
          <Banner.Icon
            icon={!hasDepositedWrappedTokens ? Icons.Rewards : Icons.Wallet}
            className={"liquidity__rewards-icon"}
          />
          <div className={"liquidity__new-user-info-text"}>
            <Banner.Title className={"liquidity-banner__title text-white"}>
              {!hasDepositedWrappedTokens
                ? "Start earning rewards"
                : "No LP Tokens found in wallet"}
            </Banner.Title>
            <Banner.Description className="liquidity-banner__info text-white">
              {!hasDepositedWrappedTokens
                ? "Get LP tokens by adding liquidity first to the"
                : "Get more by adding liquidity to the"}
              &nbsp;
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={viewPoolLink}
                className="text-white text-link"
              >
                {title === LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.label
                  ? "Saddle pool"
                  : "Uniswap pool"}
              </a>
            </Banner.Description>
          </div>
        </Banner>
      )
    )
  }

  return (
    <Card className={`liquidity__card tile ${wrapperClassName}`}>
      <div className={"liquidity__card-title"}>
        <DoubleIcon
          MainIcon={MainIcon}
          SecondaryIcon={SecondaryIcon}
          className={`liquidity__double-icon-container`}
        />
        <h2 className={"h2--alt text-grey-70"}>{title}</h2>
      </div>
      <h4 className="liquidity__card-subtitle text-grey-40">
        {title === LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.label
          ? "Saddle Pool"
          : "Uniswap Pool"}
        &nbsp;
        <a
          target="_blank"
          rel="noopener noreferrer"
          href={viewPoolLink}
          className="text-small"
        >
          View pool
        </a>
        &nbsp;
        <Tooltip
          simple
          delay={0}
          triggerComponent={Icons.MoreInfo}
          className={"liquidity__card-subtitle__tooltip"}
        >
          LP tokens represent the amount of money you&apos;ve deposited into a
          liquidity pool as a liquidity provider. KEEP rewards are proportional
          to your share of the total pool.
        </Tooltip>
      </h4>
      <div
        className={`liquidity__info${
          gt(lpBalance, 0) ? "" : "--locked"
        } mt-2 mb-2`}
      >
        <MetricsTile className="liquidity__info-tile bg-mint-10">
          <MetricsTile.Tooltip className="liquidity__info-tile__tooltip">
            <APY.TooltipContent />
          </MetricsTile.Tooltip>
          <APY
            apy={apy}
            isFetching={isAPYFetching}
            className="liquidity__info-tile__title text-mint-100"
          />
          <h6>Estimate of pool apy</h6>
        </MetricsTile>
        <MetricsTile className="liquidity__info-tile bg-mint-10">
          <ShareOfPool
            className="liquidity__info-tile__title text-mint-100"
            percentageOfTotalPool={percentageOfTotalPool}
            isFetching={isFetching}
          />
          <h6>Your share of POOL</h6>
        </MetricsTile>
      </div>
      {renderUserInfoBanner()}
      <LPTokenBalance lpTokens={lpTokens} lpTokenBalance={lpTokenBalance} />
      <div className={"liquidity__reward-balance"}>
        <h4 className={"liquidity__reward-balance__title text-grey-70"}>
          Your rewards
        </h4>
        <span className={"liquidity__reward-balance__subtitle text-grey-40"}>
          Rewards allocated on a weekly basis.
        </span>
        <div className={"liquidity__reward-balance_values text-grey-70"}>
          <h3 className={"liquidity__reward-balance_values_label"}>
            <Icons.KeepOutline />
            <span>KEEP</span>
          </h3>
          {isFetching ? (
            <Skeleton tag="h3" shining color="grey-20" className="ml-3" />
          ) : (
            <h3>
              <CountUp
                end={toTokenUnit(rewardBalance).toNumber()}
                separator={","}
                preserveValue
              />
            </h3>
          )}
        </div>
      </div>
      <SubmitButton
        className={`liquidity__add-more-tokens btn btn-primary btn-lg w-100`}
        disabled={!gt(wrappedTokenBalance, 0)}
        onSubmitAction={(awaitingPromise) =>
          addLpTokens(
            wrappedTokenBalance,
            liquidityPairContractName,
            pool,
            awaitingPromise
          )
        }
      >
        {gt(lpBalance, 0) ? "add more lp tokens" : "deposit lp tokens"}
      </SubmitButton>

      <SubmitButton
        className={"liquidity__withdraw btn btn-secondary btn-lg w-100"}
        disabled={!gt(rewardBalance, 0) && !gt(lpBalance, 0)}
        onSubmitAction={(awaitingPromise) =>
          withdrawLiquidityRewards(
            liquidityPairContractName,
            lpBalance,
            pool,
            awaitingPromise
          )
        }
      >
        withdraw all
      </SubmitButton>
      {(gt(rewardBalance, 0) || gt(lpBalance, 0)) && (
        <div className={"text-validation text-center"}>
          Withdraw includes rewards and principal
        </div>
      )}
    </Card>
  )
}

export default LiquidityRewardCard
