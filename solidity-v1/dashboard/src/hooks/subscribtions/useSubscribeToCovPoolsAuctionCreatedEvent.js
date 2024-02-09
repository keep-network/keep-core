import { useDispatch } from "react-redux"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { riskManagerAuctionCreatedEventEmitted } from "../../actions/coverage-pool"
import { RISK_MANAGER_V1_CONTRACT_NAME } from "../../lib/keep/contracts"

export const useSubscribeToCovPoolsAuctionCreatedEvent = () => {
  const dispatch = useDispatch()

  useSubscribeToExplorerModeContractEvent(
    RISK_MANAGER_V1_CONTRACT_NAME,
    "AuctionCreated",
    (event) => {
      dispatch(riskManagerAuctionCreatedEventEmitted(event))
    }
  )
}
