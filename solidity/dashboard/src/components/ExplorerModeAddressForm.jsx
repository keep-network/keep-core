import FormInput from "./FormInput"
import { SubmitButton } from "./Button"
import React from "react"
import { withFormik } from "formik"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import { getErrorsObj, validateEthAddress } from "../forms/common-validators"
import { colors } from "../constants/colors"

const ExplorerModeAddressForm = ({ submitAction, onCancel }) => {
  const onSubmit = useCustomOnSubmitFormik(submitAction)

  return (
    <form>
      <FormInput
        name="address"
        type="text"
        label="Enter an Ethereum address"
        tooltipProps={{ direction: "top" }}
        tooltipText={
          "Enter an Ethereum address to preview its read-only version of the dashboard."
        }
      />
      <div
        className="flex row center mt-2"
        style={{
          borderTop: `1px solid ${colors.grey20}`,
          margin: "0 -2rem",
          padding: "2rem 2rem 0",
        }}
      >
        <SubmitButton
          className="btn btn-lg btn-primary"
          type="submit"
          onSubmitAction={onSubmit}
          withMessageActionIsPending={false}
          triggerManuallyFetch={false}
        >
          explore
        </SubmitButton>
        <span onClick={onCancel} className="ml-1 text-link">
          Cancel
        </span>
      </div>
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
