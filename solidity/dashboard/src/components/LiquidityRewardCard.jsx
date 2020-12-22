import React, { useMemo } from "react"
import BigNumber from "bignumber.js"
import DoubleIcon from "./DoubleIcon"
import * as Icons from "./Icons"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { displayAmount } from "../utils/token.utils"
import { gt } from "../utils/arithmetics.utils"
import { Skeleton } from "./skeletons"
import { addMoreLpTokens, withdrawAllLiquidityRewards } from "../actions/web3"
import { useDispatch } from "react-redux"

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
}) => {
  const dispatch = useDispatch()

  const formattedPercentageOfTotalPool = useMemo(() => {
    const bn = new BigNumber(percentageOfTotalPool)
    return bn.isLessThan(0.01) && bn.isGreaterThan(0)
      ? "<0.01"
      : bn.decimalPlaces(2, BigNumber.ROUND_DOWN)
  }, [percentageOfTotalPool])

  // TODO: get the amount
  const addLpTokens = async (amount, awaitingPromise) => {
    dispatch(
      addMoreLpTokens(amount, liquidityPairContractName, awaitingPromise)
    )
  }

  const withdrawLiquidityRewards = async (awaitingPromise) => {
    dispatch(
      withdrawAllLiquidityRewards(liquidityPairContractName, awaitingPromise)
    )
  }

  return (
    <Card className={"tile"}>
      <div className={"liquidity__card-title-section"}>
        <DoubleIcon
          MainIcon={MainIcon}
          SecondaryIcon={SecondaryIcon}
          className={`liquidity__double-icon-container`}
        />
        <h2 className={"h2--alt text-grey-70"}>{title}</h2>
      </div>
      <div className={"liquidity-card-subtitle-section"}>
        <h4 className="text-grey-40">
          Uniswap Pool&nbsp;
          <a href={viewPoolLink} className="arrow-link text-small">
            View pool
          </a>
        </h4>
      </div>
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
      <div className={"liquidity__add-more-tokens"}>
        <SubmitButton
          className={`btn btn-primary btn-lg w-100`}
          disabled={!gt(wrappedTokenBalance, 0)}
          onSubmitAction={async (awaitingPromise) =>
            await addLpTokens(1, awaitingPromise)
          }
        >
          add more lp tokens
        </SubmitButton>
      </div>
      <div className={"liquidity__withdraw"}>
        <SubmitButton
          className={"btn btn-primary btn-lg w-100 text-black"}
          disabled={!gt(lpBalance, 0)}
          onSubmitAction={async (awaitingPromise) =>
            await withdrawLiquidityRewards(awaitingPromise)
          }
        >
          withdraw all
        </SubmitButton>
      </div>
    </Card>
  )
}

export default LiquidityRewardCard
