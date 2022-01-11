import React, { useContext, useMemo } from "react"
import DoubleIcon from "./DoubleIcon"
import OnlyIf from "./OnlyIf"
import Tooltip from "./Tooltip"
import * as Icons from "./Icons"
import { gt } from "../utils/arithmetics.utils"
import MetricsTile from "./MetricsTile"
import { APY, ShareOfPool } from "./liquidity"
import { Skeleton } from "./skeletons"
import CountUp from "react-countup"
import { KEEP } from "../utils/token.utils"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { POOL_TYPE } from "../constants/constants"
import Banner from "./Banner"

const LiquidityRewardCardContext = React.createContext({
  apy: "0",
  percentageOfTotalPool: "0",
  lpBalance: "0",
  wrappedTokenBalance: "0",
  rewardBalance: "0",
  isFetching: false,
  isAPYFetching: false,
})

const useLiquidityRewardCardContext = () => {
  const context = useContext(LiquidityRewardCardContext)
  if (!context) {
    throw new Error(
      "LiquidityRewardCardContext used outside of LiquidityRewardCard component"
    )
  }
  return context
}

const LiquidityRewardCard = ({
  title,
  MainIcon,
  SecondaryIcon,
  wrapperClassName,
  isFetching = false,
  isAPYFetching = false,
  apy = "0",
  percentageOfTotalPool = "0",
  lpBalance = "0",
  wrappedTokenBalance = "0",
  rewardBalance = "0",
  children,
}) => {
  return (
    <LiquidityRewardCardContext.Provider
      value={{
        isFetching,
        isAPYFetching,
        apy,
        percentageOfTotalPool,
        lpBalance,
        wrappedTokenBalance,
        rewardBalance,
      }}
    >
      <Card className={`liquidity__card tile ${wrapperClassName}`}>
        <div className={"liquidity__card-title"}>
          {!!SecondaryIcon ? (
            <DoubleIcon
              MainIcon={MainIcon}
              SecondaryIcon={SecondaryIcon}
              className={`liquidity__double-icon-container`}
            />
          ) : (
            <MainIcon width={24} height={24} className={"mr-1"} />
          )}
          <h2 className={"h2--alt text-grey-70"}>{title}</h2>
        </div>
        {children}
      </Card>
    </LiquidityRewardCardContext.Provider>
  )
}

const Subtitle = ({ pool, viewPoolLink, className = "" }) => {
  const poolName = useMemo(() => {
    switch (pool) {
      case POOL_TYPE.UNISWAP:
        return "Uniswap pool"
      case POOL_TYPE.SADDLE:
        return "Saddle pool"
      case POOL_TYPE.MSTABLE:
        return "mStable"
      default:
        return "Uniswap pool"
    }
  }, [pool])

  return (
    <h4 className={`liquidity__card-subtitle text-grey-40 mb-2 ${className}`}>
      {poolName}
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
        liquidity pool as a liquidity provider. KEEP rewards are proportional to
        your share of the total pool.
      </Tooltip>
    </h4>
  )
}

LiquidityRewardCard.Subtitle = Subtitle

const Metrics = () => {
  const { isFetching, isAPYFetching, apy, lpBalance, percentageOfTotalPool } =
    useLiquidityRewardCardContext()
  return (
    <div
      className={`liquidity__info${gt(lpBalance, 0) ? "" : "--locked"} mb-2`}
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
  )
}

LiquidityRewardCard.Metrics = Metrics

LiquidityRewardCard.InactivePoolBanner = ({ inactivePoolBannerProps = {} }) => {
  const bannerProps = {
    icon: Icons.Bell,
    title: "Incentives removed",
    description:
      "The incentives for this pool has been removed and you can no longer deposit the lp tokens. You can still withdraw rewards that you already earned.",
    link: null,
    linkText: "More info",
    ...inactivePoolBannerProps,
  }
  return <UserInfoBanner incentivesRemoved={true} {...bannerProps} />
}

const ActivePoolBanner = ({
  pool,
  viewPoolLink,
  userInfoBannerProps = null,
}) => {
  const { lpBalance } = useLiquidityRewardCardContext()
  const hasDepositedWrappedTokens = useMemo(() => gt(lpBalance, 0), [lpBalance])

  const bannerProps = {
    icon: !hasDepositedWrappedTokens ? Icons.Rewards : Icons.Wallet,
    title: !hasDepositedWrappedTokens
      ? "Start earning rewards"
      : "No LP Tokens found in wallet",
    description: !hasDepositedWrappedTokens
      ? "Get LP tokens by adding liquidity first to the"
      : "Get more by adding liquidity to the",
    link: viewPoolLink,
    linkText: pool === POOL_TYPE.SADDLE ? "Saddle pool" : "Uniswap pool",
    ...userInfoBannerProps,
  }

  return <UserInfoBanner incentivesRemoved={false} {...bannerProps} />
}

LiquidityRewardCard.ActivePoolBanner = ActivePoolBanner

const UserInfoBanner = ({
  incentivesRemoved = false,
  icon,
  title,
  description,
  link,
  linkText,
}) => {
  return (
    <Banner
      className={`liquidity-info-banner ${
        incentivesRemoved ? "liquidity-info-banner--warning mt-2" : ""
      }`}
    >
      <Banner.Icon
        icon={icon}
        className={`liquidity-info-banner__icon ${
          incentivesRemoved ? "liquidity-info-banner__icon--warning" : ""
        }`}
      />
      <div className={"liquidity-info-banner__content"}>
        <Banner.Title
          className={`liquidity-info-banner__content__title ${
            incentivesRemoved ? "text-black" : "text-white"
          }`}
        >
          {title}
        </Banner.Title>
        <Banner.Description
          className={`liquidity-info-banner__content__description ${
            incentivesRemoved ? "text-black" : "text-white"
          }`}
        >
          {description}
          <OnlyIf condition={link && linkText}>
            &nbsp;
            <a
              target="_blank"
              rel="noopener noreferrer"
              href={link}
              className={`text-link ${
                incentivesRemoved ? "text-grey-60" : "text-white"
              }`}
            >
              {linkText}
            </a>
          </OnlyIf>
        </Banner.Description>
      </div>
    </Banner>
  )
}

const Rewards = () => {
  const { isFetching, rewardBalance } = useLiquidityRewardCardContext()
  return (
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
              end={KEEP.toTokenUnit(rewardBalance).toNumber()}
              separator={","}
              preserveValue
            />
          </h3>
        )}
      </div>
    </div>
  )
}

LiquidityRewardCard.Rewards = Rewards

const ActionButtons = ({
  poolId,
  incentivesRemoved,
  addLpTokens,
  withdrawLiquidityRewards,
}) => {
  const { wrappedTokenBalance, lpBalance, rewardBalance } =
    useLiquidityRewardCardContext()
  return (
    <>
      <OnlyIf condition={!incentivesRemoved}>
        <SubmitButton
          className={`liquidity__add-more-tokens btn btn-primary btn-lg w-100`}
          disabled={!gt(wrappedTokenBalance, 0) || incentivesRemoved}
          onSubmitAction={(awaitingPromise) =>
            addLpTokens(poolId, wrappedTokenBalance, awaitingPromise)
          }
        >
          {gt(lpBalance, 0) ? "add more lp tokens" : "deposit lp tokens"}
        </SubmitButton>
      </OnlyIf>
      <SubmitButton
        className={"liquidity__withdraw btn btn-secondary btn-lg w-100"}
        disabled={!gt(rewardBalance, 0) && !gt(lpBalance, 0)}
        onSubmitAction={(awaitingPromise) =>
          withdrawLiquidityRewards(poolId, lpBalance, awaitingPromise)
        }
      >
        withdraw all
      </SubmitButton>
      {(gt(rewardBalance, 0) || gt(lpBalance, 0)) && (
        <div className={"text-validation text-center"}>
          Withdraw includes rewards and principal
        </div>
      )}
    </>
  )
}

LiquidityRewardCard.ActionButtons = ActionButtons

LiquidityRewardCard.GoToPoolButton = ({ viewPoolLink }) => {
  return (
    <a
      href={viewPoolLink}
      rel="noopener noreferrer"
      target="_blank"
      className={`btn btn-primary btn-lg w-100 mt-2`}
    >
      go to pool
    </a>
  )
}

export default LiquidityRewardCard
