import React from "react"
import Button from "../../../Button"
import { ModalBody, ModalHeader, ModalFooter } from "../../Modal"
import { withBaseModal } from "../../withBaseModal"

export const WarningBeforeDeposit = withBaseModal(({ onConfirm, onClose }) => {
  return (
    <>
      <ModalHeader>Deposit</ModalHeader>
      <ModalBody>
        <h3>Take note!</h3>
        The coverage pool is about to cover an event. Do you want to continue
        with this deposit?
      </ModalBody>
      <ModalFooter>
        <ModalFooter>
          <Button
            className="btn btn-primary btn-lg mr-2"
            type="submit"
            onClick={onConfirm}
          >
            continue
          </Button>
          <Button className="btn btn-unstyled" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </ModalFooter>
    </>
  )
})
