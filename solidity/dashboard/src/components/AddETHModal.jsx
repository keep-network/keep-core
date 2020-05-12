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

const AddEthModal = ({ operatorAddress, closeModal }) => {
  const web3Context = useWeb3Context()
  const { yourAddress, web3 } = web3Context
  const showMessage = useShowMessage()

  const onSubmit = useCallback(
    async (formValues, onTransactionHashCallback) => {
      const { ethAmount: value } = formValues
      try {
        await tbtcAuthorizationService.depositEthForOperator(
          web3Context,
          { operatorAddress, value },
          onTransactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Add ETH for operator transaction successfully completed",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Add ETH for operator action has failed ",
          content: error.message,
        })
        throw error
      }
    },
    [operatorAddress, showMessage, web3Context]
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
        closeModal={closeModal}
      />
    </>
  )
}

export default React.memo(AddEthModal)

const AddETHForm = ({ onSubmit, closeModal, ...formikProps }) => {
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
          add eth
        </SubmitButton>
        <span onClick={closeModal} className="ml-1 text-link">
          Cancel
        </span>
      </div>
    </form>
  )
}
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
})(AddETHForm)
