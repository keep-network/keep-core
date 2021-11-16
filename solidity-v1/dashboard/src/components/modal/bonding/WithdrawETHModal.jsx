import React, { useCallback } from "react"
import { useDispatch } from "react-redux"
import { withFormik } from "formik"
import web3Utils from "web3-utils"
import { ModalHeader, ModalBody } from "../Modal"
import { withBaseModal } from "../withBaseModal"
import AvailableTokenForm from "../../AvailableTokenForm"
import MaxAmountAddon from "../../MaxAmountAddon"
import AvailableEthAmount from "../../AvailableEthAmount"
import * as Icons from "../../Icons"
import { useWeb3Context } from "../../WithWeb3Context"
import useSetMaxAmountToken from "../../../hooks/useSetMaxAmountToken"
import { ETH } from "../../../utils/token.utils"
import { getErrorsObj } from "../../../forms/common-validators"
import {
  withdrawUnbondedEth,
  withdrawUnbondedEthAsManagedGrantee,
} from "../../../actions/web3"

const WithdrawETHModalBase = ({
  operatorAddress,
  availableETHInWei,
  availableETH,
  onClose,
  managedGrantAddress,
}) => {
  const { web3 } = useWeb3Context()
  const dispatch = useDispatch()

  const onSubmit = useCallback(
    async (formValues, awaitingPromise) => {
      const { ethAmount } = formValues
      const weiToWithdraw = web3Utils.toWei(ethAmount.toString(), "ether")

      if (managedGrantAddress) {
        dispatch(
          withdrawUnbondedEthAsManagedGrantee(
            weiToWithdraw,
            operatorAddress,
            managedGrantAddress,
            awaitingPromise
          )
        )
      } else {
        dispatch(
          withdrawUnbondedEth(weiToWithdraw, operatorAddress, awaitingPromise)
        )
      }
    },
    [operatorAddress, managedGrantAddress, dispatch]
  )

  return (
    <>
      <ModalHeader>Withdraw ETH</ModalHeader>
      <ModalBody>
        <h4 style={{ marginBottom: "0.5rem" }}>
          Amount available to withdraw.
        </h4>
        <div className="mt-1">
          <AvailableEthAmount
            availableETHInWei={availableETHInWei}
            availableETH={availableETH}
          />
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
          availableETHInWei={availableETHInWei}
          onCancel={onClose}
          submitBtnText="withdraw eth"
        />
      </ModalBody>
    </>
  )
}

export const WithdrawETHModal = React.memo(withBaseModal(WithdrawETHModalBase))

const WithdrawETHForm = (props) => {
  const { availableETHInWei } = props
  const setMaxAmount = useSetMaxAmountToken(
    "ethAmount",
    availableETHInWei,
    ETH,
    ETH.decimals
  )
  return (
    <AvailableTokenForm
      formInputProps={{
        label: "Withdraw",
        name: "ethAmount",
        inputAddon: <MaxAmountAddon onClick={setMaxAmount} text="Max Amount" />,
      }}
      {...props}
    />
  )
}

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
})(WithdrawETHForm)
