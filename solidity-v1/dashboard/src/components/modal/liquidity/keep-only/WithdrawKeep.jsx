import React from "react"
import { AmountTile } from "./components"
import { ModalHeader, ModalBody, ModalFooter } from "../../Modal"
import MaxAmountAddon from "../../../MaxAmountAddon"
import { FormInputBase } from "../../../FormInput"
import Button from "../../../Button"
import * as Icons from "../../../Icons"
import { useAddAmountFormik } from "../../../../hooks/useAddAmountFormik"
import { withBaseModal } from "../../withBaseModal"
import { useSetMaxAmountToken } from "../../../../hooks/useSetMaxAmountToken"
import { formatAmount, normalizeAmount } from "../../../../forms/form.utils"

export const WithdrawKeep = withBaseModal(
  ({ availableAmount, rewardedAmount, onClose, onConfirm }) => {
    const formik = useAddAmountFormik(availableAmount, onConfirm)
    const setMaxAmount = useSetMaxAmountToken(
      "amount",
      availableAmount,
      formik.setFieldValue
    )

    return (
      <>
        <ModalHeader>Withdraw Locked KEEP</ModalHeader>
        <ModalBody>
          <h3 className="mb-1">Amount available to withdraw.</h3>
          <div className="flex row mb-2">
            <AmountTile title="deposited" amount={availableAmount} />
            <AmountTile
              title="rewarded"
              amount={rewardedAmount}
              icon={Icons.Rewards}
            />
          </div>
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              type="text"
              name="amount"
              label="Withdraw"
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
            withdraw keep
          </Button>
          <Button className="btn btn-unstyled text-link" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
