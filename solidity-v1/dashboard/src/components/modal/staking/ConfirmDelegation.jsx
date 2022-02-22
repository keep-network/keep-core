import React from "react"
import moment from "moment"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"
import OnlyIf from "../../OnlyIf"

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
            <p className="text-black mt-1 text-bold">
              To use this stake in Threshold, you will first have to contact
              your grant manager. Please do it immediately after you stake.
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
            cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
