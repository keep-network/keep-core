import React from "react"
import {
  claimTokensFromWithdrawal,
  withdrawAssetPool,
} from "../../actions/coverage-pool"
import PendingWithdrawalsView from "./PendingWithdrawalsView"
import { useDispatch, useSelector } from "react-redux"
import ClaimTokensModal from "./ClaimTokensModal"
import { useModal } from "../../hooks/useModal"
import ReinitiateWithdrawalModal from "./ReinitiateWithdrawalModal"

const PendingWithdrawals = ({ covTokensAvailableToWithdraw }) => {
  const dispatch = useDispatch()
  const { openConfirmationModal, closeModal } = useModal()
  const {
    withdrawalDelay,
    withdrawalTimeout,
    pendingWithdrawal,
    withdrawalInitiatedTimestamp,
  } = useSelector((state) => state.coveragePool)

  const onClaimTokensSubmitButtonClick = async (covAmount, awaitingPromise) => {
    await openConfirmationModal(
      {
        closeModal: closeModal,
        submitBtnText: "claim",
        amount: covAmount,
        modalOptions: {
          title: "Claim tokens",
          classes: {
            modalWrapperClassName: "modal-wrapper__claim-tokens",
          },
        },
      },
      ClaimTokensModal
    )
    dispatch(claimTokensFromWithdrawal(awaitingPromise))
  }

  const onReinitiateWithdrawal = async (
    pendingWithdrawal,
    covTokensAvailableToWithdraw,
    awaitingPromise
  ) => {
    await openConfirmationModal(
      {
        modalOptions: {
          title: "Re-initiate withdrawal",
          classes: {
            modalWrapperClassName: "modal-wrapper__reinitiate-withdrawal",
          },
        },
        submitBtnText: "continue",
        pendingWithdrawalBalance: pendingWithdrawal,
        covTokensAvailableToWithdraw,
        containerTitle: "You are about to re-initiate this withdrawal:",
      },
      ReinitiateWithdrawalModal
    )
    dispatch(withdrawAssetPool("0", awaitingPromise))
  }

  return (
    <PendingWithdrawalsView
      onClaimTokensSubmitButtonClick={onClaimTokensSubmitButtonClick}
      onReinitiateWithdrawal={onReinitiateWithdrawal}
      withdrawalDelay={withdrawalDelay}
      withdrawalTimeout={withdrawalTimeout}
      pendingWithdrawal={pendingWithdrawal}
      withdrawalInitiatedTimestamp={withdrawalInitiatedTimestamp}
      covTokensAvailableToWithdraw={covTokensAvailableToWithdraw}
    />
  )
}

export default PendingWithdrawals
