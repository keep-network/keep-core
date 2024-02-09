import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { EVENTS } from "../../constants/events"
import { coveragePoolWithdrawalInitiatedEventEmitted } from "../../actions/coverage-pool"
import { ASSET_POOL_CONTRACT_NAME } from "../../lib/keep/contracts"

export const useSubscribeToCovPoolsWithdrawalInitiatedEvent = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    ASSET_POOL_CONTRACT_NAME,
    EVENTS.COVERAGE_POOLS.WITHDRAWAL_INITIATED,
    (event) => {
      dispatch(coveragePoolWithdrawalInitiatedEventEmitted(event))
    }
  )
}
