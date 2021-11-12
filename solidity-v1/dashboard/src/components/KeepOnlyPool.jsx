import React, { useMemo, useCallback } from "react"
import CountUp from "react-countup"
import Divider from "./Divider"
import { SubmitButton } from "./Button"
import * as Icons from "./Icons"
import { APY } from "./liquidity"
import { gt, add } from "../utils/arithmetics.utils"
import { KEEP } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import MetricsTile from "./MetricsTile"
import RewardMultiplier from "./liquidity/RewardMultiplier"
import Banner from "./Banner"
import { LINK, MODAL_TYPES } from "../constants/constants"

const poolId = "KEEP_ONLY"
const KeepOnlyPool = ({
  apy,
  lpBalance,
  rewardBalance,
  wrappedTokenBalance,
  isFetching,
  isAPYFetching,
  addLpTokens,
  withdrawLiquidityRewards,
  rewardMultiplier,
}) => {
  const { openConfirmationModal } = useModal()

  const lockedKEEP = useMemo(() => {
    return add(lpBalance, rewardBalance)
  }, [lpBalance, rewardBalance])

  const formattingFn = useCallback((value) => {
    return KEEP.displayAmount(KEEP.fromTokenUnit(value))
  }, [])

  const addKEEP = useCallback(
    async (awaitingPromise) => {
      const { amount } = await openConfirmationModal(
        MODAL_TYPES.KeepOnlyPoolAddKeep,
        {
          availableAmount: wrappedTokenBalance,
        }
      )

      addLpTokens(
        poolId,
        KEEP.fromTokenUnit(amount).toString(),
        awaitingPromise
      )
    },
    [addLpTokens, openConfirmationModal, wrappedTokenBalance]
  )

  const withdrawKEEP = useCallback(
    async (awaitingPromise) => {
      const { amount } = await openConfirmationModal(
        MODAL_TYPES.KeepOnlyPoolWithdrawKeep,
        {
          availableAmount: lpBalance,
          rewardedAmount: rewardBalance,
        }
      )

      withdrawLiquidityRewards(
        poolId,
        KEEP.fromTokenUnit(amount).toString(),
        awaitingPromise
      )
    },
    [withdrawLiquidityRewards, lpBalance, openConfirmationModal, rewardBalance]
  )

  return (
    <section className="keep-only-pool">
      <section className="tile keep-only-pool__overview">
        <section>
          <h2 className="h2--alt text-grey-70">Your KEEP Total Locked</h2>
          <h1 className="text-mint-100 mt-2">
            <CountUp
              end={KEEP.toTokenUnit(lockedKEEP).toNumber()}
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
                end={KEEP.toTokenUnit(lpBalance).toNumber()}
                preserveValue
                decimals={2}
                duration={1}
                formattingFn={formattingFn}
              />
              &nbsp;KEEP
            </h4>
          </div>
          <Divider style={{ margin: "0.5rem 0" }} />
          <div className="flex row space-between text-grey-40">
            <h4>Rewarded KEEP tokens</h4>
            <h4 className="self-end">
              <CountUp
                end={KEEP.toTokenUnit(rewardBalance).toNumber()}
                preserveValue
                decimals={2}
                duration={1}
                formattingFn={formattingFn}
              />
              &nbsp;KEEP
            </h4>
          </div>
          <Banner
            className="liquidity-info-banner liquidity-info-banner--warning mt-2"
            style={{ maxWidth: "100%" }}
          >
            <Banner.Icon
              icon={Icons.Warning}
              className="liquidity-info-banner__icon liquidity-info-banner__icon--warning"
            />
            <div className="liquidity-info-banner__content">
              <Banner.Title className="text-grey-60">
                Incentives removed
              </Banner.Title>
              <Banner.Description className="liquidity-info-banner__content__description text-grey-60">
                The incentives for this pool has been removed and you can no
                longer deposit the KEEP tokens. You can still withdraw deposited
                KEEP tokens and rewards that you already earned.&nbsp;
                <a
                  target="_blank"
                  rel="noopener noreferrer"
                  href={LINK.proposals.shiftingIncentivesToCoveragePools}
                  className="text-link text-grey-60"
                >
                  More info
                </a>
              </Banner.Description>
            </div>
          </Banner>
          <div className="flex row space-between mt-2">
            <SubmitButton
              className="btn btn-primary btn-lg"
              onSubmitAction={addKEEP}
            >
              {gt(lpBalance, 0) ? "add more keep" : "deposit keep"}
            </SubmitButton>
            <SubmitButton
              className="liquidity__withdraw btn btn-secondary btn-lg"
              disabled={!gt(rewardBalance || 0, 0) && !gt(lpBalance || 0, 0)}
              onSubmitAction={withdrawKEEP}
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
          <MetricsTile className="liquidity__info-tile bg-mint-10 mb-1">
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
            <MetricsTile.Tooltip className="liquidity__info-tile__tooltip">
              <RewardMultiplier.TooltipContent />
            </MetricsTile.Tooltip>
            <RewardMultiplier
              rewardMultiplier={rewardMultiplier}
              className="liquidity__info-tile__title text-mint-100"
            />
            <h6>reward multiplier</h6>
          </MetricsTile>
        </section>
      </section>
      <section className="keep-only-pool__icon" />
    </section>
  )
}

export default KeepOnlyPool
