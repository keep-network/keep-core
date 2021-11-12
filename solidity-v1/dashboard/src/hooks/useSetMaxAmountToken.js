import BigNumber from "bignumber.js"
import { useFormikContext } from "formik"
import { KEEP } from "../utils/token.utils"

const useSetMaxAmountTokenFormik = (
  filedName,
  availableAmount,
  token = KEEP,
  decimals = 0
) => {
  const { setFieldValue } = useFormikContext()

  return useSetMaxAmountToken(
    filedName,
    availableAmount,
    setFieldValue,
    token,
    decimals
  )
}

export const useSetMaxAmountToken = (
  filedName,
  availableAmount,
  setFieldValue,
  token = KEEP,
  decimals = 0
) => {
  const setMaxAvailableAmount = () => {
    setFieldValue(
      filedName,
      token
        .toTokenUnit(availableAmount)
        .decimalPlaces(decimals, BigNumber.ROUND_DOWN)
        .toString()
    )
  }

  return setMaxAvailableAmount
}

export default useSetMaxAmountTokenFormik
