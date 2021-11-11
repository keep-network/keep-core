import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"

export const ConfirmUndelegation = withBaseModal(({ onConfirm, onClose }) => {
  const formik = useTypeTextToConfirmFormik("UNDELEGATE", onConfirm)

  return (
    <>
      <ModalHeader>Are you sure?</ModalHeader>
      <ModalBody>
        <h3>You’re about to undelegate.</h3>
        <p className="text-grey-60 mt-1">
          Undelegating will return all of your tokens to their owner. There is
          an undelegation period of 2 months until the tokens will be completely
          undelegated.
        </p>
        <form onSubmit={formik.handleSubmit} className="mt-2">
          <FormInputBase
            name="confirmationText"
            type="text"
            onChange={formik.handleChange}
            value={formik.values.confirmationText}
            label={"Type UNDELEGATE to confirm."}
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
          undelegate
        </Button>
        <Button className="btn btn-unstyled" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
})
