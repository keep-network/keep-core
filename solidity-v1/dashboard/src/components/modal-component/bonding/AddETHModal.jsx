import React, { useCallback } from "react"
import { useDispatch } from "react-redux"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { ModalHeader, ModalBody } from "../Modal"
import AvailableETHForm from "../../AvailableTokenForm"
import { useWeb3Context } from "../../WithWeb3Context"
import { depositEthForOperator } from "../../../actions/web3"
import { getErrorsObj } from "../../../forms/common-validators"
import { withBaseModal } from "../withBaseModal"

const AddETHModalBase = ({ operatorAddress, onClose }) => {
  const dispatch = useDispatch()
  const { web3, yourAddress } = useWeb3Context()

  const onSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { ethAmount } = formValues
      const weiAmount = web3Utils.toWei(ethAmount.toString(), "ether")

      dispatch(
        depositEthForOperator(operatorAddress, weiAmount, awaitingPromise)
      )
    },
    [operatorAddress, dispatch]
  )

  return (
    <>
      <ModalHeader>Add ETH</ModalHeader>
      <ModalBody>
        <h4 style={{ marginBottom: "0.5rem" }}>Enter an amount of ETH</h4>
        <div className="text-big text-grey-60 mb-3">
          This amount of ETH will be available for bonding. An available balance
          of ETH allows you to be selected for signing groups, which bonds ETH.
        </div>
        <AddETHFormik
          web3={web3}
          onSubmit={onSubmit}
          yourAddress={yourAddress}
          onCancel={onClose}
          submitBtnText="add eth"
        />
      </ModalBody>
    </>
  )
}

export const AddETHModal = React.memo(withBaseModal(AddETHModalBase))

const AddETHFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    ethAmount: "0",
  }),
  validate: (values, { yourAddress, web3 }) => {
    return web3.eth.getBalance(yourAddress).then((balance) => {
      const { ethAmount } = values
      const errors = {}

      const ethBalance = web3Utils.toBN(balance)
      const valueInWei = web3Utils.toBN(
        web3Utils.toWei(ethAmount ? ethAmount.toString() : "0")
      )

      if (!ethAmount) {
        errors.ethAmount = "Required"
      } else if (ethBalance.isZero()) {
        errors.ethAmount = "Account ETH balance should be greater than 0"
      } else if (valueInWei.gt(ethBalance)) {
        errors.ethAmount = `The value should be less than ${web3Utils.fromWei(
          ethBalance.toString(),
          "ether"
        )}`
      } else if (valueInWei.lte(web3Utils.toBN(0))) {
        errors.ethAmount = "The value should be greater than 0"
      }

      return getErrorsObj(errors)
    })
  },
  displayName: "AddEthForm",
})(AvailableETHForm)
