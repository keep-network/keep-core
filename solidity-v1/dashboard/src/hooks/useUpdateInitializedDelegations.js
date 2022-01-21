import { useDispatch } from "react-redux"
import { useEffect } from "react"
import moment from "moment"

/**
 * Check if there are any initialized delegations and updates them to be completed
 * @param {array} delegations - array of delegations
 */
const useUpdateInitializedDelegations = (delegations) => {
  const dispatch = useDispatch()
  useEffect(() => {
    const currentDate = moment()
    const initializedDelegations = delegations.filter(
      (delegation) => !!delegation.isInInitializationPeriod
    )

    for (const initializedDelegation of initializedDelegations) {
      if (currentDate.isAfter(initializedDelegation.initializationOverAt)) {
        console.log("update bro")
        dispatch({
          type: "staking/update_delegation",
          payload: {
            operatorAddress: initializedDelegation.operatorAddress,
            values: {
              isInInitializationPeriod: false,
            },
          },
        })
      }
    }
  }, [dispatch, delegations])
}

export default useUpdateInitializedDelegations
