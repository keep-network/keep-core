import { useMemo } from "react"
import { isSameEthAddress } from "../utils/general.utils"
import { useSelector } from "react-redux"

const useDelegationsWithTAuthData = () => {
  const { delegations } = useSelector((state) => state.staking)

  const thresholdAuthState = useSelector(
    (state) => state.thresholdAuthorization
  )

  const delegationsWithTAuthData = useMemo(() => {
    if (delegations.length === 0) return []
    if (thresholdAuthState.data?.length === 0) return delegations
    return delegations.map((delegation) => {
      const tAuthData = thresholdAuthState.authData.find((data) => {
        return isSameEthAddress(
          data.operatorAddress,
          delegation.operatorAddress
        )
      })
      return {
        ...delegation,
        isTStakingContractAuthorized: tAuthData?.contract.isAuthorized || false,
        isStakedToT: tAuthData?.isStakedToT || false,
      }
    })
  }, [delegations, thresholdAuthState.authData])

  return delegationsWithTAuthData
}

export default useDelegationsWithTAuthData
