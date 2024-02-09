import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { THRESHOLD_STAKING_CONTRACT_NAME } from "../../lib/keep/contracts"
import { thresholdStakeKeepEventEmitted } from "../../actions/keep-to-t-staking"

export const useSubscribeToThresholdStakeKeepEvent = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    THRESHOLD_STAKING_CONTRACT_NAME,
    "Staked",
    (event) => {
      dispatch(thresholdStakeKeepEventEmitted(event))
    }
  )
}
