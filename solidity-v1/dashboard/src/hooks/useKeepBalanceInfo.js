import { useMemo } from "react"
import { add, sub } from "../utils/arithmetics.utils"
import { useSelector } from "react-redux"

const useKeepBalanceInfo = () => {
  const { ownedTokensDelegationsBalance, ownedTokensUndelegationsBalance } =
    useSelector((state) => state.staking)

  const keepToken = useSelector((state) => state.keepTokenBalance)

  const totalOwnedStakedBalance = useMemo(() => {
    return add(
      ownedTokensDelegationsBalance,
      ownedTokensUndelegationsBalance
    ).toString()
  }, [ownedTokensDelegationsBalance, ownedTokensUndelegationsBalance])

  const totalKeepTokenBalance = useMemo(() => {
    return add(totalOwnedStakedBalance, keepToken.value).toString()
  }, [keepToken.value, totalOwnedStakedBalance])

  const totalOwnedUnstakedBalance = useMemo(() => {
    return sub(totalKeepTokenBalance, totalOwnedStakedBalance).toString()
  }, [totalKeepTokenBalance, totalOwnedStakedBalance])

  return {
    totalOwnedStakedBalance,
    totalOwnedUnstakedBalance,
    totalKeepTokenBalance,
  }
}

export default useKeepBalanceInfo
