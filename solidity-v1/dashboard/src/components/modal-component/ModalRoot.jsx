import React from "react"
import ReactDOM from "react-dom"
import MobileUsersModal from "./MobileUsersModal"
import { useModal } from "../../hooks/useModal"
import { MODAL_TYPES } from "../../constants/constants"
import { WithdrawETHModal, AddETHModal } from "./bonding"
import {
  MetaMaskModal,
  LedgerModal,
  TrezorModal,
  ExplorerModeModal,
  WalletConnectModal,
  WalletSelectionModal,
} from "./wallets"
import {
  DelegationAlreadyExists,
  TopUpInitiatedConfirmation,
  TopUpInitialization,
  ConfirmTopUpInitialization,
  ConfirmDelegation,
  ConfirmRecovering,
  ConfirmUndelegation,
  ConfirmCancelDelegationFromGrant,
  CopyStake,
} from "./staking"
import { AddKeep, WithdrawKeep } from "./liquidity"

const MODAL_TYPE_TO_COMPONENT = {
  [MODAL_TYPES.MobileUsers]: MobileUsersModal,
  [MODAL_TYPES.BondingAddETH]: AddETHModal,
  [MODAL_TYPES.BondingWithdrawETH]: WithdrawETHModal,
  [MODAL_TYPES.MetaMask]: MetaMaskModal,
  [MODAL_TYPES.Ledger]: LedgerModal,
  [MODAL_TYPES.Trezor]: TrezorModal,
  [MODAL_TYPES.ExplorerMode]: ExplorerModeModal,
  [MODAL_TYPES.WalletConnect]: WalletConnectModal,
  [MODAL_TYPES.WalletSelection]: WalletSelectionModal,
  [MODAL_TYPES.DelegationAlreadyExists]: DelegationAlreadyExists,
  [MODAL_TYPES.TopUpInitialization]: TopUpInitialization,
  [MODAL_TYPES.ConfirmTopUpInitialization]: ConfirmTopUpInitialization,
  [MODAL_TYPES.TopUpInitiatedConfirmation]: TopUpInitiatedConfirmation,
  [MODAL_TYPES.KeepOnlyPoolAddKeep]: AddKeep,
  [MODAL_TYPES.KeepOnlyPoolWithdrawKeep]: WithdrawKeep,
  [MODAL_TYPES.ConfirmDelegation]: ConfirmDelegation,
  [MODAL_TYPES.ConfirmRecovering]: ConfirmRecovering,
  [MODAL_TYPES.ConfirmCancelDelegationFromGrant]:
    ConfirmCancelDelegationFromGrant,
  [MODAL_TYPES.ConfirmUndelegation]: ConfirmUndelegation,
  [MODAL_TYPES.CopyStake]: CopyStake,
}

const modalRoot = document.getElementById("modal-root")

export const ModalRoot = () => {
  const {
    modalType,
    modalProps,
    closeModal,
    closeConfirmationModal,
    onSubmitConfirmationModal,
  } = useModal()
  const { onClose, onConfirm, ...restProps } = modalProps

  if (!modalType) {
    return <></>
  }

  const _onConfirm = (values) => {
    if (onConfirm && typeof onClose === "function") {
      onConfirm(values)
    } else {
      onSubmitConfirmationModal(values)
    }
    if (modalProps.shouldCloseOnSubmit) {
      // Just close modal we don't want to dispatch `modal/cancel`
      // action.
      closeModal()
    }
  }

  const SpecificModal = MODAL_TYPE_TO_COMPONENT[modalType]
  return ReactDOM.createPortal(
    <SpecificModal
      onClose={() => {
        onClose && typeof onClose === "function" && onClose()
        if (modalProps.isConfirmationModal) {
          closeConfirmationModal()
        } else {
          closeModal()
        }
      }}
      onConfirm={!modalProps.isConfirmationModal ? undefined : _onConfirm}
      {...restProps}
    />,
    modalRoot
  )
}
