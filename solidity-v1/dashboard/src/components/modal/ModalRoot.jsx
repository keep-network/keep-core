import React from "react"
import ReactDOM from "react-dom"
import MobileUsersModal from "./MobileUsersModal"
import { useModal } from "../../hooks/useModal"
import { MODAL_TYPES } from "../../constants/constants"
import { WithdrawETHModal, AddETHModal } from "./bonding"
import {
  MetaMaskModal,
  TallyModal,
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
  ClaimTokens as ClaimStakingTokens,
  ConfirmRecovering,
  TokensClaimed,
  UndelegateStake,
  UndelegationInitiated,
  ConfirmCancelDelegationFromGrant,
  CopyStake,
  ConfirmReleaseTokensFromGrant,
} from "./staking"
import { AddKeep, WithdrawKeep } from "./liquidity"
import {
  InitiateWithdraw,
  WithdrawInitialized,
  ClaimTokens,
  ReInitiateWithdraw,
  ConfirmIncreaseWithdrawal,
  IncreaseWithdrawal,
} from "./coverage-pools"
import { WithdrawGrantedTokens } from "./threshold/WithdrawGrantedTokens"
import { GrantTokensWithdrawn } from "./threshold/GrantTokensWithdrawn"
import {
  AuthorizeAndStakeOnThreshold,
  StakeOnThresholdConfirmed,
  StakeOnThresholdWithoutAuthorization,
} from "./threshold/StakeOnThreshold"
import {
  ThresholdAuthorizationLoadingModal,
  ThresholdStakeConfirmationLoadingModal,
} from "./threshold/ThresholdLoadingModal"
import { AuthorizedButNotStakedToTWarning } from "./threshold/AuthorizedButNotStakedToTWarning"
import { ContactYourGrantManagerWarning } from "./threshold/ContactYourGrantManagerWarning"
import { LegacyDashboardModal } from "./LegacyDashboardModal"

const MODAL_TYPE_TO_COMPONENT = {
  [MODAL_TYPES.MobileUsers]: MobileUsersModal,
  [MODAL_TYPES.BondingAddETH]: AddETHModal,
  [MODAL_TYPES.BondingWithdrawETH]: WithdrawETHModal,
  [MODAL_TYPES.MetaMask]: MetaMaskModal,
  [MODAL_TYPES.Tally]: TallyModal,
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
  [MODAL_TYPES.ClaimStakingTokens]: ClaimStakingTokens,
  [MODAL_TYPES.StakingTokensClaimed]: TokensClaimed,
  [MODAL_TYPES.GrantTokensWithdrawn]: GrantTokensWithdrawn,
  [MODAL_TYPES.ConfirmCancelDelegationFromGrant]:
    ConfirmCancelDelegationFromGrant,
  [MODAL_TYPES.UndelegateStake]: UndelegateStake,
  [MODAL_TYPES.UndelegationInitiated]: UndelegationInitiated,
  [MODAL_TYPES.CopyStake]: CopyStake,
  [MODAL_TYPES.ConfirmReleaseTokensFromGrant]: ConfirmReleaseTokensFromGrant,
  [MODAL_TYPES.InitiateCovPoolWithdraw]: InitiateWithdraw,
  [MODAL_TYPES.CovPoolWithdrawInitialized]: WithdrawInitialized,
  [MODAL_TYPES.CovPoolClaimTokens]: ClaimTokens,
  [MODAL_TYPES.ReInitiateCovPoolWithdraw]: ReInitiateWithdraw,
  [MODAL_TYPES.ConfirmCovPoolIncreaseWithdrawal]: ConfirmIncreaseWithdrawal,
  [MODAL_TYPES.IncreaseCovPoolWithdrawal]: IncreaseWithdrawal,
  [MODAL_TYPES.WithdrawGrantedTokens]: WithdrawGrantedTokens,
  [MODAL_TYPES.AuthorizeAndStakeOnThreshold]: AuthorizeAndStakeOnThreshold,
  [MODAL_TYPES.StakeOnThresholdWithoutAuthorization]:
    StakeOnThresholdWithoutAuthorization,
  [MODAL_TYPES.StakeOnThresholdConfirmed]: StakeOnThresholdConfirmed,
  [MODAL_TYPES.ThresholdAuthorizationLoadingModal]:
    ThresholdAuthorizationLoadingModal,
  [MODAL_TYPES.ThresholdStakeConfirmationLoadingModal]:
    ThresholdStakeConfirmationLoadingModal,
  [MODAL_TYPES.AuthorizedButNotStakedToTWarningModal]:
    AuthorizedButNotStakedToTWarning,
  [MODAL_TYPES.ContactYourGrantManagerWarning]: ContactYourGrantManagerWarning,
  [MODAL_TYPES.LegacyDashboard]: LegacyDashboardModal,
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
      isOpen={!!modalType}
      {...restProps}
    />,
    modalRoot
  )
}
