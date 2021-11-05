import React from "react"
import moment from "moment"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"

export const ConfirmDelegation = withBaseModal(
  ({ initializationPeriod, onConfirm, onClose }) => {
    const formik = useTypeTextToConfirmFormik("DELEGATE", onConfirm)

    return (
      <>
        <ModalHeader>Are you sure?</ModalHeader>
        <ModalBody>
          <h3>You’re about to delegate stake.</h3>
          <p className="text-grey-60 mt-1">
            You’re delegating KEEP tokens. You will be able to cancel the
            delegation for up to
            {moment().add(initializationPeriod, "seconds").fromNow(true)}. After
            that time, you can undelegate your stake.
          </p>
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
