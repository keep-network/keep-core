import { useEffect, useState } from "react"
import useCurrentWidth from "./useCurrentWidth"
import { useModal } from "./useModal"
import { MODAL_TYPES } from "../constants/constants"

const MODAL_WINDOW_STATUS = {
  NOT_DISPLAYED: "NOT_DISPLAYED",
  IS_DISPLAYING: "IS_DISPLAYING",
  DISPLAYED: "DISPLAYED",
}

const useModalWindowForMobileUsers = () => {
  const widthThreshold = 1000
  const currentWidth = useCurrentWidth()
  const [modalWindowStatus, setModalWindowStatus] = useState(
    MODAL_WINDOW_STATUS.NOT_DISPLAYED
  )
  const { openModal, closeModal } = useModal()

  useEffect(() => {
    const customModalWindowForMobileUsersClose = () => {
      setModalWindowStatus(MODAL_WINDOW_STATUS.DISPLAYED)
      closeModal()
    }

    if (
      currentWidth < widthThreshold &&
      modalWindowStatus === MODAL_WINDOW_STATUS.NOT_DISPLAYED
    ) {
      openModal(MODAL_TYPES.MobileUsers, {
        onClose: customModalWindowForMobileUsersClose,
      })
      setModalWindowStatus(MODAL_WINDOW_STATUS.IS_DISPLAYING)
    }

    if (
      currentWidth > widthThreshold &&
      modalWindowStatus === MODAL_WINDOW_STATUS.IS_DISPLAYING
    ) {
      customModalWindowForMobileUsersClose()
    }
  }, [currentWidth, openModal, closeModal, modalWindowStatus])
}

export default useModalWindowForMobileUsers
