import React from "react"
import Button from "../Button"
import {
  Modal,
  ModalOverlay,
  ModalHeader,
  ModalCloseButton,
  ModalContent,
  ModalBody,
  ModalFooter,
} from "./Modal"

export const ExampleModal = ({ name, onClose, onConfirm }) => {
  return (
    <Modal isOpen onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Example Modal {name}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>Content of the modaaaalll!!</ModalBody>
        <ModalFooter>
          Footer here
          <Button onClick={() => onConfirm({ value: "hehee" })}>confirm</Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}
