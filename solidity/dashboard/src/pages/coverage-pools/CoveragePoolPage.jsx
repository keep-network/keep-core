import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import {
  DepositForm,
  InitiateDepositModal,
  MetricsSection,
} from "../../components/coverage-pools"
import TokenAmount from "../../components/TokenAmount"
import { useWeb3Address } from "../../components/WithWeb3Context"
import OnlyIf from "../../components/OnlyIf"
import {
  fetchTvlRequest,
  fetchCovPoolDataRequest,
  depositAssetPool,
  fetchAPYRequest,
  withdrawAssetPool,
} from "../../actions/coverage-pool"
import { useModal } from "../../hooks/useModal"
import { eq, gt } from "../../utils/arithmetics.utils"
import { covKEEP, KEEP } from "../../utils/token.utils"
import { displayPercentageValue } from "../../utils/general.utils"
import WithdrawAmountForm from "../../components/WithdrawAmountForm"
import PendingWithdrawals from "../../components/coverage-pools/PendingWithdrawals"
import Chip from "../../components/Chip"
import InitiateCovPoolsWithdrawModal from "../../components/coverage-pools/InitiateCovPoolsWithdrawModal"
import ReinitiateWithdrawalModal from "../../components/coverage-pools/ReinitiateWithdrawalModal"
import { addAdditionalDataToModal } from "../../actions/modal"
import ResourceTooltip from "../../components/ResourceTooltip"
import resourceTooltipProps from "../../constants/tooltips"
import { Keep } from "../../contracts"

const CoveragePoolPage = ({ title, withNewLabel }) => {
  const { openConfirmationModal } = useModal()
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
    pendingWithdrawal,
    withdrawalInitiatedTimestamp,
  } = useSelector((state) => state.coveragePool)

  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)

  const address = useWeb3Address()

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
    const amount = KEEP.fromTokenUnit(tokenAmount)
    await openConfirmationModal(
      {
        modalOptions: {
          title: "Deposit",
          classes: {
            modalWrapperClassName: "modal-wrapper__initiate-withdrawal",
          },
        },
        submitBtnText: "deposit",
        amount,
      },
      InitiateDepositModal
    )
    dispatch(depositAssetPool(amount, awaitingPromise))
  }

  const onSubmitWithdrawForm = async (values, awaitingPromise) => {
    const { withdrawAmount } = values
    const amount = KEEP.fromTokenUnit(withdrawAmount).toString()
    dispatch(
      addAdditionalDataToModal({
        componentProps: {
          totalValueLocked,
          covTotalSupply,
          covTokensAvailableToWithdraw,
        },
      })
    )
    if (eq(withdrawalInitiatedTimestamp, 0)) {
      await openConfirmationModal(
        {
          modalOptions: {
            title: "Withdraw",
            classes: {
              modalWrapperClassName: "modal-wrapper__initiate-withdrawal",
            },
          },
          submitBtnText: "withdraw",
          amount,
          covTotalSupply,
          totalValueLocked,
          covTokensAvailableToWithdraw,
          containerTitle: "You are about to withdraw:",
        },
        InitiateCovPoolsWithdrawModal
      )
      dispatch(withdrawAssetPool(amount, awaitingPromise))
    } else {
      const { amount: finalAmount } = await openConfirmationModal(
        {
          modalOptions: {
            title: "Re-initiate withdrawal",
            classes: {
              modalWrapperClassName: "modal-wrapper__reinitiate-withdrawal",
            },
          },
          submitBtnText: "continue",
          pendingWithdrawalBalance: pendingWithdrawal,
          initialAmountValue: amount,
          covTokensAvailableToWithdraw,
          covTotalSupply,
          totalValueLocked,
          withdrawalDelay,
          containerTitle: "You are about to re-initiate this withdrawal:",
        },
        ReinitiateWithdrawalModal
      )
      dispatch(withdrawAssetPool(finalAmount, awaitingPromise))
    }
  }

  const onCancel = () => {}

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
      <PendingWithdrawals
        covTokensAvailableToWithdraw={covTokensAvailableToWithdraw}
      />
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
          <div className={"coverage-pool__balance-title"}>
            <h3>Balance</h3>
            <OnlyIf condition={gt(withdrawalInitiatedTimestamp, 0)}>
              <Chip
                text={`Pending withdrawal`}
                size="small"
                className={"coverage-pool_pending-withdrawal-chip"}
                color="yellow"
              />
            </OnlyIf>
            <span className={"coverage-pool__share-of-pool text-grey-40"}>
              {displayPercentageValue(shareOfPool * 100, false)} of pool
            </span>
          </div>
          <TokenAmount
            wrapperClassName={"coverage-pool__token-amount"}
            amount={Keep.coveragePoolV1.estimatedBalanceFor(
              covBalance,
              covTotalSupply,
              totalValueLocked
            )}
            amountClassName={"h1 text-mint-100"}
            symbolClassName={"h2 text-mint-100"}
            token={KEEP}
            withIcon
          />
          <TokenAmount
            wrapperClassName={"coverage-pool__cov-token-amount"}
            amount={covBalance}
            amountClassName={"h3 text-grey-40"}
            symbolClassName={"h3 text-grey-40"}
            token={covKEEP}
            withIcon
            icon={() => {
              return <div style={{ width: "32px", height: "32px" }}></div>
            }}
          />
        </section>

        <section className="tile coverage-pool__withdraw-wrapper">
          <div className={"flex row center"}>
            <h3>Available to withdraw</h3>
            <ResourceTooltip
              tooltipClassName={"ml-1"}
              {...resourceTooltipProps.covPoolsAvailableToWithdraw}
            />
          </div>
          <TokenAmount
            wrapperClassName={"coverage-pool__token-amount"}
            amount={Keep.coveragePoolV1.estimatedBalanceFor(
              covTokensAvailableToWithdraw,
              covTotalSupply,
              totalValueLocked
            )}
            amountClassName={"h2 text-mint-100"}
            symbolClassName={"h3 text-mint-100"}
            token={KEEP}
            withIcon
          />
          <TokenAmount
            wrapperClassName={"coverage-pool__cov-token-amount"}
            amount={covTokensAvailableToWithdraw}
            amountClassName={"h3 text-grey-40"}
            symbolClassName={"h3 text-grey-40"}
            token={covKEEP}
            withIcon
            icon={() => {
              return <div style={{ width: "32px", height: "32px" }}></div>
            }}
          />
          <WithdrawAmountForm
            onCancel={onCancel}
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
