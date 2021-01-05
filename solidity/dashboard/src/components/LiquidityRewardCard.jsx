import React, { useMemo } from "react"
import BigNumber from "bignumber.js"
import DoubleIcon from "./DoubleIcon"
import * as Icons from "./Icons"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { displayAmount } from "../utils/token.utils"
import { gt } from "../utils/arithmetics.utils"
import { Skeleton } from "./skeletons"
import Tooltip from "./Tooltip"

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
}) => {
  const formattedApy = useMemo(() => {
    const bn = new BigNumber(apy)
    if (bn.isEqualTo(Infinity)) return <span>&#8734;</span>

    return bn.isLessThan(0.01) && bn.isGreaterThan(0)
      ? "<0.01%"
      : bn.multipliedBy(100).decimalPlaces(2, BigNumber.ROUND_DOWN).toString() +
          "%"
  }, [apy])

  const formattedPercentageOfTotalPool = useMemo(() => {
    const bn = new BigNumber(percentageOfTotalPool)
    return bn.isLessThan(0.01) && bn.isGreaterThan(0)
      ? "<0.01"
      : bn.decimalPlaces(2, BigNumber.ROUND_DOWN)
  }, [percentageOfTotalPool])

  return (
    <Card className={`liquidity__card tile ${wrapperClassName}`}>
      <Icons.SantaHat className="liquidity-card__santa-hat" />
      <div className={"liquidity__card-title"}>
        <DoubleIcon
          MainIcon={MainIcon}
          SecondaryIcon={SecondaryIcon}
          className={`liquidity__double-icon-container`}
        />
        <h2 className={"h2--alt text-grey-70"}>{title}</h2>
      </div>
      <h4 className="liquidity__card-subtitle text-grey-40">
        Uniswap Pool&nbsp;
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
      <div className={"liquidity__info text-grey-60 mt-2 mb-2"}>
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
            &nbsp;to fetch the the total pool value and KEEP token in USD.
          </Tooltip>
          <h2 className={"liquidity__info-tile__title text-mint-100"}>
            {formattedApy}
          </h2>
          <h6>Estimate of pool apy</h6>
        </div>
        <div className={"liquidity__info-tile bg-mint-10"}>
          {isFetching ? (
            <Skeleton tag="h2" shining color="mint-20" />
          ) : (
            <h2
              className={"liquidity__info-tile__title text-mint-100"}
            >{`${formattedPercentageOfTotalPool}%`}</h2>
          )}
          <h6>% of total pool</h6>
        </div>
      </div>
      {!gt(wrappedTokenBalance, 0) && (
        <div className={"liquidity__new-user-info"}>
          <Icons.Rewards className={"liquidity__rewards-icon"} />
          <div className={"liquidity__new-user-info-text"}>
            <h4>Start earning rewards</h4>
            <span>
              Get LP tokens by adding liquidity first to the&nbsp;
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={viewPoolLink}
                className="text-white text-link"
              >
                Uniswap pool
              </a>
            </span>
          </div>
        </div>
      )}
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
            <h3>{displayAmount(rewardBalance)}</h3>
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
            awaitingPromise
          )
        }
      >
        add more lp tokens
      </SubmitButton>

      <SubmitButton
        className={"liquidity__withdraw btn btn-secondary btn-lg w-100"}
        disabled={!gt(rewardBalance, 0)}
        onSubmitAction={(awaitingPromise) =>
          withdrawLiquidityRewards(liquidityPairContractName, awaitingPromise)
        }
      >
        withdraw all
      </SubmitButton>
      {gt(rewardBalance, 0) && (
        <div className={"text-validation text-center"}>
          Withdraw includes rewards and principal
        </div>
      )}
    </Card>
  )
}

export default LiquidityRewardCard
