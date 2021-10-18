import React, { useContext } from "react"
import * as Icons from "../Icons"

const ModalContext = React.createContext({
  isOpen: false,
  onClose: () => {},
  isCentered: true,
  closeOnOverlayClick: true,
})

const useModalContext = () => {
  const context = useContext(ModalContext)

  if (!context) {
    throw new Error("ModalContext used outside of Modal component")
  }

  return context
}

export const Modal = ({
  isOpen,
  onClose,
  isCentered,
  closeOnOverlayClick,
  children,
}) => {
  return (
    <ModalContext.Provider
      value={{ isOpen, onClose, isCentered, closeOnOverlayClick }}
    >
      <div className="modal">{children}</div>
    </ModalContext.Provider>
  )
}

export const ModalOverlay = ({
  className = "",
  color = "",
  onClick = () => {},
}) => {
  const { onClose, closeOnOverlayClick } = useModalContext()

  const _onClose = (event) => {
    event.stopPropagation()
    if (closeOnOverlayClick) {
      onClick()
    }
    onClose()
  }

  return (
    <div
      className={`modal__overlay ${
        color && `modal__overlay--${color}`
      } ${className}`}
      onClick={_onClose}
    />
  )
}

export const ModalContent = ({ className = "", ...restProps }) => {
  return (
    <div className="modal__content-wrapper">
      <section className="modal__content" {...restProps} />
    </div>
  )
}

export const ModalHeader = ({ className = "", ...restProps }) => {
  return <header className={`modal__header ${className}`} {...restProps} />
}

export const ModalBody = ({ className = "", ...restProps }) => {
  return <div className={`modal__body ${className}`} {...restProps} />
}

export const ModalFooter = ({ className = "", ...restProps }) => {
  return <footer className={`modal__footer ${className}`} {...restProps} />
}

export const ModalCloseButton = ({
  onClick = () => {},
  className = "",
  isDisabled,
  children,
  ...restProps
}) => {
  const { onClose } = useModalContext()

  const _onClose = (event) => {
    event.stopPropagation()
    onClick()
    onClose()
  }

  return (
    <button
      className={`modal__close-btn ${
        isDisabled ? "modal__close-btn--disabled" : ""
      } ${className}`}
      type="button"
      disabled={isDisabled}
      onClick={_onClose}
      {...restProps}
    >
      {children || <Icons.Cross width={15} height={15} />}
    </button>
  )
}
