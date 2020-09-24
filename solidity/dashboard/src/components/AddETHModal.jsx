import React, { useCallback } from "react"
import AvailableETHForm from "./AvailableETHForm"
import { getErrorsObj } from "../forms/common-validators"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { depositEthForOperator } from "../actions/web3"
import { connect } from "react-redux"
import { useWeb3Context } from "./WithWeb3Context"

const AddEthModal = ({
  operatorAddress,
  closeModal,
  depositEthForOperator,
}) => {
  const { web3, yourAddress } = useWeb3Context()

  const onSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { ethAmount } = formValues
      const weiAmount = web3Utils.toWei(ethAmount.toString(), "ether")

      depositEthForOperator(operatorAddress, weiAmount, awaitingPromise)
    },
    [operatorAddress, depositEthForOperator]
  )

  return (
    <>
      <h4 style={{ marginBottom: "0.5rem" }}>Enter an amount of ETH</h4>
      <div className="text-big text-grey-60 mb-3">
        This amount of ETH will be available for bonding. An available balance
        of ETH allows you to be selected for signing groups, which bonds ETH.
      </div>
      <AddETHFormik
        web3={web3}
        onSubmit={onSubmit}
        yourAddress={yourAddress}
        onCancel={closeModal}
        submitBtnText="add eth"
      />
    </>
  )
}

const mapDispatchToProps = {
  depositEthForOperator,
}

export default React.memo(connect(null, mapDispatchToProps)(AddEthModal))

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
