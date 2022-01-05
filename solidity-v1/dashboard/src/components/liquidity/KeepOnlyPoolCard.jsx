import React, { useCallback, useMemo } from "react"
import * as Icons from "../Icons"
import { MODAL_TYPES } from "../../constants/constants"
import CountUp from "react-countup"
import { KEEP } from "../../utils/token.utils"
import { add, gt } from "../../utils/arithmetics.utils"
import Divider from "../Divider"
import { SubmitButton } from "../Button"
import { useModal } from "../../hooks/useModal"
import LiquidityRewardCard from "../LiquidityRewardCard"

const KeepOnlyPoolCard = ({
  poolId = "KEEP_ONLY",
  title,
  MainIcon,
  lpBalance,
  rewardBalance,
  withdrawLiquidityRewards,
}) => {
  const { openConfirmationModal } = useModal()

  const lockedKEEP = useMemo(() => {
    return add(lpBalance, rewardBalance)
  }, [lpBalance, rewardBalance])

  const formattingFn = useCallback((value) => {
    return KEEP.displayAmount(KEEP.fromTokenUnit(value))
  }, [])

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
    [
      withdrawLiquidityRewards,
      lpBalance,
      openConfirmationModal,
      rewardBalance,
      poolId,
    ]
  )

  return (
    <LiquidityRewardCard
      title={title}
      MainIcon={MainIcon}
      lpBalance={lpBalance}
      rewardBalance={rewardBalance}
    >
      <LiquidityRewardCard.InactivePoolBanner
        description={
          "The incentives for this pool have been removed and you can no longer deposit KEEP tokens. You can still withdraw deposited KEEP tokens and rewards that you already earned"
        }
      />
      <div className={"liquidity__reward-balance"}>
        <h4 className={"liquidity__reward-balance__title text-grey-70 mb-1"}>
          Your KEEP pool balance
        </h4>
        <div className={"liquidity__reward-balance_values text-grey-70"}>
          <h3 className={"liquidity__reward-balance_values_label"}>
            <Icons.KeepOutline />
            <span>KEEP</span>
          </h3>
          <h3>
            <CountUp
              end={KEEP.toTokenUnit(lockedKEEP).toNumber()}
              separator={","}
              preserveValue
            />
          </h3>
        </div>
        <div className="flex row space-between text-grey-40 mt-1">
          <h4>Deposited</h4>
          <h4 className="self-end">
            <CountUp
              end={KEEP.toTokenUnit(lpBalance).toNumber()}
              preserveValue
              decimals={2}
              duration={1}
              formattingFn={formattingFn}
            />
          </h4>
        </div>
        <Divider style={{ margin: "0.5rem 0" }} />
        <div className="flex row space-between text-grey-40 mb-2">
          <h4>Rewarded</h4>
          <h4 className="self-end">
            <CountUp
              end={KEEP.toTokenUnit(rewardBalance).toNumber()}
              preserveValue
              decimals={2}
              duration={1}
              formattingFn={formattingFn}
            />
          </h4>
        </div>
      </div>
      <SubmitButton
        className="liquidity__withdraw btn btn-secondary btn-lg w-100"
        disabled={!gt(rewardBalance || 0, 0) && !gt(lpBalance || 0, 0)}
        onSubmitAction={withdrawKEEP}
      >
        withdraw all
      </SubmitButton>
    </LiquidityRewardCard>
  )
}

export default KeepOnlyPoolCard
