import React, { useCallback } from "react"
import { useDispatch } from "react-redux"
import { useModal } from "./useModal"
import { ContractsLoaded } from "../contracts"
import { ViewAddressInBlockExplorer } from "../components/ViewInBlockExplorer"
import { releaseTokens } from "../actions/web3"
import { withConfirmationModal } from "../components/ConfirmationModal"
import { isEmptyArray } from "../utils/array.utils"

const useReleaseTokens = () => {
  const dispatch = useDispatch()
  const { openConfirmationModal } = useModal()

  const onWithdrawTokens = useCallback(
    async (grant, awaitingPromise) => {
      const { escrowOperatorsToWithdraw } = grant

      if (!isEmptyArray(escrowOperatorsToWithdraw)) {
        const { tokenStakingEscrow } = await ContractsLoaded
        await openConfirmationModal(
          {
            modalOptions: { title: "Are you sure?" },
            title: "You’re about to release tokens.",
            escrowAddress: tokenStakingEscrow.options.address,
            btnText: "release",
            confirmationText: "RELEASE",
          },
          withConfirmationModal(ConfirmWithdrawModal)
        )
      }
      dispatch(releaseTokens(grant, awaitingPromise))
    },
    [dispatch, openConfirmationModal]
  )

  return onWithdrawTokens
}

const ConfirmWithdrawModal = ({ escrowAddress }) => {
  return (
    <>
      <span>You have deposited tokens in the</span>&nbsp;
      <ViewAddressInBlockExplorer
        text="TokenStakingEscrow contract"
        address={escrowAddress}
      />
      <p>
        To withdraw all tokens it may be necessary to confirm more than one
        transaction.
      </p>
    </>
  )
}

export default useReleaseTokens
