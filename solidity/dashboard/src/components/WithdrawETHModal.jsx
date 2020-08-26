import React, { useCallback } from "react"
import { getErrorsObj } from "../forms/common-validators"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { useWeb3Context } from "./WithWeb3Context"
import * as Icons from "./Icons"
import AvailableEthAmount from "./AvailableEthAmount"
import AvailableETHForm from "./AvailableETHForm"
import {
  withdrawUnbondedEth,
  withdrawUnbondedEthAsManagedGrantee,
} from "../actions/web3"
import { connect } from "react-redux"

const WithdrawETHModal = ({
  operatorAddress,
  availableETH,
  closeModal,
  managedGrantAddress,
  withdrawUnbondedEth,
  withdrawUnbondedEthAsManagedGrantee,
}) => {
  const { web3 } = useWeb3Context()

  const onSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { ethAmount } = formValues
      const weiToWithdraw = web3Utils.toWei(ethAmount.toString(), "ether")

      if (managedGrantAddress) {
        withdrawUnbondedEthAsManagedGrantee(
          weiToWithdraw,
          operatorAddress,
          managedGrantAddress,
          awaitingPromise
        )
      } else {
        withdrawUnbondedEth(weiToWithdraw, operatorAddress, awaitingPromise)
      }
    },
    [
      operatorAddress,
      managedGrantAddress,
      withdrawUnbondedEth,
      withdrawUnbondedEthAsManagedGrantee,
    ]
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
          Withdrawn ETH will go to the beneficiary address.
        </span>
      </div>
      <WithdrawETHFormik
        web3={web3}
        onSubmit={onSubmit}
        availableETH={availableETH}
        onCancel={closeModal}
        submitBtnText="withdraw eth"
      />
    </>
  )
}

const mapDispatchToProps = {
  withdrawUnbondedEth,
  withdrawUnbondedEthAsManagedGrantee,
}

export default React.memo(connect(null, mapDispatchToProps)(WithdrawETHModal))

const WithdrawETHFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    ethAmount: "0",
  }),
  validate: (values, { availableETH }) => {
    const { ethAmount } = values
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
