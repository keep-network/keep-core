import React from "react"
import { withFormik } from "formik"
import FormInput from "../../components/FormInput"
import { SubmitButton } from "../../components/Button"
import Divider from "../../components/Divider"
import MaxAmountAddon from "../MaxAmountAddon"
import { normalizeFloatingAmount } from "../../forms/form.utils"
import { KEEP } from "../../utils/token.utils"
import TokenAmount from "../TokenAmount"
import { useCustomOnSubmitFormik } from "../../hooks/useCustomOnSubmitFormik"
import {
  validateAmountInRange,
  getErrorsObj,
} from "../../forms/common-validators"
import { lte } from "../../utils/arithmetics.utils"
import useSetMaxAmountToken from "../../hooks/useSetMaxAmountToken"
import { Keep } from "../../contracts"

const AddAmountToWithdrawalForm = ({
  tokenAmount,
  onSubmit,
  initialValue = "0",
  totalValueLocked,
  covTotalSupply,
  ...formikProps
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)
  const onAddonClick = useSetMaxAmountToken("tokenAmount", tokenAmount)

  return (
    <form className="add-amount-to-withdraw-form">
      <div className="add-amount-to-withdraw-form__token-amount-wrapper">
        <h4>Add your available balance?</h4>
        <div className={"add-amount-to-withdraw-form__available-balance"}>
          <TokenAmount
            amount={Keep.coveragePoolV1.estimatedBalanceFor(
              tokenAmount,
              covTotalSupply,
              totalValueLocked
            )}
            wrapperClassName={"add-amount-to-withdraw-form__token-amount"}
            amountClassName={"h3 text-mint-100"}
            symbolClassName={"h3 text-mint-100"}
            token={KEEP}
            withIcon
          />
          <h4
            className={
              "add-amount-to-withdraw-form__cov-token-amount text-grey-70"
            }
          >
            {KEEP.toFormat(KEEP.toTokenUnit(tokenAmount)).toString()} covKEEP
          </h4>
        </div>
        <FormInput
          name="tokenAmount"
          type="text"
          label="Amount"
          normalize={normalizeFloatingAmount}
          inputAddon={
            <MaxAmountAddon onClick={onAddonClick} text="Max Stake" />
          }
        />
      </div>
      <Divider className="divider divider--tile-fluid" />
      <SubmitButton
        className="btn btn-lg btn-primary w-100"
        onSubmitAction={onSubmitBtn}
        disabled={!formikProps.isValid}
      >
        continue
      </SubmitButton>
    </form>
  )
}

export default withFormik({
  validateOnChange: true,
  validateOnBlur: true,
  mapPropsToValues: ({ initialValue }) => ({
    tokenAmount: KEEP.toTokenUnit(initialValue).toString(),
  }),
  validate: (values, props) => {
    const { tokenAmount } = values
    const errors = {}

    // TODO: Remove default 0 value

    if (lte(props.tokenAmount || 0, 0)) {
      errors.tokenAmount = "Insufficient funds"
    } else {
      errors.tokenAmount = validateAmountInRange(
        tokenAmount,
        props.tokenAmount,
        0
      )
    }

    return getErrorsObj(errors)
  },
  displayName: "CovPoolsAddAmountToWithdrawalForm",
})(AddAmountToWithdrawalForm)
