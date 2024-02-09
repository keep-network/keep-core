import { useDispatch, useSelector } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { MODAL_TYPES } from "../../constants/constants"
import { getEventsFromTransaction } from "../../utils/ethereum.utils"
import { isSameEthAddress } from "../../utils/general.utils"
import { showModal } from "../../actions/modal"
import {
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../../lib/keep/contracts"
import { KeepExplorerMode } from "../../contracts"

export const useSubscribeToTopUpInitiatedEvent = () => {
  const delegations = useSelector((state) => state.tokenGrants.delegations)
  const dispatch = useDispatch()
  const {
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
    [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: { instance: tokenStakingEscrow },
  } = KeepExplorerMode

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "TopUpInitiated",
    async (event) => {
      try {
        // Other events may also be emitted with the `TopUpInitiated` event.
        const eventsToCheck = [
          [grantContract, "TokenGrantStaked"],
          [tokenStakingEscrow, "DepositRedelegated"],
        ]

        const {
          transactionHash,
          returnValues: { operator, topUp },
        } = event

        const emittedEvents = await getEventsFromTransaction(
          eventsToCheck,
          transactionHash
        )

        // Find existing delegation in the app context
        const delegation = delegations.find(({ operatorAddress }) =>
          isSameEthAddress(operatorAddress, operator)
        )

        if (delegation) {
          dispatch(
            showModal({
              modalType: MODAL_TYPES.TopUpInitiatedConfirmation,
              modalProps: {
                addedAmount: topUp,
                currentAmount: delegation.amount,
                authorizerAddress: delegation.authorizerAddress,
                beneficiary: delegation.beneficiary,
                operatorAddress: delegation.operatorAddress,
              },
            })
          )
          dispatch({
            type: "staking/top_up_initiated",
            payload: { operator, topUp },
          })

          if (
            emittedEvents.DepositRedelegated ||
            emittedEvents.TokenGrantStaked
          ) {
            const { grantId, amount } =
              emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked
            dispatch({
              type: "token-grant/grant_staked",
              payload: { grantId, value: amount },
            })
          }
        }
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode TopUpInitiated event`,
          error
        )
      }
    }
  )
}
