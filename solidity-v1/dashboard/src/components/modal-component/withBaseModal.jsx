import React from "react"
import { Modal, ModalContent, ModalOverlay, ModalCloseButton } from "."

function withBaseModal(WrappedModalContent) {
  return (props) => {
    return (
      <Modal isOpen onClose={props.onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalCloseButton />
          <WrappedModalContent {...props} />
        </ModalContent>
      </Modal>
    )
  }
}

export default withBaseModal
