import React, { useCallback } from "react"
import { getErrorsObj } from "../forms/common-validators"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { useWeb3Context } from "./WithWeb3Context"
import * as Icons from "./Icons"
import AvailableBondingTokenAmount from "./AvailableBondingTokenAmount"
import AvailableBondingTokenForm from "./AvailableBondingTokenForm"
import {
  withdrawUnbondedBondingToken,
  withdrawUnbondedBondingTokenAsManagedGrantee,
} from "../actions/web3"
import { connect } from "react-redux"

const WithdrawBondingTokenModal = ({
  operatorAddress,
  availableTokensInWei,
  availableTokens,
  closeModal,
  managedGrantAddress,
  withdrawUnbondedBondingToken,
  withdrawUnbondedBondingTokenAsManagedGrantee,
}) => {
  const { web3 } = useWeb3Context()

  const onSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { tokenAmount } = formValues
      const weiToWithdraw = web3Utils.toWei(tokenAmount.toString(), "ether")

      if (managedGrantAddress) {
        withdrawUnbondedBondingTokenAsManagedGrantee(
          weiToWithdraw,
          operatorAddress,
          managedGrantAddress,
          awaitingPromise
        )
      } else {
        withdrawUnbondedBondingToken(weiToWithdraw, operatorAddress, awaitingPromise)
      }
    },
    [
      operatorAddress,
      managedGrantAddress,
      withdrawUnbondedBondingToken,
      withdrawUnbondedBondingTokenAsManagedGrantee,
    ]
  )

  return (
    <>
      <h4 style={{ marginBottom: "0.5rem" }}>Amount available to withdraw.</h4>
      <div className="mt-1">
        <AvailableBondingTokenAmount
          availableTokensInWei={availableTokensInWei}
          availableTokens={availableTokens}
        />
      </div>
      <div className="text-validation mb-1 mt-2 flex row center">
        <Icons.Diamond />
        <span className="pl-1">
          Withdrawn ERC20 will go to the beneficiary address.
        </span>
      </div>
      <WithdrawBondingTokenModalFormik
        web3={web3}
        onSubmit={onSubmit}
        availableTokens={availableTokens}
        onCancel={closeModal}
        submitBtnText="withdraw ERC20"
      />
    </>
  )
}

const mapDispatchToProps = {
  withdrawUnbondedBondingToken,
  withdrawUnbondedBondingTokenAsManagedGrantee,
}

export default React.memo(connect(null, mapDispatchToProps)(WithdrawBondingTokenModal))

const WithdrawBondingTokenModalFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    tokenAmount: "0",
  }),
  validate: (values, { availableTokens }) => {
    const { tokenAmount } = values
    const errors = {}

    if (isNaN(tokenAmount)) {
      errors.tokenAmount = "A valid number must be provided"
      return getErrorsObj(errors)
    }

    const unbondedValueInWei = web3Utils.toBN(
      web3Utils.toWei(availableTokens ? availableTokens.toString() : "0")
    )
    const valueToWithdrawInWei = web3Utils.toBN(
      web3Utils.toWei(tokenAmount.toString())
    )

    if (valueToWithdrawInWei.gt(unbondedValueInWei)) {
      errors.tokenAmount = `The withdrawable amount should be less than ${availableTokens} Eth`
    }

    if (valueToWithdrawInWei.lte(web3Utils.toBN(0))) {
      errors.tokenAmount = "The withdrawable amount should be greater than 0"
    }

    return getErrorsObj(errors)
  },
  displayName: "WithdrawBondingTokenModalForm",
})(AvailableBondingTokenForm)
