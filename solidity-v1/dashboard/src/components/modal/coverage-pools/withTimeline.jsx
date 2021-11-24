import React from "react"
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  ModalCloseButton,
} from "../Modal"
import { CovPoolTimeline } from "./components"

export const withTimeline =
  ({ title, step, withDescription }) =>
  (WrappedModalContent) => {
    return (props) => {
      return (
        <Modal isOpen onClose={props.onClose} size="xl">
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>{title}</ModalHeader>
            <ModalCloseButton />
            <div className="modal-with-timeline__content-wrapper">
              <CovPoolTimeline step={step} withDescription={withDescription} />
              <WrappedModalContent {...props} />
            </div>
          </ModalContent>
        </Modal>
      )
    }
  }
