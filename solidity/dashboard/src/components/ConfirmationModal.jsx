import React from "react"
import Button from "./Button"
import FormInput from "./FormInput"
import { colors } from "../constants/colors"
import { withFormik } from "formik"
import { getErrorsObj } from "../forms/common-validators"

const ConfirmationModal = ({
  title,
  subtitle,
  confirmationText,
  btnText,
  onBtnClick,
  onCancel,
}) => {
  return (
    <>
      <h4 className="mb-1">{title}</h4>
      <div className="text-big text-grey-60 mb-3">{subtitle}</div>
      <ConfirmationFormFormik
        confirmationText={confirmationText}
        btnText={btnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
      />
    </>
  )
}

export default React.memo(ConfirmationModal)

export const withConfirmationModal = (WrappedComponent) => ({
  title,
  subtitle,
  confirmationText,
  btnText,
  onBtnClick,
  onCancel,
  ...restProps
}) => {
  return (
    <ConfirmationModal
      title={title}
      btnText={btnText}
      confirmationText={confirmationText}
      onCancel={onCancel}
      onBtnClick={onBtnClick}
      subtitle={<WrappedComponent {...restProps} />}
    />
  )
}

const ConfirmationForm = ({
  confirmationText,
  btnText,
  onCancel,
  ...formikProps
}) => {
  return (
    <form>
      <FormInput
        name="confirmationText"
        type="text"
        label={`Type ${confirmationText} to confirm.`}
        placeholder=""
      />
      <div
        className="flex row center mt-2"
        style={{
          borderTop: `1px solid ${colors.grey20}`,
          margin: "0 -2rem",
          padding: "2rem 2rem 0",
        }}
      >
        <Button
          className="btn btn-primary"
          type="submit"
          disabled={!(formikProps.isValid && formikProps.dirty)}
          onClick={formikProps.handleSubmit}
        >
          {btnText}
        </Button>
        <span onClick={onCancel} className="ml-1 text-link">
          Cancel
        </span>
      </div>
    </form>
  )
}

const ConfirmationFormFormik = withFormik({
  mapPropsToValues: () => ({
    confirmationText: "",
  }),
  validate: (values, { confirmationText }) => {
    const errors = {}

    if (values.confirmationText !== confirmationText) {
      errors.confirmationText = "Does not match"
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props }) => props.onBtnClick(values),
  displayName: "ConfirmationForm",
})(ConfirmationForm)
