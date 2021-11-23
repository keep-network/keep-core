import { useMemo } from "react"
import { add, sub } from "../utils/arithmetics.utils"
import { useSelector } from "react-redux"

const useGrantedBalanceInfo = () => {
  const { delegations, undelegations } = useSelector((state) => state.staking)

  const { grants } = useSelector((state) => state.tokenGrants)

  const totalGrantedStakedBalance = useMemo(() => {
    return [...delegations, ...undelegations]
      .filter((delegation) => delegation.isFromGrant)
      .map(({ amount }) => amount)
      .reduce(add, "0")
      .toString()
  }, [delegations, undelegations])

  const totalGrantedTokenBalance = useMemo(() => {
    const grantedBalance = grants
      .map(({ amount, released }) => sub(amount, released))
      .reduce(add, "0")
      .toString()
    return grantedBalance
  }, [grants])

  const totalGrantedUnstakedBalance = useMemo(() => {
    return sub(totalGrantedTokenBalance, totalGrantedStakedBalance).toString()
  }, [totalGrantedTokenBalance, totalGrantedStakedBalance])

  return {
    totalGrantedStakedBalance,
    totalGrantedUnstakedBalance,
    totalGrantedTokenBalance,
  }
}

export default useGrantedBalanceInfo
