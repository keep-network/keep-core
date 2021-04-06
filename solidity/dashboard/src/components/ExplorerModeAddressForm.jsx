import FormInput from "./FormInput"
import { SubmitButton } from "./Button"
import React from "react"
import { withFormik } from "formik"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import { getErrorsObj, validateEthAddress } from "../forms/common-validators"

const ExplorerModeAddressForm = ({ submitAction }) => {
  const onSubmit = useCustomOnSubmitFormik(submitAction)

  return (
    <form>
      <FormInput
        name="address"
        type="text"
        label="Enter an address"
        tooltipText={<>tooltip text</>}
      />
      <SubmitButton
        className="btn btn-primary"
        type="submit"
        onSubmitAction={onSubmit}
        withMessageActionIsPending={false}
        triggerManuallyFetch={false}
      >
        explore
      </SubmitButton>
    </form>
  )
}

const connectWithFormik = withFormik({
  mapPropsToValues: () => ({
    address: "",
  }),
  validate: (values, props) => {
    const { address } = values
    const errors = {}

    errors.address = validateEthAddress(address)

    return getErrorsObj(errors)
  },
})(ExplorerModeAddressForm)

export default connectWithFormik
