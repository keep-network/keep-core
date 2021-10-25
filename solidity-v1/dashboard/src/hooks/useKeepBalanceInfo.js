import { useMemo } from "react"
import { add } from "../utils/arithmetics.utils"
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

  return {
    totalOwnedStakedBalance,
    totalKeepTokenBalance,
  }
}

export default useKeepBalanceInfo
