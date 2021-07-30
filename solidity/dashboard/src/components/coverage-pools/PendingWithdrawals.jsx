import React from "react"
import {
  claimTokensFromWithdrawal,
  withdrawAssetPool,
} from "../../actions/coverage-pool"
import PendingWithdrawalsView from "./PendingWithdrawalsView"
import { useDispatch, useSelector } from "react-redux"
import ClaimTokensModal from "./ClaimTokensModal"
import { useModal } from "../../hooks/useModal"

const PendingWithdrawals = () => {
  const dispatch = useDispatch()
  const { openConfirmationModal, closeModal } = useModal()
  const {
    withdrawalDelay,
    withdrawalTimeout,
    pendingWithdrawals,
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

  const onReinitiateWithdrawal = async (awaitingPromise) => {
    dispatch(withdrawAssetPool("0", awaitingPromise))
  }

  return (
    <PendingWithdrawalsView
      onClaimTokensSubmitButtonClick={onClaimTokensSubmitButtonClick}
      onReinitiateWithdrawal={onReinitiateWithdrawal}
      withdrawalDelay={withdrawalDelay}
      withdrawalTimeout={withdrawalTimeout}
      pendingWithdrawals={pendingWithdrawals}
    />
  )
}

export default PendingWithdrawals
