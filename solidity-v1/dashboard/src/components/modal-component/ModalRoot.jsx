import React from "react"
import ReactDOM from "react-dom"
import { ExampleModal } from "."
import { useModal } from "../../hooks/useModal"
import { MODAL_TYPES } from "../../constants/constants"

const MODAL_TYPE_TO_COMPONENT = {
  [MODAL_TYPES.Example]: ExampleModal,
}
const modalRoot = document.getElementById("modal-root")

export const ModalRoot = () => {
  const { modalType, modalProps, closeModal } = useModal()
  console.log("modalType", modalType, modalProps, closeModal)

  if (!modalType) {
    return <></>
  }
  const SpecificModal = MODAL_TYPE_TO_COMPONENT[modalType]
  return ReactDOM.createPortal(
    <SpecificModal closeModal={closeModal} {...modalProps} />,
    modalRoot
  )
}
