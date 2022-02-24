import React from "react"
import { useFormik } from "formik"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import TokenAmount from "../../TokenAmount"
import Button from "../../Button"
import AddAmountToWithdrawalForm from "../../coverage-pools/AddAmountToWithdrawalForm"
import { withBaseModal } from "../withBaseModal"
import { covKEEP, KEEP } from "../../../utils/token.utils"
import {
  validateAmountInRange,
  getErrorsObj,
} from "../../../forms/common-validators"
import { Keep } from "../../../contracts"

const ReInitiateWithdrawComponent = ({
  pendingWithdrawal,
  availableToWithdraw,
  totalValueLocked,
  covTotalSupply,
  onConfirm,
  onClose,
}) => {
  const formik = useFormik({
    validateOnChange: true,
    validateOnBlur: true,
    initialValues: {
      amount: "0",
    },
    validate: (values) => {
      const { amount } = values
      const errors = {}

      errors.tokenAmount = validateAmountInRange(
        amount,
        availableToWithdraw,
        0,
        covKEEP
      )

      return getErrorsObj(errors)
    },
    onSubmit: (values) => {
      onConfirm(values)
    },
  })
  return (
    <>
      <ModalHeader>Re-initiate withdrawal</ModalHeader>
      <ModalBody>
        <h3 className="mb-1">You are about to re-initiate this withdrawal:</h3>
        <TokenAmount amount={pendingWithdrawal} token={covKEEP} />
        <TokenAmount
          amount={Keep.coveragePoolV1.estimatedBalanceFor(
            pendingWithdrawal,
            covTotalSupply,
            totalValueLocked
          )}
          token={KEEP}
          amountClassName="text-grey-60"
          symbolClassName="text-grey-60"
        />
        <AddAmountToWithdrawalForm
          tokenAmount={availableToWithdraw}
          onSubmit={formik.handleSubmit}
          setMaxAmount={formik.setFieldValue}
          totalValueLocked={totalValueLocked}
          covTotalSupply={covTotalSupply}
          inputProps={{
            name: "amount",
            onChange: (event, formattedValue) => {
              formik.handleChange(event)
              formik.setFieldValue("amount", formattedValue)
            },
            value: formik.values.amount,
            hasError: formik.errors.amount,
            errorMsg: formik.errors.amount,
          }}
        />
      </ModalBody>
      <ModalFooter>
        <Button
          className="btn btn-primary btn-lg mr-2"
          type="submit"
          onClick={formik.handleSubmit}
          disabled={!formik.isValid}
        >
          continue
        </Button>
        <Button className="btn btn-unstyled text-link" onClick={onClose}>
          Cancel
        </Button>
      </ModalFooter>
    </>
  )
}

const ReInitiateWithdrawWithBaseModal = withBaseModal(
  ReInitiateWithdrawComponent
)

export const ReInitiateWithdraw = (props) => (
  <ReInitiateWithdrawWithBaseModal size="lg" {...props} />
)
