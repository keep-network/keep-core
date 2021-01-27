import React, { useMemo, useCallback } from "react"
import CountUp from "react-countup"
import BigNumber from "bignumber.js"
import DoubleIcon from "./DoubleIcon"
import * as Icons from "./Icons"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { Skeleton } from "./skeletons"
import Tooltip from "./Tooltip"
import Banner from "./Banner"
import { toTokenUnit } from "../utils/token.utils"
import { gt } from "../utils/arithmetics.utils"
import { formatPercentage } from "../utils/general.utils"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"

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
  isFetching,
  wrapperClassName = "",
  addLpTokens,
  withdrawLiquidityRewards,
  isAPYFetching,
  pool,
}) => {
  const formattedApy = useMemo(() => {
    const bn = new BigNumber(apy).multipliedBy(100)
    if (bn.isEqualTo(Infinity)) {
      return Infinity
    } else if (bn.isLessThan(0.01) && bn.isGreaterThan(0)) {
      return 0.01
    } else if (bn.isGreaterThan(999)) {
      return 999
    }

    return formatPercentage(bn)
  }, [apy])

  const formattedPercentageOfTotalPool = useMemo(() => {
    const bn = new BigNumber(percentageOfTotalPool)
    return bn.isLessThan(0.01) && bn.isGreaterThan(0)
      ? 0.01
      : formatPercentage(bn)
  }, [percentageOfTotalPool])

  const formattingFn = useCallback((value) => {
    let prefix = ""
    if (value === 0.01) {
      prefix = `<`
    } else if (value >= 999) {
      prefix = `>`
    }
    return `${prefix}${value}%`
  }, [])

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
        <div className={"liquidity__info-tile bg-mint-10"}>
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
              href={"https://thegraph.com/explorer/subgraph/uniswap/uniswap-v2"}
              className="text-white text-link"
            >
              Uniswap subgraph API
            </a>
            &nbsp;to fetch the total pool value and KEEP token in USD.
          </Tooltip>
          {isAPYFetching ? (
            <Skeleton tag="h2" shining color="grey-10" />
          ) : (
            <h2 className={"liquidity__info-tile__title text-mint-100"}>
              {formattedApy === Infinity ? (
                <span>&#8734;</span>
              ) : (
                <CountUp
                  end={formattedApy}
                  // Save previously ended number to start every new animation from it.
                  preserveValue
                  decimals={2}
                  duration={1}
                  formattingFn={formattingFn}
                />
              )}
            </h2>
          )}
          <h6>Estimate of pool apy</h6>
        </div>
        <div className={"liquidity__info-tile bg-mint-10"}>
          {isFetching ? (
            <Skeleton tag="h2" shining color="grey-10" />
          ) : (
            <h2 className={"liquidity__info-tile__title text-mint-100"}>
              <CountUp
                end={formattedPercentageOfTotalPool}
                // Save previously ended number to start every new animation from it.
                preserveValue
                decimals={2}
                duration={1}
                formattingFn={formattingFn}
              />
            </h2>
          )}
          <h6>Your share of POOL</h6>
        </div>
      </div>
      {renderUserInfoBanner()}
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
          withdrawLiquidityRewards(liquidityPairContractName, awaitingPromise)
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
