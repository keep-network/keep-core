import React from "react"
import moment from "moment"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"
import OnlyIf from "../../OnlyIf"
import { GRANT_MANAGER_EMAIL } from "../../../constants/constants"
import { colors } from "../../../constants/colors"

const styles = {
  warning: {
    padding: "1rem",
    backgroundColor: colors.yellowPrimary,
    color: colors.grey70,
    border: `1px solid ${colors.yellowSecondary}`,
    borderRadius: "0.5rem",
  },
  grantManager: {
    color: colors.yellowSecondary,
  },
}

export const ConfirmDelegation = withBaseModal(
  ({ initializationPeriod, onConfirm, isFromGrant = false, onClose }) => {
    const formik = useTypeTextToConfirmFormik("DELEGATE", onConfirm)

    return (
      <>
        <ModalHeader>Are you sure?</ModalHeader>
        <ModalBody>
          <h3>You’re about to delegate stake.</h3>
          <p className="text-grey-60 mt-1">
            You’re delegating KEEP tokens. You will be able to cancel the
            delegation for up to{" "}
            {moment().add(initializationPeriod, "seconds").fromNow(true)}. After
            that time, you can undelegate your stake.
          </p>
          <OnlyIf condition={isFromGrant}>
            <p style={styles.warning}>
              Please do not forget to contact your Grant Manager{" "}
              <span className="text-yellow-100">(</span>
              <a
                style={styles.grantManager}
                href={`mailto:${GRANT_MANAGER_EMAIL}`}
              >
                {GRANT_MANAGER_EMAIL}
              </a>
              <span className="text-yellow-100">)</span>&nbsp;immediately!
            </p>
          </OnlyIf>
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              name="confirmationText"
              type="text"
              onChange={formik.handleChange}
              value={formik.values.confirmationText}
              label={"Type DELEGATE to confirm."}
              placeholder=""
              hasError={formik.errors.confirmationText}
              errorMsg={formik.errors.confirmationText}
            />
          </form>
        </ModalBody>
        <ModalFooter>
          <Button
            className="btn btn-primary btn-lg mr-2"
            type="submit"
            onClick={formik.handleSubmit}
          >
            delegate
          </Button>
          <Button className="btn btn-unstyled" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
