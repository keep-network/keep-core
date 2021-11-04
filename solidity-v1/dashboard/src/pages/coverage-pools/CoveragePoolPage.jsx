import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import { DepositForm, MetricsSection } from "../../components/coverage-pools"
import TokenAmount from "../../components/TokenAmount"
import { useWeb3Address } from "../../components/WithWeb3Context"
import OnlyIf from "../../components/OnlyIf"
import {
  fetchTvlRequest,
  fetchCovPoolDataRequest,
  fetchAPYRequest,
} from "../../actions/coverage-pool"
import { useModal } from "../../hooks/useModal"
import { gt } from "../../utils/arithmetics.utils"
import { covKEEP, KEEP } from "../../utils/token.utils"
import { displayPercentageValue } from "../../utils/general.utils"
import WithdrawAmountForm from "../../components/WithdrawAmountForm"
import PendingWithdrawals from "../../components/coverage-pools/PendingWithdrawals"
import Chip from "../../components/Chip"
import ResourceTooltip from "../../components/ResourceTooltip"
import resourceTooltipProps from "../../constants/tooltips"
import { Keep } from "../../contracts"
import { MODAL_TYPES } from "../../constants/constants"

const CoveragePoolPage = ({ title, withNewLabel }) => {
  const { openConfirmationModal, openModal } = useModal()
  const dispatch = useDispatch()
  const {
    totalValueLocked,
    totalValueLockedInUSD,
    isTotalValueLockedFetching,
    shareOfPool,
    covBalance,
    covTokensAvailableToWithdraw,
    covTotalSupply,
    apy,
    isApyFetching,
    totalAllocatedRewards,
    totalCoverageClaimed,
    withdrawalDelay,
    // pendingWithdrawal,
    withdrawalInitiatedTimestamp,
    hasRiskManagerOpenAuctions,
  } = useSelector((state) => state.coveragePool)

  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)

  const address = useWeb3Address()

  const hasCovKEEPTokens = gt(covBalance, 0)

  useEffect(() => {
    dispatch(fetchTvlRequest())
    dispatch(fetchAPYRequest())
  }, [dispatch])

  useEffect(() => {
    if (address) {
      dispatch(fetchCovPoolDataRequest(address))
    }
  }, [dispatch, address])

  const onSubmitDepositForm = async (values, awaitingPromise) => {
    const { tokenAmount } = values
    const amount = KEEP.fromTokenUnit(tokenAmount).toString()
    if (hasRiskManagerOpenAuctions) {
      await openConfirmationModal(MODAL_TYPES.WarningBeforeCovPoolDeposit)
    }
    openModal(MODAL_TYPES.InitiateCovPoolDeposit, {
      amount,
      estimatedBalanceAmountInKeep: Keep.coveragePoolV1.estimatedBalanceFor(
        covBalance,
        covTotalSupply,
        totalValueLocked
      ),
      totalValueLocked,
      covTotalSupply,
    })
  }

  const onSubmitWithdrawForm = async (values, awaitingPromise) => {
    const { withdrawAmount } = values
    const amount = KEEP.fromTokenUnit(withdrawAmount).toString()
    openModal(MODAL_TYPES.InitiateCovPoolWithdraw, {
      totalValueLocked,
      covTotalSupply,
      covBalanceOf: covBalance,
      estimatedBalanceAmountInKeep: Keep.coveragePoolV1.estimatedBalanceFor(
        covBalance,
        covTotalSupply,
        totalValueLocked
      ),
      amount,
    })
    // dispatch(
    //   addAdditionalDataToModal({
    //     componentProps: {
    //       totalValueLocked,
    //       covTotalSupply,
    //       covTokensAvailableToWithdraw,
    //     },
    //   })
    // )
    // if (eq(withdrawalInitiatedTimestamp, 0)) {
    //   await openConfirmationModal(
    //     {
    //       modalOptions: {
    //         title: "Withdraw",
    //         classes: {
    //           modalWrapperClassName: "modal-wrapper__initiate-withdrawal",
    //         },
    //       },
    //       submitBtnText: "withdraw",
    //       amount,
    //       covTotalSupply,
    //       totalValueLocked,
    //       covTokensAvailableToWithdraw,
    //       containerTitle: "You are about to withdraw:",
    //     },
    //     InitiateCovPoolsWithdrawModal
    //   )
    //   dispatch(withdrawAssetPool(amount, awaitingPromise))
    // } else {
    //   const { amount: finalAmount } = await openConfirmationModal(
    //     {
    //       modalOptions: {
    //         title: "Re-initiate withdrawal",
    //         classes: {
    //           modalWrapperClassName: "modal-wrapper__reinitiate-withdrawal",
    //         },
    //       },
    //       submitBtnText: "continue",
    //       pendingWithdrawalBalance: pendingWithdrawal,
    //       initialAmountValue: amount,
    //       covTokensAvailableToWithdraw,
    //       covTotalSupply,
    //       totalValueLocked,
    //       withdrawalDelay,
    //       containerTitle: "You are about to re-initiate this withdrawal:",
    //     },
    //     ReinitiateWithdrawalModal
    //   )
    //   dispatch(withdrawAssetPool(finalAmount, awaitingPromise))
    // }
  }

  return (
    <>
      <MetricsSection
        tvl={totalValueLocked}
        tvlInUSD={totalValueLockedInUSD}
        rewardRate={apy}
        isRewardRateFetching={isApyFetching}
        totalAllocatedRewards={totalAllocatedRewards}
        isTotalAllocatedRewardsFetching={isTotalValueLockedFetching}
        lifetimeCovered={totalCoverageClaimed}
        isLifetimeCoveredFetching={isTotalValueLockedFetching}
      />
      <OnlyIf condition={withdrawalInitiatedTimestamp > 0}>
        <PendingWithdrawals
          covTokensAvailableToWithdraw={covTokensAvailableToWithdraw}
        />
      </OnlyIf>
      <section className="coverage-pool__deposit-wrapper">
        <section className="tile coverage-pool__deposit-form">
          <div className={"flex row center"}>
            <h3>Deposit</h3>
            <ResourceTooltip
              tooltipClassName={"ml-1"}
              {...resourceTooltipProps.covPoolsDeposit}
            />
          </div>
          <DepositForm
            onSubmit={onSubmitDepositForm}
            tokenAmount={keepTokenBalance.value}
            apy={apy}
          />
        </section>

        <section className="tile coverage-pool__balance">
          <div className="coverage-pool__balance-title mb-1">
            <h3>Balance</h3>
            <OnlyIf condition={gt(withdrawalInitiatedTimestamp, 0)}>
              <Chip
                text={`Pending withdrawal`}
                size="small"
                className={"coverage-pool_pending-withdrawal-chip"}
                color="yellow"
              />
            </OnlyIf>
            <OnlyIf condition={hasCovKEEPTokens}>
              <span className={"coverage-pool__share-of-pool text-grey-40"}>
                {displayPercentageValue(shareOfPool * 100, false)} of pool
              </span>
            </OnlyIf>
          </div>
          <OnlyIf condition={!hasCovKEEPTokens}>
            <h4 className="text-center text-grey-50 mt-3">
              You have no share of the pool yet.
              <p>Deposit KEEP to see your balance.</p>
            </h4>
          </OnlyIf>
          <OnlyIf condition={hasCovKEEPTokens}>
            <TokenAmount
              wrapperClassName={"coverage-pool__token-amount"}
              amount={covBalance}
              amountClassName={"h1 text-mint-100"}
              symbolClassName={"h2 text-mint-100"}
              token={covKEEP}
            />
            <TokenAmount
              amount={Keep.coveragePoolV1.estimatedBalanceFor(
                covBalance,
                covTotalSupply,
                totalValueLocked
              )}
              amountClassName={"h3 text-grey-40"}
              symbolClassName={"h3 text-grey-40"}
            />
          </OnlyIf>
        </section>

        <OnlyIf condition={hasCovKEEPTokens}>
          <section className="tile coverage-pool__withdraw-wrapper">
            <div className={"flex row center mb-1"}>
              <h3>Available to withdraw</h3>
              <ResourceTooltip
                tooltipClassName={"ml-1"}
                {...resourceTooltipProps.covPoolsAvailableToWithdraw}
              />
            </div>
            <TokenAmount
              wrapperClassName={"coverage-pool__token-amount"}
              amount={covTokensAvailableToWithdraw}
              amountClassName={"h2 text-mint-100"}
              symbolClassName={"h3 text-mint-100"}
              token={covKEEP}
            />
            <TokenAmount
              wrapperClassName={"coverage-pool__cov-token-amount"}
              amount={Keep.coveragePoolV1.estimatedBalanceFor(
                covTokensAvailableToWithdraw,
                covTotalSupply,
                totalValueLocked
              )}
              amountClassName={"h3 text-grey-40"}
              symbolClassName={"h3 text-grey-40"}
              token={KEEP}
            />
            <WithdrawAmountForm
              submitBtnText={
                gt(withdrawalInitiatedTimestamp, "0")
                  ? "increase withdrawal"
                  : "withdraw"
              }
              withdrawAmount={covTokensAvailableToWithdraw}
              onSubmit={onSubmitWithdrawForm}
              withdrawalDelay={withdrawalDelay}
            />
          </section>
        </OnlyIf>
      </section>
    </>
  )
}

CoveragePoolPage.route = {
  title: "Deposit",
  path: "/coverage-pools/deposit",
  exact: true,
}

export default CoveragePoolPage
