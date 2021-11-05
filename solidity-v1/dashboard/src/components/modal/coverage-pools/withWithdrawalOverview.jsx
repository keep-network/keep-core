import React from "react"
import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  ModalCloseButton,
} from "../Modal"
import { WithdrawalOverview } from "./components"

export const withWithdrawalOverview =
  ({ title }) =>
  (WrappedModalContent) => {
    return (props) => {
      return (
        <Modal isOpen onClose={props.onClose} size="xl">
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>{title}</ModalHeader>
            <ModalCloseButton />
            <div className="modal-with-withdraw-overview__content-wrapper">
              <WithdrawalOverview
                existingWithdrawalCovAmount={props.existingWithdrawalCovAmount}
                covAmountToAdd={props.covAmountToAdd}
                withdrawalDelay={props.withdrawalDelay}
                withdrawalInitiatedTimestamp={
                  props.withdrawalInitiatedTimestamp
                }
                withdrawalStatus={props.withdrawalStatus}
              />
              <WrappedModalContent {...props} />
            </div>
          </ModalContent>
        </Modal>
      )
    }
  }
