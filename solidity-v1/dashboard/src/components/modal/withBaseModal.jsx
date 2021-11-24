import React from "react"
import { Modal, ModalContent, ModalOverlay, ModalCloseButton } from "."

export function withBaseModal(WrappedModalContent) {
  return (props) => {
    return (
      <Modal isOpen={props.isOpen} onClose={props.onClose} size={props.size}>
        <ModalOverlay />
        <ModalContent>
          <ModalCloseButton />
          <WrappedModalContent {...props} />
        </ModalContent>
      </Modal>
    )
  }
}
