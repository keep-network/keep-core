import React, { useMemo, useCallback } from "react"
import CountUp from "react-countup"
import Divider from "./Divider"
import { SubmitButton } from "./Button"
import * as Icons from "./Icons"
import { APY, ShareOfPool } from "./liquidity"
import { gt } from "../utils/arithmetics.utils"
import {
  //   displayAmountWithMetricSuffix,
  //   getNumberWithMetricSuffix,
  toTokenUnit,
  displayAmount,
  fromTokenUnit,
} from "../utils/token.utils"

const KeepOnlyPool = ({
  apy,
  lpBalance,
  rewardBalance,
  wrappedTokenBalance,
  isFetching,
  isAPYFetching,
  percentageOfTotalPool,
  addLpTokens,
  withdrawLiquidityRewards,
  liquidityContractName,
  pool,
}) => {
  const lockedKEEP = useMemo(() => {
    return add(lpBalance, rewardBalance)
  }, [lpBalance, rewardBalance])

  const formattingFn = useCallback((value) => {
    return displayAmount(fromTokenUnit(value))
  }, [])

  return (
    <section className="keep-only-pool">
      <section className="tile keep-only-pool__overview">
        <section>
          <h2 className="h2--alt text-grey-70">Your KEEP Total Locked</h2>
          <h1 className="text-mint-100 mt-2">
            <CountUp
              end={toTokenUnit(lockedKEEP).toNumber()}
              preserveValue
              decimals={2}
              duration={1}
              formattingFn={formattingFn}
            />
            &nbsp;<span className="h2">KEEP</span>
          </h1>
          <div className="flex row space-between text-grey-40 mt-1">
            <h4>Deposited KEEP tokens</h4>
            <h4 className="self-end">
              <CountUp
                end={toTokenUnit(lpBalance).toNumber()}
                preserveValue
                decimals={2}
                duration={1}
                formattingFn={formattingFn}
              />
              KEEP
            </h4>
          </div>
          <Divider style={{ margin: "0.5rem 0" }} />
          <div className="flex row space-between text-grey-40">
            <h4>Rewarded KEEP tokens</h4>
            <h4 className="self-end">
              <CountUp
                end={toTokenUnit(rewardBalance).toNumber()}
                preserveValue
                decimals={2}
                duration={1}
                formattingFn={formattingFn}
              />
              KEEP
            </h4>
          </div>

          <div className="flex row space-between mt-2">
            <SubmitButton
              className="btn btn-primary btn-lg"
              disabled={!gt(wrappedTokenBalance || 0, 0)}
              onSubmitAction={(awaitingPromise) =>
                addLpTokens(
                  wrappedTokenBalance,
                  liquidityContractName,
                  pool,
                  awaitingPromise
                )
              }
            >
              {gt(lpBalance, 0) ? "add more keep" : "deposit keep"}
            </SubmitButton>
            <SubmitButton
              className="liquidity__withdraw btn btn-secondary btn-lg"
              disabled={!gt(rewardBalance || 0, 0) && !gt(lpBalance || 0, 0)}
              onSubmitAction={(awaitingPromise) =>
                withdrawLiquidityRewards(
                  liquidityContractName,
                  lpBalance,
                  pool,
                  awaitingPromise
                )
              }
            >
              withdraw all
            </SubmitButton>
          </div>
        </section>
        <section
          className={`keep-only-pool__overview__info-tiles liquidity__info${
            gt(lpBalance, 0) ? "" : "--locked"
          }`}
        >
          <div className="liquidity__info-tile bg-mint-10 mb-1">
            <APY.Tooltip className="liquidity__info-tile__tooltip" />
            <APY
              apy={apy}
              isFetching={isAPYFetching}
              className="liquidity__info-tile__title text-mint-100"
            />
            <h6>Estimate of pool apy</h6>
          </div>
          <div className="liquidity__info-tile bg-mint-10">
            <ShareOfPool
              className="liquidity__info-tile__title text-mint-100"
              percentageOfTotalPool={percentageOfTotalPool}
              isFetching={isFetching}
            />
            <h6>your keep rewards</h6>
          </div>
        </section>
      </section>
      <section className="keep-only-pool__icon">
        <Icons.KeepOnlyPool preserveAspectRatio="none" />
      </section>
    </section>
  )
}

export default KeepOnlyPool
