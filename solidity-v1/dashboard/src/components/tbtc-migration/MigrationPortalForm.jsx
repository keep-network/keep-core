import React from "react"
import { withFormik, useField, useFormikContext } from "formik"
import * as Icons from "../Icons"
import FormInput from "../FormInput"
import MaxAmountAddon from "../MaxAmountAddon"
import { Keep } from "../../contracts"
import { normalizeFloatingAmount } from "../../forms/form.utils"
import {
  getErrorsObj,
  validateAmountInRange,
} from "../../forms/common-validators"
import { TBTC } from "../../utils/token.utils"
import { sub } from "../../utils/arithmetics.utils"
import { colors } from "../../constants/colors"
import { TBTC_TOKEN_VERSION } from "../../constants/constants"
import Button from "../Button"

const styles = {
  tokenLabel: { margin: "0.5rem 0" },
  fromBox: {
    color: colors.white,
    borderRadius: "0.5rem",
    backgroundColor: colors.night,
    padding: "0 0.5rem",
  },
  toBox: {
    color: colors.night,
    borderRadius: "0.5rem",
    backgroundColor: colors.night,
    border: `1px solid ${colors.night}`,
    padding: "0 0.5rem",
  },
}

const MigrationPortalForm = ({
  mintingFee = 0,
  tbtcV1Balance = 0,
  tbtcV2Balance = 0,
  ...formikProps
}) => {
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
      <Button
        className="btn btn-primary btn-lg w-100 mt-1"
        onClick={formikProps.handleSubmit}
      >
        {from === TBTC_TOKEN_VERSION.v1 ? "upgrade" : "downgrade"}
      </Button>
    </form>
  )
}

export default withFormik({
  mapPropsToValues: () => ({
    amount: 0,
    from: TBTC_TOKEN_VERSION.v1,
    to: TBTC_TOKEN_VERSION.v2,
  }),
  handleSubmit: (values, { props }) => {
    props.onSubmit(values)
  },
  validate: (values, props) => {
    return getMaxAmount(values, props).then(
      ({ tokenBalance, maxTokenBalance }) => {
        const errors = {}
        const formattedValue = TBTC.fromTokenUnit(values.amount)

        if (formattedValue.gt(tokenBalance)) {
          errors.amount = "Insufficient funds"
        } else if (
          values.from === TBTC_TOKEN_VERSION.v2 &&
          formattedValue.gt(maxTokenBalance)
        ) {
          errors.amount = "You don't have enough funds to cover minting fee."
        } else {
          errors.amount = validateAmountInRange(
            values.amount,
            maxTokenBalance,
            TBTC.fromTokenUnit(TBTC.toTokenUnit(1)).toString(), // 1 wei,
            TBTC,
            true
          )
        }

        return getErrorsObj(errors)
      }
    )
  },
  displayName: "TBTCMigrationPortalForm",
})(MigrationPortalForm)

const getMaxAmount = async (values, props) => {
  const { amount, from } = values
  const { mintingFee, tbtcV1Balance, tbtcV2Balance } = props

  if (from === TBTC_TOKEN_VERSION.v1) {
    return { tokenBalance: tbtcV1Balance, maxTokenBalance: tbtcV1Balance }
  }

  const unmintFeeFor = await Keep.tBTCV2Migration.unmintFeeFor(
    TBTC.fromTokenUnit(amount).toString(),
    mintingFee
  )

  return {
    tokenBalance: tbtcV2Balance,
    maxTokenBalance: sub(tbtcV2Balance, unmintFeeFor).toString(),
    unmintFee: unmintFeeFor,
  }
}
