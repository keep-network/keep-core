import { useCallback, useEffect } from "react"
import { useLocalStorage } from "./useLocalStorage"
import { useModal } from "./useModal"
import { MODAL_TYPES } from "../constants/constants"

const KEY = "shouldShowLegacyDappModal"

export const useShouldShowLegacyDappModal = () => {
  const [shouldShowModal, setShouldShowModal] = useLocalStorage(KEY, true)

  const modalHasBeenClosed = useCallback(() => {
    setShouldShowModal(false)
  }, [setShouldShowModal])

  return { shouldShowModal, modalHasBeenClosed }
}

export const useShowLegacyDappModal = () => {
  const { shouldShowModal } = useShouldShowLegacyDappModal()
  const { openModal } = useModal()

  useEffect(() => {
    if (!shouldShowModal) return

    openModal(MODAL_TYPES.LegacyDashboard)
  }, [shouldShowModal, openModal])
}
