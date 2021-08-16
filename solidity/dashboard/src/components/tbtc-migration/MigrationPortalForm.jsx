import React, { useState } from "react"
import * as Icons from "../Icons"
import FormInput from "../FormInput"
import MaxAmountAddon from "../MaxAmountAddon"
import { formatAmount, normalizeAmount } from "../../forms/form.utils"
import { getErrorsObj } from "../../forms/common-validators"
import Chip from "../Chip"
import { withFormik } from "formik"
import { TBTC } from "../../utils/token.utils"
import { SubmitButton } from "../Button"

const MigrationPortalForm = ({
  mintingFee = 0,
  tbtcV1Balance = 0,
  tbtcV2Balance = 0,
  onSubmit = () => {},
}) => {
  const [from, setFrom] = useState("v1")
  const [to, setTo] = useState("v2")

  const onSwapBtn = (event) => {
    event.preventDefault()
    if (from === "v1") {
      setFrom("v2")
      setTo("v1")
    } else {
      setFrom("v1")
      setTo("v2")
    }
  }

  return (
    <form className="tbtc-migration-portal-form">
      <div className="tbtc-migration-portal-form__inputs-wrapper">
        <div
          className={`tbtc-token-container tbtc-token-container--${from} tbtc-token-container--top`}
        >
          <Chip text="from" size="big" color="primary" />
          <h3 className="mt-1 flex row center">
            <Icons.TBTC />
            &nbsp;tBTC {from}
          </h3>
          <FormInput
            name="amount"
            type="text"
            label="Amount"
            normalize={normalizeAmount}
            format={formatAmount}
            placeholder="0"
            additionalInfoText={`Balance: ${TBTC.displayAmount(
              to === "v1" ? tbtcV1Balance : tbtcV2Balance
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
          <Chip text="to" size="big" color="black" />
          <h3 className="mt-1 flex row center">
            <Icons.TBTC />
            &nbsp;tBTC {to}
          </h3>
          <FormInput
            name="amount"
            type="text"
            label="Amount"
            normalize={normalizeAmount}
            format={formatAmount}
            placeholder="0"
            disabled
            additionalInfoText={`Balance: ${TBTC.displayAmount(
              to === "v1" ? tbtcV1Balance : tbtcV2Balance
            )}`}
          />
        </div>
      </div>

      <p className="text-smaller text-secondary mb-0">
        {`Minting Fee: ${from === "v2" ? mintingFee : 0}`}
      </p>
      <SubmitButton
        className="btn btn-primary btn-lg w-100 mt-1"
        onSubmitAction={onSubmit}
      >
        {from === "v1" ? "upgrade" : "downgrade"}
      </SubmitButton>
    </form>
  )
}

export default withFormik({
  mapPropsToValues: () => ({
    amount: 0,
  }),
  validate: (values, props) => {
    // const { amount } = values
    const errors = {}

    return getErrorsObj(errors)
  },
  displayName: "TBTCMigrationPortalForm",
})(MigrationPortalForm)
