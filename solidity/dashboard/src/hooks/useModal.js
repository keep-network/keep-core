import React, { useState, useCallback } from "react"
import Modal from "../components/Modal"

export const useModal = () => {
  const [isOpen, setIsOpen] = useState(false)

  const openModal = useCallback(() => {
    setIsOpen(true)
  }, [])
  const closeModal = useCallback(() => {
    setIsOpen(false)
  }, [])

  const ModalComponent = useCallback(
    (props) => {
      return <Modal isOpen={isOpen} closeModal={closeModal} {...props} />
    },
    [closeModal, isOpen]
  )

  return { openModal, closeModal, ModalComponent }
}
