import { useEffect } from "react"
import moment from "moment"
import { useDispatch } from "react-redux"

/**
 * Check if there are any pending undelegations and updates them to be completed
 * @param {array} undelegations - array of undelegations
 */
const useUpdatePendingUndelegations = (undelegations) => {
  const dispatch = useDispatch()
  useEffect(() => {
    const currentDate = moment()
    const pendingUndelegations = undelegations.filter(
      (undelegation) => !undelegation.canRecoverStake
    )

    for (const pendingUndelegation of pendingUndelegations) {
      if (currentDate.isAfter(pendingUndelegation.undelegationCompleteAt)) {
        dispatch({
          type: "staking/update_undelegation",
          payload: {
            operatorAddress: pendingUndelegation.operatorAddress,
            values: {
              canRecoverStake: true,
            },
          },
        })
      }
    }
  }, [dispatch, undelegations])
}

export default useUpdatePendingUndelegations
