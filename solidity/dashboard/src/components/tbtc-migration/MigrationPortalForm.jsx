import React from "react"
import { withFormik, useField, useFormikContext } from "formik"
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
import { TBTC_TOKEN_VERSION } from "../../constants/constants"

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
  const { setFieldValue } = useFormikContext()

  const [fromField, , fromHelpers] = useField("from")
  const [toField, , toHelpers] = useField("to")
  const from = fromField.value
  const to = toField.value

  const onSwapBtn = (event) => {
    event.preventDefault()
    if (from === TBTC_TOKEN_VERSION.v1) {
      fromHelpers.setValue(TBTC_TOKEN_VERSION.v2)
      toHelpers.setValue(TBTC_TOKEN_VERSION.v1)
    } else {
      fromHelpers.setValue(TBTC_TOKEN_VERSION.v1)
      toHelpers.setValue(TBTC_TOKEN_VERSION.v2)
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
            <span className="h5" style={styles.fromBox}>
              {from}
            </span>
          </h3>
          <FormInput
            name="amount"
            type="text"
            label="Amount"
            normalize={normalizeFloatingAmount}
            placeholder="0"
            additionalInfoText={`Balance: ${TBTC.displayAmount(
              from === TBTC_TOKEN_VERSION.v1 ? tbtcV1Balance : tbtcV2Balance
            )}`}
            inputAddon={
              <MaxAmountAddon
                onClick={() => {
                  setFieldValue(
                    "amount",
                    TBTC.toTokenUnit(
                      from === TBTC_TOKEN_VERSION.v1
                        ? tbtcV1Balance
                        : tbtcV2Balance
                    ).toString()
                  )
                }}
                text="Max"
              />
            }
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
            <span className="h5" style={styles.toBox}>
              {to}
            </span>
          </h3>
          <FormInput
            name="amount"
            type="text"
            label="Amount"
            normalize={normalizeFloatingAmount}
            placeholder="0"
            disabled
            additionalInfoText={`Balance: ${TBTC.displayAmount(
              to === TBTC_TOKEN_VERSION.v1 ? tbtcV1Balance : tbtcV2Balance
            )}`}
          />
        </div>
      </div>

      <p className="text-smaller text-secondary mb-0">
        {`Minting Fee: ${
          from === TBTC_TOKEN_VERSION.v2 ? TBTC.displayAmount(mintingFee) : 0
        }`}
      </p>
      <SubmitButton
        className="btn btn-primary btn-lg w-100 mt-1"
        onSubmitAction={onSubmitBtn}
      >
        {from === TBTC_TOKEN_VERSION.v1 ? "upgrade" : "downgrade"}
      </SubmitButton>
    </form>
  )
}

export default withFormik({
  mapPropsToValues: () => ({
    amount: 0,
    from: TBTC_TOKEN_VERSION.v1,
    to: TBTC_TOKEN_VERSION.v2,
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

  if (from === TBTC_TOKEN_VERSION.v1) {
    return tbtcV1Balance
  }

  const unmintFeeFor = await Keep.tBTCV2Migration.unmintFeeFor(
    TBTC.fromTokenUnit(amount).toString(),
    mintingFee
  )

  return sub(tbtcV2Balance, unmintFeeFor).toString()
}
