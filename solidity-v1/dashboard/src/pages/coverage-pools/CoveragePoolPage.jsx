import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import {
  DepositForm,
  MetricsSection,
  WithdrawAmountForm,
} from "../../components/coverage-pools"
import TokenAmount from "../../components/TokenAmount"
import { useWeb3Address } from "../../components/WithWeb3Context"
import OnlyIf from "../../components/OnlyIf"
import {
  fetchTvlRequest,
  fetchCovPoolDataRequest,
  fetchAPYRequest,
  increaseWithdrawal,
} from "../../actions/coverage-pool"
import { useModal } from "../../hooks/useModal"
import { gt, eq } from "../../utils/arithmetics.utils"
import { covKEEP, KEEP } from "../../utils/token.utils"
import { displayPercentageValue } from "../../utils/general.utils"
import PendingWithdrawals from "../../components/coverage-pools/PendingWithdrawals"
import Chip from "../../components/Chip"
import ResourceTooltip from "../../components/ResourceTooltip"
import resourceTooltipProps from "../../constants/tooltips"
import { Keep } from "../../contracts"
import { MODAL_TYPES } from "../../constants/constants"
import { CoveragePoolV1ExchangeRate } from "../../components/coverage-pools/ExchangeRate"
import * as Icons from "../../components/Icons"

const CoveragePoolPage = () => {
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

  const onSubmitDepositForm = async (values) => {
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

  const onSubmitWithdrawForm = async (values) => {
    const { withdrawAmount } = values
    const amount = KEEP.fromTokenUnit(withdrawAmount).toString()
    if (eq(withdrawalInitiatedTimestamp, 0)) {
      openModal(MODAL_TYPES.InitiateCovPoolWithdraw, {
        totalValueLocked,
        covTotalSupply,
        covBalanceOf: covBalance,
        amount,
      })
    } else {
      dispatch(increaseWithdrawal(amount))
    }
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
              <span className={"coverage-pool__share-of-pool text-grey-70"}>
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
              wrapperClassName={"mb-2"}
            />
          </OnlyIf>
          <span className={"text-grey-60 flex row"}>
            <Icons.KeepSymbol
              width={20}
              height={20}
              style={{ marginRight: "0.5rem", marginTop: "0.3rem" }}
            />{" "}
            <span>
              Add the covKEEP&nbsp;
              <a
                target="_blank"
                rel="noopener noreferrer"
                href={"https://google.com"}
                className={"no-arrow"}
              >
                token address
              </a>
              &nbsp;to your Ethereum wallet to view balance in the wallet.
            </span>
          </span>
        </section>

        <OnlyIf condition={hasCovKEEPTokens}>
          <section className="tile coverage-pool__withdraw-wrapper">
            <div className={"flex row center mb-1"}>
              <h3>Available to withdraw</h3>
              <ResourceTooltip
                tooltipClassName={"ml-1"}
                {...resourceTooltipProps.covPoolsAvailableToWithdraw}
              />
              <CoveragePoolV1ExchangeRate
                covToken={covKEEP}
                collateralToken={KEEP}
                covTotalSupply={covTotalSupply}
                totalValueLocked={totalValueLocked}
                className="ml-a text-grey-70"
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
