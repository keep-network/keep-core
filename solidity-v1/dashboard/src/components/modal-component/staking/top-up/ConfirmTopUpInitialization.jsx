import React from "react"
import { DelegationDetails } from "./components"
import { ModalHeader, ModalBody } from "../../Modal"
import { withBaseModal } from "../../withBaseModal"
import * as Icons from "../../../Icons"
import TokenAmount from "../../../TokenAmount"
import { KEEP } from "../../../../utils/token.utils"
import { ModalFooter } from "../../Modal"
import Button from "../../../Button"
import { useTypeTextToConfirmFormik } from "../../../../hooks/useTypeTextToConfirmFormik"
import { FormInputBase } from "../../../FormInput"

export const ConfirmTopUpInitialization = withBaseModal(
  ({
    authorizerAddress,
    beneficiary,
    operatorAddress,
    amountToAdd,
    newAmount,
    onConfirm,
    onClose,
  }) => {
    const formik = useTypeTextToConfirmFormik("CONFIRM", onConfirm)
    console.log("formik.isValid", formik)
    return (
      <>
        <ModalHeader>Add KEEP</ModalHeader>
        <ModalBody>
          <TokenAmount
            amount={amountToAdd}
            withIcon
            icon={Icons.Plus}
            iconProps={{ width: 24, height: 24, className: "plus-icon" }}
          />
          <h4 className="text-grey-70 mb-1">
            New delegation balance: {KEEP.displayAmountWithSymbol(newAmount)}
          </h4>
          <DelegationDetails
            authorizerAddress={authorizerAddress}
            beneficiary={beneficiary}
            operatorAddress={operatorAddress}
          />
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              name="confirmationText"
              type="text"
              onChange={formik.handleChange}
              value={formik.values.confirmationText}
              label={"Type CONFIRM to add KEEP."}
              placeholder=""
              hasError={formik.errors.confirmationText}
              errorMsg={formik.errors.confirmationText}
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
            confirm
          </Button>
          <Button className="btn btn-unstyled text-link" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
