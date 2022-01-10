import React from "react"
import { FormInputBase } from "../../components/FormInput"
import MaxAmountAddon from "../MaxAmountAddon"
import {
  formatFloatingAmount,
  normalizeFloatingAmount,
} from "../../forms/form.utils"
import { covKEEP, KEEP } from "../../utils/token.utils"
import TokenAmount from "../TokenAmount"
import { useSetMaxAmountToken } from "../../hooks/useSetMaxAmountToken"
import { Keep } from "../../contracts"
import { COV_POOLS_FORMS_MAX_DECIMAL_PLACES } from "../../constants/constants"

const AddAmountToWithdrawalForm = ({
  tokenAmount,
  onSubmit,
  totalValueLocked,
  covTotalSupply,
  setMaxAmount,
  inputProps,
}) => {
  const onAddonClick = useSetMaxAmountToken(
    inputProps.name,
    tokenAmount,
    setMaxAmount,
    covKEEP,
    covKEEP.decimals
  )

  return (
    <form className="add-amount-to-withdraw-form" onSubmit={onSubmit}>
      <div className="add-amount-to-withdraw-form__token-amount-wrapper">
        <h4>Add your available balance?</h4>
        <div className={"add-amount-to-withdraw-form__available-balance"}>
          <TokenAmount
            amount={tokenAmount}
            wrapperClassName={"add-amount-to-withdraw-form__token-amount"}
            amountClassName={"h3 text-mint-100"}
            symbolClassName={"h3 text-mint-100"}
            token={covKEEP}
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
        <FormInputBase
          {...inputProps}
          type="text"
          label="Amount"
          placeholder="0"
          normalize={(value) => {
            return normalizeFloatingAmount(
              value,
              COV_POOLS_FORMS_MAX_DECIMAL_PLACES
            )
          }}
          format={(value) => {
            return formatFloatingAmount(
              value,
              COV_POOLS_FORMS_MAX_DECIMAL_PLACES
            )
          }}
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
    </form>
  )
}

export default AddAmountToWithdrawalForm
