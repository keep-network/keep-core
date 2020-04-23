import React, { useCallback, useEffect } from "react"
import { SubmitButton } from "./Button"
import FormInput from "./FormInput"
import { withFormik, useFormikContext } from "formik"
import {
  validateAmountInRange,
  validateEthAddress,
  getErrorsObj,
  validateRequiredValue,
} from "../forms/common-validators"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import { displayAmount, formatAmount } from "../utils/general.utils"
import ProgressBar from "./ProgressBar"
import { colors } from "../constants/colors"
import Dropdown from "./Dropdown"
import SelectedGrantDropdown from "./SelectedGrantDropdown"
import {
  normalizeAmount,
  formatAmount as formatFormAmount,
} from "../forms/form.utils.js"
import { lte } from "../utils/arithmetics.utils"

const DelegateStakeForm = ({
  onSubmit,
  minStake,
  keepBalance,
  grants,
  ...formikProps
}) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)
  const stakeTokensValue = formatAmount(formikProps.values.stakeTokens)
  const {
    values: {
      selectedGrant: { id: selectedGrantId },
    },
    setFieldValue,
  } = formikProps

  useEffect(() => {
    if (!selectedGrantId && grants.length > 0) {
      setFieldValue("selectedGrant", grants[0], true)
    }
  }, [grants, selectedGrantId, setFieldValue])

  const getContextBalance = () => {
    const {
      values: { context, selectedGrant },
    } = formikProps

    if (context === "granted") {
      return selectedGrant.availableToStake || "0"
    }
    return keepBalance
  }

  const isGrantContext = () => {
    const {
      values: { context },
    } = formikProps
    return context === "granted"
  }

  return (
    <form className="delegate-stake-form flex column">
      <div className="flex row">
        <div className="stake-token-section flex column flex-1">
          <div className="text-big text-black">Select Token Amount</div>
          <ContextSwitch />
          <div className="input-wrapper">
            {isGrantContext() && (
              <Dropdown
                onSelect={(selectedGrant) =>
                  formikProps.setFieldValue(
                    "selectedGrant",
                    selectedGrant,
                    true
                  )
                }
                options={grants}
                valuePropertyName="id"
                labelPropertyName="id"
                selectedItem={formikProps.values.selectedGrant}
                labelPrefix="Grant ID"
                noItemSelectedText="Select Grant"
                label={`Choose Grant (${grants.length})`}
                selectedItemComponent={
                  <SelectedGrantDropdown
                    grant={formikProps.values.selectedGrant}
                  />
                }
              />
            )}
            <FormInput
              name="stakeTokens"
              type="text"
              label="Token Amount"
              normalize={normalizeAmount}
              format={formatFormAmount}
            />
            <ProgressBar
              total={getContextBalance()}
              items={[{ value: stakeTokensValue, color: colors.primary }]}
            />
            <div className="text-small text-grey-50">
              {displayAmount(getContextBalance())} KEEP available
            </div>
            <div className="text-smaller text-grey-30 mb-1">
              Min Stake: {displayAmount(minStake)} KEEP
            </div>
          </div>
        </div>
        <div className="addresses-section flex column flex-1 self-start">
          <div className="text-big text-black mb-1">Enter Addresses</div>
          <FormInput
            name="beneficiaryAddress"
            type="text"
            label="Beneficiary Address"
          />
          <FormInput
            name="operatorAddress"
            type="text"
            label="Operator Address"
          />
          <FormInput
            name="authorizerAddress"
            type="text"
            label="Authorizer Address"
          />
        </div>
      </div>
      <div>
        <SubmitButton
          className="btn btn-primary btn-large"
          type="submit"
          onSubmitAction={onSubmitBtn}
          withMessageActionIsPending={false}
          triggerManuallyFetch={true}
        >
          delegate stake
        </SubmitButton>
      </div>
    </form>
  )
}

const ContextSwitch = (props) => {
  const { setFieldValue, values } = useFormikContext()

  const getClassName = useCallback(
    (contextName) => {
      return values.context === contextName ? "active" : "inactive"
    },
    [values.context]
  )

  const onClick = useCallback(
    (event) => {
      setFieldValue("context", event.target.id, true)
    },
    [setFieldValue]
  )

  return (
    <div className="tabs flex">
      <div
        id="granted"
        className={`tab text-label ${getClassName("granted")}`}
        onClick={onClick}
      >
        granted
      </div>
      <div
        id="owned"
        className={`tab text-label ${getClassName("owned")}`}
        onClick={onClick}
      >
        owned
      </div>
    </div>
  )
}

const connectedWithFormik = withFormik({
  mapPropsToValues: () => ({
    selectedGrant: { id: null, amount: "0" },
    beneficiaryAddress: "",
    stakeTokens: "0",
    operatorAddress: "",
    authorizerAddress: "",
    context: "granted",
  }),
  validate: (values, props) => {
    const {
      beneficiaryAddress,
      operatorAddress,
      authorizerAddress,
      context,
      selectedGrant,
    } = values
    const errors = {}

    errors.stakeTokens = getStakeTokensError(props, values)
    errors.selectedGrant =
      context === "granted" && validateRequiredValue(selectedGrant.id)
    errors.beneficiaryAddress = validateEthAddress(beneficiaryAddress)
    errors.operatorAddress = validateEthAddress(operatorAddress)
    errors.authorizerAddress = validateEthAddress(authorizerAddress)

    return getErrorsObj(errors)
  },
  displayName: "DelegateStakeForm",
})(DelegateStakeForm)

const getStakeTokensError = (
  props,
  { stakeTokens, selectedGrant, context }
) => {
  const { keepBalance, minStake } = props
  const contextBalance =
    context === "granted" ? selectedGrant.availableToStake : keepBalance

  if (context === "granted") {
    if (!selectedGrant.id) {
      return "Please select grant"
    } else if (lte(selectedGrant.availableToStake, 0)) {
      return "Insufficient funds"
    } else {
      return validateAmountInRange(stakeTokens, contextBalance, minStake)
    }
  } else {
    if (lte(keepBalance, 0)) {
      return "Insufficient funds"
    } else {
      return validateAmountInRange(stakeTokens, contextBalance, minStake)
    }
  }
}

export default connectedWithFormik
