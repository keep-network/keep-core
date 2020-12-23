import React, { useMemo } from "react"
import BigNumber from "bignumber.js"
import DoubleIcon from "./DoubleIcon"
import * as Icons from "./Icons"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { displayAmount } from "../utils/token.utils"
import { gt } from "../utils/arithmetics.utils"
import { Skeleton } from "./skeletons"

const LiquidityRewardCard = ({
  title,
  liquidityPairContractName,
  MainIcon,
  SecondaryIcon,
  viewPoolLink,
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
        <a href={viewPoolLink} className="text-small">
          View pool
        </a>
      </h4>
      <div className={"liquidity__info text-grey-60"}>
        <div className={"liquidity__info-tile bg-mint-10"}>
          <h2 className={"liquidity__info-tile__title text-mint-100"}>200%</h2>
          <h6>Anual % yield</h6>
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
      <div className={"liquidity__token-balance"}>
        <span className={"liquidity__token-balance_title text-grey-70"}>
          Reward
        </span>
        <div className={"liquidity__token-balance_values text-grey-70"}>
          <h3 className={"liquidity__token-balance_values_label"}>
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
    </Card>
  )
}

export default LiquidityRewardCard
