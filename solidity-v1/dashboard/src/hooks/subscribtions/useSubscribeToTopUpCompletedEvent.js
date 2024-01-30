import { useDispatch, useSelector } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { getEventsFromTransaction } from "../../utils/ethereum.utils"
import {
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../../lib/keep/contracts"
import { KeepExplorerMode } from "../../contracts"

export const useSubscribeToTopUpCompletedEvent = () => {
  const grants = useSelector((state) => state.tokenGrants.grants)
  const dispatch = useDispatch()
  const {
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
    [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: { instance: tokenStakingEscrow },
  } = KeepExplorerMode

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "TopUpCompleted",
    async (event) => {
      try {
        const eventsToCheck = [
          [grantContract, "TokenGrantStaked"],
          [tokenStakingEscrow, "DepositRedelegated"],
        ]

        const {
          transactionHash,
          returnValues: { operator, newAmount },
        } = event

        const emittedEvents = await getEventsFromTransaction(
          eventsToCheck,
          transactionHash
        )

        dispatch({
          type: "staking/top_up_completed",
          payload: { operator, newAmount },
        })
        if (
          emittedEvents.DepositRedelegated ||
          emittedEvents.TokenGrantStaked
        ) {
          const { grantId, amount } =
            emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked
          dispatch({
            type: "token-grant/grant_satked",
            payload: { grantId, value: amount },
          })
        }
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode TopUpCompleted event`,
          error
        )
      }
    }
  )
}
