import React from "react"
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  ModalCloseButton,
} from "./Modal"
import { CovPoolTimeline } from "./coverage-pools/components"

export const withTimeline =
  ({
    title,
    timelineComponent: TimelineComponent = CovPoolTimeline,
    timelineProps = {},
  }) =>
  (WrappedModalContent) => {
    return (props) => {
      return (
        <Modal isOpen onClose={props.onClose} size="xl">
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>{title}</ModalHeader>
            <ModalCloseButton />
            <div className="modal-with-timeline__content-wrapper">
              <TimelineComponent {...timelineProps} />
              <WrappedModalContent {...props} />
            </div>
          </ModalContent>
        </Modal>
      )
    }
  }
