import React from "react"
import { withFormik } from "formik"
import FormInput from "../../components/FormInput"
import { SubmitButton } from "../../components/Button"
import Divider from "../../components/Divider"
import MaxAmountAddon from "../MaxAmountAddon"
import {
  formatFloatingAmount,
  normalizeFloatingAmount,
} from "../../forms/form.utils"
import { covKEEP, KEEP } from "../../utils/token.utils"
import TokenAmount from "../TokenAmount"
import { useCustomOnSubmitFormik } from "../../hooks/useCustomOnSubmitFormik"
import {
  validateAmountInRange,
  getErrorsObj,
} from "../../forms/common-validators"
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
  const onAddonClick = useSetMaxAmountToken(
    "tokenAmount",
    tokenAmount,
    KEEP,
    KEEP.decimals
  )

  return (
    <form className="add-amount-to-withdraw-form">
      <div className="add-amount-to-withdraw-form__token-amount-wrapper">
        <h4>Add your available balance?</h4>
        <div className={"add-amount-to-withdraw-form__available-balance"}>
          <TokenAmount
            amount={tokenAmount}
            wrapperClassName={"add-amount-to-withdraw-form__token-amount"}
            amountClassName={"h3 text-mint-100"}
            symbolClassName={"h3 text-mint-100"}
            token={covKEEP}
            withIcon
          />
          <TokenAmount
            amount={Keep.coveragePoolV1.estimatedBalanceFor(
              tokenAmount,
              covTotalSupply,
              totalValueLocked
            )}
            wrapperClassName={"add-amount-to-withdraw-form__cov-token-amount"}
            amountClassName={"h4 text-grey-70"}
            symbolClassName={"h4 text-grey-70"}
            token={KEEP}
          />
        </div>
        <FormInput
          name="tokenAmount"
          type="text"
          label="Amount"
          normalize={normalizeFloatingAmount}
          format={formatFloatingAmount}
          inputAddon={
            <MaxAmountAddon onClick={onAddonClick} text="Max Amount" />
          }
          leftIconComponent={
            <span className={"form-input__left-icon__cov-keep-amount"}>
              covKEEP
            </span>
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

    errors.tokenAmount = validateAmountInRange(
      tokenAmount,
      props.tokenAmount,
      0,
      covKEEP
    )

    return getErrorsObj(errors)
  },
  displayName: "CovPoolsAddAmountToWithdrawalForm",
})(AddAmountToWithdrawalForm)
