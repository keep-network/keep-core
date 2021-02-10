import BigNumber from "bignumber.js"
import { useFormikContext } from "formik"
import { toTokenUnit } from "../utils/token.utils"

const useSetMaxAmountToken = (filedName, availableAmount) => {
  const { setFieldValue } = useFormikContext()

  const setMaxAvailableAmount = () => {
    setFieldValue(
      filedName,
      toTokenUnit(availableAmount).toFixed(0, BigNumber.ROUND_DOWN)
    )
  }

  return setMaxAvailableAmount
}

export default useSetMaxAmountToken
