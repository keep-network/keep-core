import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import {
  OLD_TOKEN_STAKING_CONTRACT_NAME,
  STAKING_PORT_BACKER_CONTRACT_NAME,
} from "../../lib/keep/contracts"

export const useSubscribeToCopyStakeEvents = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    OLD_TOKEN_STAKING_CONTRACT_NAME,
    "Undelegated",
    async (event) => {
      try {
        const {
          returnValues: { operator },
        } = event

        dispatch({
          type: "copy-stake/remove_old_delegation",
          payload: operator,
        })
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode Undelegated event`,
          error
        )
      }
    }
  )

  useSubscribeToExplorerModeContractEvent(
    STAKING_PORT_BACKER_CONTRACT_NAME,
    "StakeCopied",
    async (event) => {
      try {
        const {
          returnValues: { operator },
        } = event

        dispatch({
          type: "copy-stake/remove_old_delegation",
          payload: operator,
        })
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode StakeCopied event`,
          error
        )
      }
    }
  )

  useSubscribeToExplorerModeContractEvent(
    OLD_TOKEN_STAKING_CONTRACT_NAME,
    "RecoveredStake",
    async (event) => {
      try {
        const {
          returnValues: { operator },
        } = event

        dispatch({
          type: "copy-stake/remove_old_delegation",
          payload: operator,
        })
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode Recovered event`,
          error
        )
      }
    },
    {},
    true
  )
}
