import React, { useEffect, useRef, useCallback, useState } from "react"
import ReactDOM from "react-dom"
import * as Icons from "./Icons"
import ConfirmationModal from "./ConfirmationModal"

const modalRoot = document.getElementById("modal-root")
const crossIconHeight = 15
const crossIconWidth = 15

const Modal = React.memo(({ isOpen, closeModal, ...props }) => {
  const modalOverlay = useRef(null)
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden"
    }

    return () => {
      document.body.style.overflow = "scroll"
    }
  }, [isOpen])

  const onOverlayClick = useCallback(
    (event) => {
      if (modalOverlay.current === event.target) {
        closeModal()
      }
    },
    [closeModal]
  )

  return isOpen
    ? ReactDOM.createPortal(
        <div
          ref={modalOverlay}
          className="modal-overlay"
          onClick={onOverlayClick}
        >
          <div className="modal-wrapper">
            <div className="modal-title">
              <h4 className="text-darker-grey">{props.title}</h4>
              <div className="modal-close" onClick={closeModal}>
                <Icons.Cross width={crossIconWidth} height={crossIconHeight} />
              </div>
            </div>
            <div className="modal-content">{props.children}</div>
          </div>
        </div>,
        modalRoot
      )
    : null
})

export default Modal

export const ModalContext = React.createContext({
  openModal: (component) => {},
  closeModal: () => {},
  showConfirmationModal: (modalComponent) => {},
})

export const ModalContextProvider = ({ children }) => {
  const [modalComponent, setModalComponent] = useState()
  const [isOpen, setIsOpen] = useState(false)
  const [modalOptions, setModalOptions] = useState(null)
  const awaitingPromiseRef = useRef()

  const openModal = useCallback((modalComponent, modalOptions = {}) => {
    setModalComponent(modalComponent)
    setModalOptions(modalOptions)
    setIsOpen(true)
  }, [])

  const closeModal = useCallback(() => {
    if (awaitingPromiseRef.current) {
      awaitingPromiseRef.current.reject()
    }
    setModalOptions(null)
    setIsOpen(false)
  }, [])

  const onSubmitConfirmationModal = useCallback(
    (values) => {
      if (awaitingPromiseRef.current) {
        awaitingPromiseRef.current.resolve(values)
      }
      closeModal()
    },
    [closeModal]
  )

  const openConfirmationModal = useCallback(
    (options, ConfirmationModalComponent = ConfirmationModal) => {
      const confirmationModal = (
        <ConfirmationModalComponent
          onCancel={closeModal}
          onBtnClick={onSubmitConfirmationModal}
          {...options}
        />
      )
      openModal(confirmationModal)

      return new Promise((resolve, reject) => {
        awaitingPromiseRef.current = { resolve, reject }
      })
    },
    [openModal, onSubmitConfirmationModal, closeModal]
  )

  return (
    <ModalContext.Provider
      value={{
        openConfirmationModal,
        openModal,
        closeModal,
      }}
    >
      <Modal isOpen={isOpen} closeModal={closeModal} {...modalOptions}>
        {modalComponent}
      </Modal>
      {children}
    </ModalContext.Provider>
  )
}
