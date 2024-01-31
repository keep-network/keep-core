import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { ECDSA_REWARDS_DISTRRIBUTOR_CONTRACT_NAME } from "../../lib/keep/contracts"

export const useSubscribeToECDSARewardsClaimedEvent = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    ECDSA_REWARDS_DISTRRIBUTOR_CONTRACT_NAME,
    "RewardsClaimed",
    async (event) => {
      try {
        const {
          returnValues: { merkleRoot, index, operator, amount },
        } = event

        dispatch({
          type: "rewards/ecdsa_withdrawn",
          payload: { merkleRoot, index, operator, amount },
        })
      } catch (error) {
        console.error(
          `Failed subscribing to Explorer Mode RewardsClaimed event`,
          error
        )
      }
    }
  )
}
