import BigNumber from "bignumber.js"
import { useFormikContext } from "formik"
import { KEEP } from "../utils/token.utils"

const useSetMaxAmountToken = (
  filedName,
  availableAmount,
  token = KEEP,
  decimals = 0
) => {
  const { setFieldValue } = useFormikContext()

  const setMaxAvailableAmount = () => {
    setFieldValue(
      filedName,
      token.toTokenUnit(availableAmount).toFixed(decimals, BigNumber.ROUND_DOWN)
    )
  }

  return setMaxAvailableAmount
}

export default useSetMaxAmountToken
