import React from "react"
import DoubleIcon from "./DoubleIcon"
import * as Icons from "./Icons"
import { SubmitButton } from "./Button"
import Card from "./Card"
import { addMoreLpTokens, withdrawAllLiquidityRewards } from "../actions/web3"
import { useDispatch } from "react-redux"

const LiquidityRewardCard = ({
  title,
  MainIcon,
  SecondaryIcon,
  viewPoolLink,
}) => {
  const dispatch = useDispatch()

  const addLpTokens = async (awaitingPromise) => {
    // TODO: get the amount
    const amount = 1
    dispatch(addMoreLpTokens(amount, awaitingPromise))
  }

  const withdrawLiquidityRewards = async (awaitingPromise) => {
    dispatch(withdrawAllLiquidityRewards(awaitingPromise))
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
        <span className="text-grey-40">Uniswap Pool&nbsp;</span>
        <a
          href={viewPoolLink}
          className="arrow-link text-small"
          style={{ marginLeft: "auto", marginRight: "2rem" }}
        >
          View pool
        </a>
      </div>
      <div className={"liquidity__info text-grey-60"}>
        <div className={"liquidity__info-tile bg-mint-10"}>
          <h2 className={"liquidity__info-tile__title text-mint-100"}>200%</h2>
          <h6>Anual % yield</h6>
        </div>
        <div className={"liquidity__info-tile bg-mint-10"}>
          <h2 className={"liquidity__info-tile__title text-mint-100"}>10%</h2>
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
          <h3>1,000,000</h3>
        </div>
      </div>
      <div className={"liquidity__add-more-tokens"}>
        <SubmitButton
          className={`btn btn-primary btn-lg w-100`}
          onSubmitAction={async (awaitingPromise) =>
            await addLpTokens(awaitingPromise)
          }
        >
          add more lp tokens
        </SubmitButton>
      </div>
      <div className={"liquidity__withdraw"}>
        <SubmitButton
          className={"btn btn-primary btn-lg w-100 text-black"}
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
