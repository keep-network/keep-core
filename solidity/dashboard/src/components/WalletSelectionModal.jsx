import React, { useEffect } from "react"
import { useModal } from "../hooks/useModal"
import WalletOptions from "./WalletOptions"
import Tooltip from "./Tooltip"

export const WalletSelectionModalContent = ({ closeModal }) => {
  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        <h3 className="ml-1">
          Select the wallet you want to connect with to proceed
        </h3>
      </div>
      <span className="text-center mt-1">
        {
          <Tooltip
            direction="top"
            simple
            className="empty-state__wallet-options-tooltip"
            triggerComponent={() => (
              <span
                className={`btn btn-primary btn-lg empty-state__connect-wallet-btn`}
              >
                Connect wallet
              </span>
            )}
          >
            <WalletOptions />
          </Tooltip>
        }
      </span>
    </div>
  )
}

const WalletSelectionModal = () => {
  const { openModal, closeModal } = useModal()

  const component = <WalletSelectionModalContent closeModal={closeModal} />

  useEffect(() => {
    openModal(component, { title: "Select wallet" })
  }, [])

  return null
}

export default WalletSelectionModal
