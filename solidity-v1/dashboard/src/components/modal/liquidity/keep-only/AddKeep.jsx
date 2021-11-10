import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../../Modal"
import Button from "../../../Button"
import TokenAmount from "../../../TokenAmount"
import MaxAmountAddon from "../../../MaxAmountAddon"
import { FormInputBase } from "../../../FormInput"
import { useAddAmountFormik } from "../../../../hooks/useAddAmountFormik"
import { withBaseModal } from "../../withBaseModal"
import { useSetMaxAmountToken } from "../../../../hooks/useSetMaxAmountToken"
import { formatAmount, normalizeAmount } from "../../../../forms/form.utils"

export const AddKeep = withBaseModal(
  ({ availableAmount, onClose, onConfirm }) => {
    const formik = useAddAmountFormik(availableAmount, onConfirm)
    const setMaxAmount = useSetMaxAmountToken(
      "amount",
      availableAmount,
      formik.setFieldValue
    )

    return (
      <>
        <ModalHeader>Deposit KEEP</ModalHeader>
        <ModalBody>
          <h3 className="mb-1">Amount available to deposit.</h3>
          <TokenAmount amount={availableAmount} withIcon withMetricSuffix />
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              type="text"
              name="amount"
              label="Deposit"
              placeholder="0"
              normalize={normalizeAmount}
              format={formatAmount}
              onChange={(event, formatedValue) => {
                formik.handleChange(event)
                formik.setFieldValue("amount", formatedValue)
              }}
              value={formik.values.amount}
              hasError={formik.errors.amount}
              errorMsg={formik.errors.amount}
              inputAddon={
                <MaxAmountAddon onClick={setMaxAmount} text="Max KEEP" />
              }
            />
          </form>
        </ModalBody>
        <ModalFooter>
          <Button
            className="btn btn-primary btn-lg mr-2"
            type="submit"
            onClick={formik.handleSubmit}
            disabled={!(formik.isValid && formik.dirty)}
          >
            deposit keep
          </Button>
          <Button className="btn btn-unstyled text-link" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
