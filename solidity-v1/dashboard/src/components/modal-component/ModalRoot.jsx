import React from "react"
import ReactDOM from "react-dom"
import { ExampleModal } from "."
import MobileUsersModal from "./MobileUsersModal"
import { useModal } from "../../hooks/useModal"
import { MODAL_TYPES } from "../../constants/constants"

const MODAL_TYPE_TO_COMPONENT = {
  [MODAL_TYPES.Example]: ExampleModal,
  [MODAL_TYPES.MobileUsers]: MobileUsersModal,
}
const modalRoot = document.getElementById("modal-root")

export const ModalRoot = () => {
  const { modalType, modalProps, closeModal } = useModal()

  if (!modalType) {
    return <></>
  }
  const SpecificModal = MODAL_TYPE_TO_COMPONENT[modalType]
  return ReactDOM.createPortal(
    <SpecificModal onClose={closeModal} {...modalProps} />,
    modalRoot
  )
}
