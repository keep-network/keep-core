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
      &nbsp;ETH
    </>
  )
})

const WithdrawETHModal = ({ operatorAddress, availableETH, closeModal }) => {
  const web3Context = useWeb3Context()
  const { yourAddress, web3 } = web3Context
  const showMessage = useShowMessage()

  const onSubmit = useCallback(
    async (formValues, onTransactionHashCallback) => {
      const { ethAmount: ethToWithdraw } = formValues
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
      <AvailableEthCell availableETH={availableETH} />
      {/* TODO: add a "diamond" icon */}
      <div className="text-validation">
        Withdrawn ETH will go the beneficiary address.
      </div>
      <WithdrawETHFormik
        web3={web3}
        onSubmit={onSubmit}
        yourAddress={yourAddress}
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
        name="ethAmount"
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
  displayName: "WithdrawETHForm",
})(WithdrawETHForm)
