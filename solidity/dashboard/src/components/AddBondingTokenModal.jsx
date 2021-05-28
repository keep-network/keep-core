import React, { useCallback } from "react"
import { useDispatch } from "react-redux"
import AvailableBondingTokenForm from "./AvailableBondingTokenForm"
import { getErrorsObj } from "../forms/common-validators"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { useWeb3Context } from "./WithWeb3Context"

const AddBondingTokenModal = ({
  operatorAddress,
  closeModal
}) => {
  const dispatch = useDispatch()
  const { web3, yourAddress } = useWeb3Context()

  const handleDepositForOperator = useCallback((amount, meta) => {
    dispatch({
      type: 'bonding/deposit_start',
      payload: { amount, operatorAddress },
      meta
    })
  }, [dispatch])

  const onSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { tokenAmount } = formValues
      const weiAmount = web3Utils.toWei(tokenAmount.toString(), "ether")

      handleDepositForOperator(weiAmount, awaitingPromise)
    },
    [operatorAddress, handleDepositForOperator]
  )

  return (
    <>
      <h4 style={{ marginBottom: "0.5rem" }}>Enter an amount of ERC20</h4>
      <div className="text-big text-grey-60 mb-3">
        This amount of ERC20 will be available for bonding. An available balance
        of ERC20 allows you to be selected for signing groups, which bonds ERC20.
      </div>
      <AddBondingTokenFormik
        web3={web3}
        onSubmit={onSubmit}
        yourAddress={yourAddress}
        onCancel={closeModal}
        submitBtnText="add erc20"
      />
    </>
  )
}

export default React.memo(AddBondingTokenModal)

const AddBondingTokenFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    tokenAmount: "0",
  }),
  validate: (values) => {
    const { tokenAmount } = values
    const errors = {}

    const valueInWei = web3Utils.toBN(
      web3Utils.toWei(tokenAmount ? tokenAmount.toString() : "0")
    )

    if (!tokenAmount) {
      errors.tokenAmount = "Required"
    } else if (valueInWei.lte(web3Utils.toBN(0))) {
      errors.tokenAmount = "The value should be greater than 0"
    }

    return getErrorsObj(errors)
  },
  displayName: "AddBondingTokenForm",
})(AvailableBondingTokenForm)
