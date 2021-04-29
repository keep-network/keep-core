import React, { useMemo, useCallback } from "react"
import CountUp from "react-countup"
import { withFormik } from "formik"
import Divider from "./Divider"
import { SubmitButton } from "./Button"
import * as Icons from "./Icons"
import { APY } from "./liquidity"
import { gt, add, lte } from "../utils/arithmetics.utils"
import { KEEP } from "../utils/token.utils"
import {
  normalizeAmount,
  formatAmount as formatFormAmount,
} from "../forms/form.utils.js"
import MaxAmountAddon from "./MaxAmountAddon"
import useSetMaxAmountToken from "../hooks/useSetMaxAmountToken"
import AvailableTokenForm from "./AvailableTokenForm"
import { validateAmountInRange, getErrorsObj } from "../forms/common-validators"
import { useModal } from "../hooks/useModal"
import TokenAmount from "./TokenAmount"
import MetricsTile from "./MetricsTile"
import RewardMultiplier from "./liquidity/RewardMultiplier"

const KeepOnlyPool = ({
  apy,
  lpBalance,
  rewardBalance,
  wrappedTokenBalance,
  isFetching,
  isAPYFetching,
  addLpTokens,
  withdrawLiquidityRewards,
  liquidityContractName,
  pool,
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
        {
          modalOptions: { title: "Deposit KEEP" },
          availableAmount: wrappedTokenBalance,
        },
        AddKEEPFormik
      )

      addLpTokens(
        KEEP.fromTokenUnit(amount).toString(),
        liquidityContractName,
        pool,
        awaitingPromise
      )
    },
    [
      addLpTokens,
      liquidityContractName,
      pool,
      openConfirmationModal,
      wrappedTokenBalance,
    ]
  )

  const withdrawKEEP = useCallback(
    async (awaitingPromise) => {
      const { amount } = await openConfirmationModal(
        {
          modalOptions: { title: "Withdraw Locked KEEP" },
          availableAmount: lpBalance,
          rewardedAmount: rewardBalance,
        },
        WithdrawKEEPFormik
      )

      withdrawLiquidityRewards(
        liquidityContractName,
        KEEP.fromTokenUnit(amount).toString(),
        pool,
        awaitingPromise
      )
    },
    [
      withdrawLiquidityRewards,
      lpBalance,
      openConfirmationModal,
      pool,
      liquidityContractName,
      rewardBalance,
    ]
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

          <div className="flex row space-between mt-2">
            <SubmitButton
              className="btn btn-primary btn-lg"
              disabled={!gt(wrappedTokenBalance || 0, 0)}
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

const AddKEEPForm = (props) => {
  const { availableAmount, onCancel, ...formikProps } = props
  const setMaxAmount = useSetMaxAmountToken("amount", availableAmount)

  return (
    <>
      <h3 className="mb-1">Amount available to deposit.</h3>
      <TokenAmount amount={availableAmount} withIcon withMetricSuffix />
      <AvailableTokenForm
        onSubmit={formikProps.handleSubmit}
        onCancel={onCancel}
        submitBtnText="deposit keep"
        formInputProps={{
          name: "amount",
          type: "text",
          label: "Deposit",
          normalize: normalizeAmount,
          format: formatFormAmount,
          placeholder: "0",
          inputAddon: <MaxAmountAddon onClick={setMaxAmount} text="Max KEEP" />,
        }}
        {...formikProps}
      />
    </>
  )
}

const WithdrawKEEPForm = (props) => {
  const { availableAmount, rewardedAmount, onCancel, ...formikProps } = props
  const setMaxAmount = useSetMaxAmountToken("amount", availableAmount)

  return (
    <>
      <h3 className="mb-1">Amount available to withdraw.</h3>
      <div className="flex row mb-2">
        <AmountTile title="deposited" amount={availableAmount} />
        <AmountTile
          title="rewarded"
          amount={rewardedAmount}
          icon={Icons.Rewards}
        />
      </div>
      <AvailableTokenForm
        onSubmit={formikProps.handleSubmit}
        onCancel={onCancel}
        submitBtnText="withdraw keep"
        formInputProps={{
          name: "amount",
          type: "text",
          label: "Withdraw",
          normalize: normalizeAmount,
          format: formatFormAmount,
          placeholder: "0",
          inputAddon: <MaxAmountAddon onClick={setMaxAmount} text="Max KEEP" />,
        }}
        {...formikProps}
      />
    </>
  )
}

const styles = {
  amountTileWrapper: {
    justifyContent: "flex-start",
    flexGrow: "1",
    padding: "0.5rem",
    height: "auto",
  },
}
const AmountTile = ({ amount, title, icon }) => {
  return (
    <MetricsTile
      className="bg-grey-10 self-start"
      style={styles.amountTileWrapper}
    >
      <h5 className="text-grey-40 text-left mb-1">{title}</h5>
      <TokenAmount
        wrapperClassName="mb-1"
        amount={amount}
        icon={icon}
        withIcon
        withMetricSuffix
      />
    </MetricsTile>
  )
}

const commonFormikOptions = {
  mapPropsToValues: () => ({
    amount: "0",
  }),
  validate: ({ amount }, { availableAmount }) => {
    const errors = {}
    const minAmount = KEEP.fromTokenUnit(1)
    if (lte(availableAmount || 0, 0)) {
      errors.amount = "Insufficient funds"
    } else {
      errors.amount = validateAmountInRange(amount, availableAmount, minAmount)
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props }) => props.onBtnClick(values),
}
const WithdrawKEEPFormik = withFormik({
  ...commonFormikOptions,
  displayName: "WithdrawKEEPFormik",
})(WithdrawKEEPForm)

const AddKEEPFormik = withFormik({
  ...commonFormikOptions,
  displayName: "AddKEEPFormik",
})(AddKEEPForm)
