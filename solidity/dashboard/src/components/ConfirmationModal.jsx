import React from "react"
import Button from "./Button"
import FormInput from "./FormInput"
import { colors } from "../constants/colors"
import { withFormik } from "formik"
import { getErrorsObj } from "../forms/common-validators"
import Divider from "./Divider"
import FormCheckbox from "./FormCheckbox"

const ConfirmationModal = ({
  title,
  subtitle,
  confirmationText,
  btnText,
  onBtnClick,
  onCancel,
  getLabelText,
}) => {
  return (
    <>
      <h3 className="mb-1">{title}</h3>
      <div className="text-big text-grey-60 mb-3">{subtitle}</div>
      <ConfirmationFormFormik
        confirmationText={confirmationText}
        btnText={btnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
        getLabelText={getLabelText}
      />
    </>
  )
}

const withConfirmationModal = (WrappedComponent) => ({
  title,
  subtitle,
  confirmationText,
  btnText,
  onBtnClick,
  onCancel,
  getLabelText,
  ...restProps
}) => {
  return (
    <ConfirmationModal
      title={title}
      btnText={btnText}
      confirmationText={confirmationText}
      onCancel={onCancel}
      onBtnClick={onBtnClick}
      getLabelText={getLabelText}
      subtitle={<WrappedComponent {...restProps} />}
    />
  )
}

const ConfirmationForm = ({
  confirmationText,
  btnText,
  onCancel,
  getLabelText = (confirmationText) => `Type ${confirmationText} to confirm.`,
  ...formikProps
}) => {
  return (
    <form>
      <FormInput
        name="confirmationText"
        type="text"
        label={getLabelText(confirmationText)}
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
          className="btn btn-lg btn-primary"
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

const AcceptTermForm = ({
  termText,
  btnText,
  onSubmit,
  onCancel,
  ...formikProps
}) => {
  return (
    <form>
      <Divider className="divider divider--tile-fluid" />
      <FormCheckbox label={termText} name="checked" type="checkbox" />
      <div className="flex row center mt-2">
        <Button
          className="btn btn-lg btn-primary"
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

const AcceptTermFormik = withFormik({
  mapPropsToValues: () => ({
    checked: false,
  }),
  validate: (values, { confirmationText }) => {
    const errors = {}

    if (!values.checked) {
      errors.checked = "Required"
    }

    return getErrorsObj(errors)
  },
  handleSubmit: (values, { props }) => props.onBtnClick(values),
  displayName: "AcceptTermForm",
})(AcceptTermForm)

const AcceptTermConfirmationModal = ({
  title,
  termText,
  btnText,
  onBtnClick,
  onCancel,
  children,
}) => {
  return (
    <section>
      <h3 className="mb-1">{title}</h3>
      {children}
      <AcceptTermFormik
        termText={termText}
        btnText={btnText}
        onBtnClick={onBtnClick}
        onCancel={onCancel}
      />
    </section>
  )
}

export default React.memo(ConfirmationModal)

export { withConfirmationModal, AcceptTermConfirmationModal }
