import React from "react"
import { useDispatch } from "react-redux"
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
import { KEEP } from "../../utils/token.utils"
import { useModal } from "../../hooks/useModal"
import { depositAssetPool } from "../../actions/web3"
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
  const mockedData = [
    {
      covAmount: "1000000000000000000000",
      timestamp: "1624425850",
    },
  ]
  const cooldownDurationInDays = 14
  const withdrawAvailableDurationInDays = 3
  const shareOfPool = "0"
  const rewards = "0"

  const { openConfirmationModal } = useModal()
  const dispatch = useDispatch()

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

  const isWithdrawalCooldownOver = (pendingWithdrawal) => {
    const currentDate = moment()
    const endOfCooldownDate = moment
      .unix(pendingWithdrawal.timestamp)
      .add(cooldownDurationInDays, "days")

    return currentDate.isAfter(endOfCooldownDate)
  }

  const renderProgressBar = (
    withdrawalDate,
    endOfCooldownDate,
    currentDate
  ) => {
    const progressBarValueInMinutes = currentDate.diff(
      withdrawalDate,
      "minutes"
    )
    const progressBarTotalInMinutes = endOfCooldownDate.diff(
      withdrawalDate,
      "minutes"
    )
    return (
      <ProgressBar
        value={progressBarValueInMinutes}
        total={progressBarTotalInMinutes}
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
      .add(cooldownDurationInDays, "days")
    const days = endOfCooldownDate.diff(currentDate, "days")
    const hours = endOfCooldownDate.diff(currentDate, "hours") % 24
    const minutes = endOfCooldownDate.diff(currentDate, "minutes") % 60

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

  const onMaxAmountClick = () => {}

  const onCancel = () => {}

  const onSubmit = () => {}

  return (
    <PageWrapper title={title} newPage={withNewLabel}>
      <CheckListBanner />

      <section className="tile coverage-pool__overview">
        <section className="coverage-pool__overview__tvl">
          <h2 className="h2--alt text-grey-70 mb-1">Total Value Locked</h2>
          <TokenAmount
            amount="900000000000000000000000000"
            amountClassName="h1 text-mint-100"
            symbolClassName="h2 text-mint-100"
            withIcon
          />
        </section>
        <section className="coverage-pool__overview__apy">
          <h3 className="text-grey-70 mb-1">Pool APY</h3>
          <section className="apy__values">
            <MetricsTile className="bg-mint-10">
              <APY apy="0.15" className="text-mint-100" />
              <h5 className="text-grey-60">weekly</h5>
            </MetricsTile>
            <MetricsTile className="bg-mint-10">
              <APY apy="0.50" className="text-mint-100 " />
              <h5 className="text-grey-60">monthly</h5>
            </MetricsTile>
            <MetricsTile className="bg-mint-10">
              <APY apy="1.40" className="text-mint-100" />
              <h5 className="text-grey-60">annual</h5>
            </MetricsTile>
          </section>
        </section>
      </section>

      <section className="coverage-pool__deposit-wrapper">
        <section className="tile coverage-pool__deposit-form">
          <h3>Deposit</h3>
          <DepositForm onSubmit={onSubmitDepositForm} />
        </section>

        <section className="tile coverage-pool__share-of-pool">
          <h4 className="text-grey-70">Your Share of Pool</h4>
        </section>

        <section className="tile coverage-pool__rewards">
          <h4 className="text-grey-70">Your Rewards</h4>
        </section>

        {/*<HowDoesItWorkBanner />*/}

        <section className="tile coverage-pool__withdraw-wrapper">
          <h3>Available to withdraw</h3>
          <TokenAmount
            wrapperClassName={"coverage-pool__token-amount"}
            amount={"100000000000000000000000"}
            withIcon
          />
          <WithdrawAmountForm
            onCancel={onCancel}
            submitBtnText="add keep"
            availableAmount={"10000000000000000"}
            currentAmount={"10000000000000000"}
            onBtnClick={onSubmit}
          />
        </section>
      </section>

      <section className={"tile pending-withdrawal"}>
        <DataTable
          data={mockedData}
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
            renderContent={() => (
              <div className={"pending-withdrawal__button-container"}>
                <SubmitButton
                  className="btn btn-lg btn-primary"
                  onSubmitAction={onSubmitBtn}
                  disabled={!isWithdrawalCooldownOver(mockedData[0])}
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
