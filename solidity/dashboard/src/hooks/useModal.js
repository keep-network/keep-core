import React, { useState } from "react"
import Modal from "../components/Modal"

export const useModal = () => {
  const [isVisible, setIsVisible] = useState(false)

  const showModal = () => setIsVisible(true)
  const hideModal = () => setIsVisible(false)

  const ModalComponent = (props) =>
    isVisible && <Modal closeModal={hideModal} {...props} />

  return { showModal, hideModal, ModalComponent }
}
