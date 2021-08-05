import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import {
  CheckListBanner,
  HowDoesItWorkBanner,
  DepositForm,
  InitiateDepositModal,
} from "../../components/coverage-pools"
import TokenAmount from "../../components/TokenAmount"
import MetricsTile from "../../components/MetricsTile"
import { APY } from "../../components/liquidity"
import { Skeleton } from "../../components/skeletons"
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
import { lte } from "../../utils/arithmetics.utils"
import { covKEEP, KEEP } from "../../utils/token.utils"
import { displayPercentageValue } from "../../utils/general.utils"
import WithdrawAmountForm from "../../components/WithdrawAmountForm"
import PendingWithdrawals from "../../components/coverage-pools/PendingWithdrawals"
import Chip from "../../components/Chip"
import Tooltip from "../../components/Tooltip"
import * as Icons from "../../components/Icons"
import Divider from "../../components/Divider"
import InitiateCovPoolsWithdrawModal
  from "../../components/coverage-pools/InitiateCovPoolsWithdrawModal";

const CoveragePoolPage = ({ title, withNewLabel }) => {
  const { openConfirmationModal } = useModal()
  const dispatch = useDispatch()
  const {
    totalValueLocked,
    totalValueLockedInUSD,
    isTotalValueLockedFetching,
    // isDataFetching,
    shareOfPool,
    covBalance,
    covTokensAvailableToWithdraw,
    // covTotalSupply,
    // error,
    estimatedRewards,
    estimatedKeepBalance,
    apy,
    isApyFetching,
    totalAllocatedRewards,
    withdrawalDelay,
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
    const amount = KEEP.fromTokenUnit(withdrawAmount)
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
        containerTitle: "You are about to withdraw:"
      },
      InitiateCovPoolsWithdrawModal,
    )
    dispatch(withdrawAssetPool(amount, awaitingPromise))
  }

  const onCancel = () => {}

  return (
    <>
      <CheckListBanner />
      <section className="tile coverage-pool__overview">
        <section className="coverage-pool__overview__tvl">
          <h2 className="h2--alt text-grey-70 mb-1">Total Value Locked</h2>
          <TokenAmount
            amount={totalValueLocked}
            amountClassName="h1 text-mint-100"
            symbolClassName="h2 text-mint-100"
            withIcon
          />
          <h3 className="tvl tvl--usd">
            {`$${totalValueLockedInUSD.toString()} USD`}
          </h3>
        </section>
        <div className="coverage-pool__overview__metrics">
          <section className="metrics__apy">
            <h4 className="text-grey-70 mb-1">Rewards Rate</h4>

            <MetricsTile className="bg-mint-10 mr-2">
              <APY
                apy={apy}
                isFetching={isApyFetching}
                className="text-mint-100"
              />
              <h5 className="text-grey-60">annual</h5>
            </MetricsTile>
          </section>
          <section className="metrics__total-rewards">
            <h4 className="text-grey-70 mb-1">Total Rewards</h4>

            <MetricsTile className="bg-mint-10">
              {isTotalValueLockedFetching ? (
                <Skeleton tag="h2" shining color="grey-10" />
              ) : (
                <TokenAmount
                  amount={totalAllocatedRewards}
                  withIcon
                  withSymbol={false}
                  withMetricSuffix
                />
              )}
              <h5 className="text-grey-60">pool lifetime</h5>
            </MetricsTile>
          </section>
        </div>

        {/* TODO add more metrics according to the Figma vies */}
      </section>

      <PendingWithdrawals covTokensAvailableToWithdraw={covTokensAvailableToWithdraw}/>

      <section className="coverage-pool__deposit-wrapper">
        <section className="tile coverage-pool__deposit-form">
          <h3>Deposit</h3>
          <DepositForm
            onSubmit={onSubmitDepositForm}
            tokenAmount={keepTokenBalance.value}
            apy={apy}
          />
        </section>

        <section className="tile coverage-pool__balance">
          <div className={"coverage-pool__balance-title"}>
            <h3>Balance</h3>
            <OnlyIf condition={withdrawalInitiatedTimestamp > 0}>
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
            amount={covBalance}
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
          <div className={"coverage-pool__deposits-and-earned"}>
            <div className={"coverage-pool__deposit"}>
              <h4 className={"text-grey-70"}>Your deposits &nbsp;</h4>
              <Tooltip
                simple
                delay={0}
                triggerComponent={Icons.MoreInfo}
                className={"lp-balance__tooltip"}
              >
                Your deposits tooltip
              </Tooltip>
              <h4 className={"coverage-pool__deposit-amount text-grey-40"}>
                1,000 KEEP
              </h4>
            </div>
            <Divider style={{ margin: "0.5rem 0" }} />
            <div className={"coverage-pool__earned"}>
              <h4 className={"text-grey-70"}>Earned &nbsp;</h4>
              <Tooltip
                simple
                delay={0}
                triggerComponent={Icons.MoreInfo}
                className={"lp-balance__tooltip"}
              >
                Rewards earned tooltip
              </Tooltip>
              <h4 className={"coverage-pool__earned-amount text-grey-40"}>
                0 KEEP
              </h4>
            </div>
          </div>
        </section>

        <section className="tile coverage-pool__withdraw-wrapper">
          <h3>Available to withdraw</h3>
          <TokenAmount
            wrapperClassName={"coverage-pool__token-amount"}
            amount={covTokensAvailableToWithdraw}
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
            submitBtnText="add keep"
            withdrawAmount={covBalance}
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
