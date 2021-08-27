import React from "react"
import { withFormik, useField } from "formik"
import * as Icons from "../Icons"
import FormInput from "../FormInput"
import MaxAmountAddon from "../MaxAmountAddon"
import { SubmitButton } from "../Button"
import { useCustomOnSubmitFormik } from "../../hooks/useCustomOnSubmitFormik"
import { Keep } from "../../contracts"
import { normalizeFloatingAmount } from "../../forms/form.utils"
import {
  getErrorsObj,
  validateAmountInRange,
} from "../../forms/common-validators"
import { TBTC } from "../../utils/token.utils"
import { gt, sub } from "../../utils/arithmetics.utils"
import { colors } from "../../constants/colors"

const styles = {
  tokenLabel: { margin: "0.5rem 0" },
  fromBox: {
    color: colors.white,
    borderRadius: "0.5rem",
    backgroundColor: colors.black,
    padding: "0 0.5rem",
  },
  toBox: {
    color: colors.black,
    borderRadius: "0.5rem",
    backgroundColor: colors.white,
    border: `1px solid ${colors.black}`,
    padding: "0 0.5rem",
  },
}

const MigrationPortalForm = ({
  mintingFee = 0,
  tbtcV1Balance = 0,
  tbtcV2Balance = 0,
  onSubmit = () => {},
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)

  const [fromField, , fromHelpers] = useField("from")
  const [toField, , toHelpers] = useField("to")
  const from = fromField.value
  const to = toField.value

  const onSwapBtn = (event) => {
    event.preventDefault()
    if (from === "v1") {
      fromHelpers.setValue("v2")
      toHelpers.setValue("v1")
    } else {
      fromHelpers.setValue("v1")
      toHelpers.setValue("v2")
    }
  }

  return (
    <form className="tbtc-migration-portal-form">
      <div className="tbtc-migration-portal-form__inputs-wrapper">
        <div
          className={`tbtc-token-container tbtc-token-container--${from} tbtc-token-container--top`}
        >
          <h5>from</h5>
          <h3 className="flex row center" style={styles.tokenLabel}>
            <Icons.TBTC />
            &nbsp;tBTC&nbsp;
            <h5 style={styles.fromBox}>{from}</h5>
          </h3>
          <FormInput
            name="amount"
            type="text"
            label="Amount"
            normalize={normalizeFloatingAmount}
            placeholder="0"
            additionalInfoText={`Balance: ${TBTC.displayAmount(
              from === "v1" ? tbtcV1Balance : tbtcV2Balance
            )}`}
            inputAddon={<MaxAmountAddon onClick={() => {}} text="Max" />}
          />
        </div>
        <button className="from-to-switcher" onClick={onSwapBtn}>
          <Icons.Swap />
        </button>
        <div
          className={`tbtc-token-container tbtc-token-container--${to} tbtc-token-container--bottom`}
        >
          <h5>to</h5>
          <h3 className="flex row center" style={styles.tokenLabel}>
            <Icons.TBTC />
            &nbsp;tBTC&nbsp;
            <h5 className="self-center" style={styles.toBox}>
              {to}
            </h5>
          </h3>
          <FormInput
            name="amount"
            type="text"
            label="Amount"
            normalize={normalizeFloatingAmount}
            placeholder="0"
            disabled
            additionalInfoText={`Balance: ${TBTC.displayAmount(
              to === "v1" ? tbtcV1Balance : tbtcV2Balance
            )}`}
          />
        </div>
      </div>

      <p className="text-smaller text-secondary mb-0">
        {`Minting Fee: ${from === "v2" ? TBTC.displayAmount(mintingFee) : 0}`}
      </p>
      <SubmitButton
        className="btn btn-primary btn-lg w-100 mt-1"
        onSubmitAction={onSubmitBtn}
      >
        {from === "v1" ? "upgrade" : "downgrade"}
      </SubmitButton>
    </form>
  )
}

export default withFormik({
  mapPropsToValues: () => ({
    amount: 0,
    from: "v1",
    to: "v2",
  }),
  validate: (values, props) => {
    return getMaxAmount(values, props).then((maxAmount) => {
      const errors = {}
      if (gt(TBTC.fromTokenUnit(values.amount).toString(), maxAmount)) {
        errors.amount = "Insufficient funds"
      } else {
        errors.amount = validateAmountInRange(
          values.amount,
          maxAmount,
          1,
          TBTC,
          true
        )
      }

      return getErrorsObj(errors)
    })
  },
  displayName: "TBTCMigrationPortalForm",
})(MigrationPortalForm)

const getMaxAmount = async (values, props) => {
  const { amount, from } = values
  const { mintingFee, tbtcV1Balance, tbtcV2Balance } = props

  if (from === "v1") {
    return tbtcV1Balance
  }

  const unmintFeeFor = await Keep.tBTCV2Migration.unmintFeeFor(
    TBTC.fromTokenUnit(amount).toString(),
    mintingFee
  )

  return sub(tbtcV2Balance, unmintFeeFor).toString()
}
