import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { EVENTS } from "../../constants/events"
import { coveragePoolWithdrawalCompletedEventEmitted } from "../../actions/coverage-pool"
import { ASSET_POOL_CONTRACT_NAME } from "../../lib/keep/contracts"

export const useSubscribeToCovPoolsWithdrawalCompletedEvent = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    ASSET_POOL_CONTRACT_NAME,
    EVENTS.COVERAGE_POOLS.WITHDRAWAL_COMPLETED,
    (event) => {
      dispatch(coveragePoolWithdrawalCompletedEventEmitted(event))
    }
  )
}
