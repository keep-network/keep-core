import React, { useEffect, useRef, useCallback, useState } from "react"
import ReactDOM from "react-dom"
import * as Icons from "./Icons"
import ConfirmationModal from "./ConfirmationModal"
import { useDispatch, useSelector } from "react-redux"
import ClaimTokensModal from "./coverage-pools/ClaimTokensModal"
import InitiateCovPoolsWithdrawModal from "./coverage-pools/InitiateCovPoolsWithdrawModal"
import { InitiateDepositModal } from "./coverage-pools"
import IncreaseWithdrawalModal from "./coverage-pools/IncreaseWithdrawalModal"
import { clearAdditionalDataFromModal, hideModal } from "../actions/modal"
import { MigrationCompletedModal } from "./tbtc-migration"

const modalRoot = document.getElementById("modal-root")
const crossIconHeight = 15
const crossIconWidth = 15

export const modalComponentType = {
  COV_POOLS: {
    KEEP_DEPOSITED_SUCCESS: InitiateDepositModal,
    INITIATE_WITHDRAWAL: InitiateCovPoolsWithdrawModal,
    RE_INITIATE_WITHDRAWAL: InitiateCovPoolsWithdrawModal,
    INCREASE_WITHDRAWAL: IncreaseWithdrawalModal,
    WITHDRAWAL_COMPLETED: ClaimTokensModal,
  },
  TBTC_MIGRATION: {
    MIGRATION_COMPLETED: MigrationCompletedModal,
  },
}

const Modal = React.memo(
  ({ isOpen, closeModal, isFullScreen, hideTitleBar, classes, ...props }) => {
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
            <div
              className={`modal-wrapper${isFullScreen ? "--full-screen" : ""} ${
                classes?.modalWrapperClassName
                  ? classes.modalWrapperClassName
                  : ""
              }`}
            >
              {!isFullScreen && !hideTitleBar && (
                <div className="modal-title">
                  <h4 className="text-darker-grey">{props.title}</h4>
                  <div className="modal-close" onClick={closeModal}>
                    <Icons.Cross
                      width={crossIconWidth}
                      height={crossIconHeight}
                    />
                  </div>
                </div>
              )}
              <div className="modal-content">{props.children}</div>
            </div>
          </div>,
          modalRoot
        )
      : null
  }
)

export default Modal

export const ModalContext = React.createContext({
  openModal: (component) => {},
  closeModal: () => {},
  showConfirmationModal: (modalComponent) => {},
})

export const ModalContextProvider = ({ children }) => {
  const dispatch = useDispatch()
  const [modalComponent, setModalComponent] = useState()
  const [isOpen, setIsOpen] = useState(false)
  const [modalOptions, setModalOptions] = useState(null)
  const awaitingPromiseRef = useRef()
  const modal = useSelector((state) => state.modal)

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

  const closeEventModal = useCallback(() => {
    dispatch(clearAdditionalDataFromModal())
    dispatch(hideModal())
    closeModal()
  }, [dispatch, closeModal])

  const openSpecificModal = useCallback(() => {
    const SpecificComponent = modal.modalComponentType
    openModal(
      <SpecificComponent
        onCancel={closeEventModal}
        {...modal.componentProps}
      />,
      {
        closeModal: closeEventModal,
        ...modal.modalProps,
      }
    )
  }, [
    closeEventModal,
    modal.componentProps,
    modal.modalComponentType,
    modal.modalProps,
    openModal,
  ])

  useEffect(() => {
    if (modal.isOpen) {
      openSpecificModal()
    }
  }, [modal.isOpen, openSpecificModal])

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
      const { modalOptions, ...confirmationModalOptions } = options
      const confirmationModal = (
        <ConfirmationModalComponent
          onCancel={closeModal}
          onBtnClick={onSubmitConfirmationModal}
          {...confirmationModalOptions}
        />
      )
      openModal(confirmationModal, modalOptions)

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
