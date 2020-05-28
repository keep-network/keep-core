import React, { useCallback } from "react"
import { getErrorsObj } from "../forms/common-validators"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { useWeb3Context } from "./WithWeb3Context"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useShowMessage, messageType } from "./Message"
import * as Icons from "./Icons"
import AvailableEthAmount from "./AvailableEthAmount"
import AvailableETHForm from "./AvailableETHForm"

const WithdrawETHModal = ({ operatorAddress, availableETH, closeModal }) => {
  const web3Context = useWeb3Context()
  const { web3 } = web3Context
  const showMessage = useShowMessage()

  const onSubmit = useCallback(
    async (formValues, onTransactionHashCallback) => {
      const { ethAmount } = formValues
      try {
        await tbtcAuthorizationService.withdrawUnbondedEth(
          web3Context,
          { operatorAddress, ethAmount },
          onTransactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Withdrawal of ETH successfully completed",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Withdrawal of ETH has failed ",
          content: error.message,
        })
        throw error
      }
    },
    [operatorAddress, showMessage, web3Context]
  )

  return (
    <>
      <h4 style={{ marginBottom: "0.5rem" }}>Amount available to withdraw.</h4>
      <div className="mt-1">
        <AvailableEthAmount availableETH={availableETH} />
      </div>
      <div className="text-validation mb-1 mt-2 flex row center">
        <Icons.Diamond />
        <span className="pl-1">
          Withdrawn ETH will go the beneficiary address.
        </span>
      </div>
      <WithdrawETHFormik
        web3={web3}
        onSubmit={onSubmit}
        availableETH={availableETH}
        closeModal={closeModal}
        action="withdraw"
      />
    </>
  )
}

export default React.memo(WithdrawETHModal)

const WithdrawETHFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    ethAmount: "0",
  }),
  validate: (values, { availableETH }) => {
    const { ethAmount } = values

    console.log("ethAmount in withdrawal", ethAmount)
    const errors = {}

    if (isNaN(ethAmount)) {
      errors.ethAmount = "A valid number must be provided"
      return getErrorsObj(errors)
    }

    const unbondedValueInWei = web3Utils.toBN(
      web3Utils.toWei(availableETH ? availableETH.toString() : "0")
    )
    const valueToWithdrawInWei = web3Utils.toBN(
      web3Utils.toWei(ethAmount.toString())
    )

    if (valueToWithdrawInWei.gt(unbondedValueInWei)) {
      errors.ethAmount = `The withdrawable amount should be less than ${availableETH} Eth`
    }

    if (valueToWithdrawInWei.lte(web3Utils.toBN(0))) {
      errors.ethAmount = "The withdrawable amount should be greater than 0"
    }

    return getErrorsObj(errors)
  },
  displayName: "WithdrawETHForm",
})(AvailableETHForm)
