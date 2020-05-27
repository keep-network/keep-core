import React, { useCallback } from "react"
import { SubmitButton } from "./Button"
import { useCustomOnSubmitFormik } from "../hooks/useCustomOnSubmitFormik"
import FormInput from "./FormInput"
import { getErrorsObj } from "../forms/common-validators"
import { withFormik } from "formik"
import { colors } from "../constants/colors"
import web3Utils from "web3-utils"
import { useWeb3Context } from "./WithWeb3Context"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useShowMessage, messageType } from "./Message"
import * as Icons from "./Icons"

const AvailableEthCell = React.memo(({ availableETH }) => {
  return (
    <>
      <span
        className="text-big text-grey-70"
        style={{
          textAlign: "right",
          padding: "0.25rem 1rem",
          paddingLeft: "2rem",
          borderRadius: "100px",
          border: `1px solid ${colors.grey20}`,
          backgroundColor: `${colors.grey10}`,
        }}
      >
        {availableETH}
      </span>
      <span style={{ color: `${colors.grey60}` }}>&nbsp;ETH</span>
    </>
  )
})

const WithdrawETHModal = ({ operatorAddress, availableETH, closeModal }) => {
  const web3Context = useWeb3Context()
  const { web3 } = web3Context
  const showMessage = useShowMessage()

  const onSubmit = useCallback(
    async (formValues, onTransactionHashCallback) => {
      const { ethToWithdraw: ethToWithdraw } = formValues
      try {
        await tbtcAuthorizationService.withdrawUnbondedEth(
          web3Context,
          { operatorAddress, ethToWithdraw },
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
        <AvailableEthCell availableETH={availableETH} />
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
      />
    </>
  )
}

export default React.memo(WithdrawETHModal)

const WithdrawETHForm = ({ onSubmit, closeModal, ...formikProps }) => {
  const onSubmitBtn = useCustomOnSubmitFormik(onSubmit)

  return (
    <form>
      <FormInput
        name="ethToWithdraw"
        type="text"
        label="ETH Amount"
        placeholder="0"
      />
      <div
        className="flex row center mt-2"
        style={{
          borderTop: `1px solid ${colors.grey20}`,
          margin: "0 -2rem",
          padding: "2rem 2rem 0",
        }}
      >
        <SubmitButton
          className="btn btn-primary"
          type="submit"
          onSubmitAction={onSubmitBtn}
          withMessageActionIsPending={false}
          triggerManuallyFetch={true}
          disabled={!formikProps.dirty}
        >
          withdraw eth
        </SubmitButton>
        <span onClick={closeModal} className="ml-1 text-link">
          Cancel
        </span>
      </div>
    </form>
  )
}
const WithdrawETHFormik = withFormik({
  validateOnChange: false,
  validateOnBlur: false,
  mapPropsToValues: () => ({
    ethToWithdraw: "0",
  }),
  validate: (values, { availableETH }) => {
    const { ethToWithdraw } = values
    const errors = {}

    if (isNaN(ethToWithdraw)) {
      errors.ethToWithdraw = "A valid number must be provided"
      return getErrorsObj(errors)
    }

    const unbondedValueInWei = web3Utils.toBN(
      web3Utils.toWei(availableETH ? availableETH.toString() : "0")
    )
    const valueToWithdrawInWei = web3Utils.toBN(
      web3Utils.toWei(ethToWithdraw.toString())
    )

    if (valueToWithdrawInWei.gt(unbondedValueInWei)) {
      errors.ethToWithdraw = `The withdrawable amount should be less than ${availableETH} Eth`
    }

    if (valueToWithdrawInWei.lte(web3Utils.toBN(0))) {
      errors.ethToWithdraw = "The withdrawable amount should be greater than 0"
    }

    return getErrorsObj(errors)
  },
  displayName: "WithdrawETHForm",
})(WithdrawETHForm)
