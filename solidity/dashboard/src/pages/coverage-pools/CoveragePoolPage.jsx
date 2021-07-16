import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import PageWrapper from "../../components/PageWrapper"
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
import resourceTooltipProps from "../../constants/tooltips"
import { Column, DataTable } from "../../components/DataTable"
import moment from "moment"
import { SubmitButton } from "../../components/Button"
import { colors } from "../../constants/colors"
import ProgressBar from "../../components/ProgressBar"
import * as Icons from "../../components/Icons"
import Chip from "../../components/Chip"

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
    // covTotalSupply,
    // error,
    estimatedRewards,
    estimatedKeepBalance,
    apy,
    isApyFetching,
    totalAllocatedRewards,
    withdrawalDelay,
    withdrawalTimeout,
    pendingWithdrawals,
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
        modalOptions: { title: "Initiate Deposit" },
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
        modalOptions: { title: "Initiate Deposit" },
        submitBtnText: "withdraw",
        amount,
      },
      InitiateDepositModal
    )
    dispatch(withdrawAssetPool(amount, awaitingPromise))
  }

  const isWithdrawalCooldownOver = (pendingWithdrawalTimestamp) => {
    const currentDate = moment()
    const endOfCooldownDate = moment
      .unix(pendingWithdrawalTimestamp)
      .add(withdrawalDelay, "seconds")

    return currentDate.isAfter(endOfCooldownDate)
  }

  const renderProgressBar = (
    withdrawalDate,
    endOfCooldownDate,
    currentDate
  ) => {
    const progressBarValueInSeconds = currentDate.diff(
      withdrawalDate,
      "seconds"
    )
    const progressBarTotalInSeconds = endOfCooldownDate.diff(
      withdrawalDate,
      "seconds"
    )
    return (
      <ProgressBar
        value={progressBarValueInSeconds}
        total={progressBarTotalInSeconds}
        color={colors.secondary}
        bgColor={colors.bgSecondary}
      >
        <ProgressBar.Inline
          height={20}
          className={"pending-withdrawal__progress-bar"}
        />
      </ProgressBar>
    )
  }

  const renderCooldownStatus = (timestamp) => {
    const withdrawalDate = moment.unix(timestamp)
    const currentDate = moment()
    const endOfCooldownDate = moment
      .unix(timestamp)
      .add(withdrawalDelay, "seconds")
    const days = endOfCooldownDate.diff(currentDate, "days")
    const hours = moment.duration(endOfCooldownDate.diff(currentDate)).hours()
    const minutes = moment
      .duration(endOfCooldownDate.diff(currentDate))
      .minutes()

    let cooldownStatus = <></>
    if (days >= 0 && hours >= 0 && minutes >= 0) {
      cooldownStatus = (
        <>
          {renderProgressBar(withdrawalDate, endOfCooldownDate, currentDate)}
          <div className={"pending-withdrawal__cooldown-time-container"}>
            <Icons.Time
              width="16"
              height="16"
              className="time-icon time-icon--grey-30"
            />
            <span>
              {days}d {hours}h {minutes}m until available
            </span>
          </div>
        </>
      )
    } else {
      cooldownStatus = <Chip text={"cooldown completed"} size="small" />
    }

    return (
      <div className={"pending-withdrawal__cooldown-status"}>
        {cooldownStatus}
      </div>
    )
  }

  const onSubmitBtn = () => {}

  const onCancel = () => {}

  return (
    <PageWrapper title={title} newPage={withNewLabel}>
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

      <section className="coverage-pool__deposit-wrapper">
        <section className="tile coverage-pool__deposit-form">
          <h3>Deposit</h3>
          <DepositForm
            onSubmit={onSubmitDepositForm}
            tokenAmount={keepTokenBalance.value}
            apy={apy}
          />
        </section>

        <section className="tile coverage-pool__share-of-pool">
          <h4 className="text-grey-70 mb-3">Your Share of Pool</h4>

          <OnlyIf condition={shareOfPool <= 0}>
            <div className="text-grey-30 text-center">
              You have no balance yet.&nbsp;
              <br />
              <u>Deposit KEEP</u>&nbsp;to see balance.
            </div>
          </OnlyIf>
          <OnlyIf condition={shareOfPool > 0}>
            <div className="flex column center">
              <TokenAmount amount={estimatedKeepBalance} withSymbol={false} />
              <h4 className="text-mint-100">{KEEP.symbol}</h4>
              <div className="text-grey-40 mt-2">
                <b>{displayPercentageValue(shareOfPool * 100, false)}</b>
                &nbsp;of Pool
              </div>
            </div>
          </OnlyIf>
        </section>

        <section className="tile coverage-pool__rewards">
          <h4 className="text-grey-70 mb-3">Your Rewards</h4>
          <OnlyIf condition={lte(estimatedRewards, 0) && shareOfPool <= 0}>
            <div className="text-grey-30 text-center">
              You have no rewards yet.&nbsp;
              <br />
              <u>Deposit KEEP</u>&nbsp;to see rewards.
            </div>
          </OnlyIf>
          <OnlyIf condition={shareOfPool > 0}>
            <div className="flex column center">
              <TokenAmount amount={estimatedRewards} withSymbol={false} />
              <h4 className="text-mint-100">{KEEP.symbol}</h4>
            </div>
          </OnlyIf>
        </section>

        {/* <HowDoesItWorkBanner />*/}

        <section className="tile coverage-pool__withdraw-wrapper">
          <h3>Available to withdraw</h3>
          <TokenAmount
            wrapperClassName={"coverage-pool__token-amount"}
            amount={covBalance}
            token={covKEEP}
            withIcon
          />
          <WithdrawAmountForm
            onCancel={onCancel}
            submitBtnText="add keep"
            withdrawAmountBalance={covBalance}
            onSubmit={onSubmitWithdrawForm}
            withdrawalDelay={withdrawalDelay}
          />
        </section>
      </section>

      <section className={"tile pending-withdrawal"}>
        <DataTable
          data={pendingWithdrawals}
          itemFieldId="pendingWithdrawalId"
          title="Pending withdrawal"
          withTooltip
          tooltipProps={resourceTooltipProps.pendingWithdrawal}
          noDataMessage="No pending withdrawals."
        >
          <Column
            header="amount"
            field="covAmount"
            renderContent={({ covAmount }) => {
              return <TokenAmount amount={covAmount} />
            }}
          />
          <Column
            header="withdrawal initiated"
            field="timestamp"
            renderContent={({ timestamp }) => {
              const withdrawalDate = moment.unix(timestamp)
              return (
                <div className={"pending-withdrawal__date"}>
                  <span>{withdrawalDate.format("DD-MM-YYYY")}</span>
                  <span>{withdrawalDate.format("HH:mm:ss")}</span>
                </div>
              )
            }}
          />
          <Column
            header="cooldown status"
            field="timestamp"
            tdClassName={"cooldown-status-column"}
            renderContent={({ timestamp }) => {
              return renderCooldownStatus(timestamp)
            }}
          />
          <Column
            header=""
            field="timestamp"
            renderContent={({ timestamp }) => (
              <div className={"pending-withdrawal__button-container"}>
                <SubmitButton
                  className="btn btn-lg btn-primary"
                  onSubmitAction={onSubmitBtn}
                  disabled={!isWithdrawalCooldownOver(timestamp)}
                >
                  claim tokens
                </SubmitButton>
              </div>
            )}
          />
        </DataTable>
      </section>
    </PageWrapper>
  )
}

export default CoveragePoolPage
