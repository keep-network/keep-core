import React from "react"
import { SubmitButton } from "./Button"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import FormInput from "./FormInput"
import { colors } from "../constants/colors"

const AvailableETHForm = ({ onSubmit, closeModal, action, ...formikProps }) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)

  return (
    <form>
      <FormInput
        name="ethToWithdraw"
        type="text"
        label="ETH Amount"
        placeholder="0"
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
          className="btn btn-primary"
          type="submit"
          onSubmitAction={onSubmitBtn}
          withMessageActionIsPending={false}
          triggerManuallyFetch={true}
          disabled={!formikProps.dirty}
        >
          {action}
        </SubmitButton>
        <span onClick={closeModal} className="ml-1 text-link">
          Cancel
        </span>
      </div>
    </form>
  )
}

export default React.memo(AvailableETHForm)
