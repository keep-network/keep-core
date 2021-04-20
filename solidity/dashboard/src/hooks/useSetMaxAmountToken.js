import BigNumber from "bignumber.js"
import { useFormikContext } from "formik"
import { KEEP } from "../utils/token.utils"

const useSetMaxAmountToken = (filedName, availableAmount, token = KEEP) => {
  const { setFieldValue } = useFormikContext()

  const setMaxAvailableAmount = () => {
    setFieldValue(
      filedName,
      token.toTokenUnit(availableAmount).toFixed(0, BigNumber.ROUND_DOWN)
    )
  }

  return setMaxAvailableAmount
}

export default useSetMaxAmountToken
