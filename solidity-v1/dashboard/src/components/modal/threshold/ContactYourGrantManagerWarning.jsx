import React from "react"
import { withBaseModal } from "../withBaseModal"
import { ModalBody, ModalFooter, ModalHeader } from "../Modal"
import Button from "../../Button"

export const ContactYourGrantManagerWarning = withBaseModal(
  ({
    header = "Contact your grant manager",
    bodyTitle = "To enable staking into Threshold please contact your grant manager first.",
    onClose,
  }) => {
    return (
      <>
        <ModalHeader>{header}</ModalHeader>
        <ModalBody>
          <h3>{bodyTitle}</h3>
        </ModalBody>
        <ModalFooter>
          <Button className={"btn btn-secondary btn-lg"} onClick={onClose}>
            close
          </Button>
        </ModalFooter>
      </>
    )
  }
)
