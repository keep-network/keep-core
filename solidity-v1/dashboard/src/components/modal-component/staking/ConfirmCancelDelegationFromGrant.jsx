import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import { ViewAddressInBlockExplorer } from "../../ViewInBlockExplorer"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"

export const ConfirmCancelDelegationFromGrant = withBaseModal(
  ({ tokenStakingEscrowAddress, onConfirm, onCancel }) => {
    const formik = useTypeTextToConfirmFormik("CANCEL", onConfirm)

    return (
      <>
        <ModalHeader>Are you sure?</ModalHeader>
        <ModalBody>
          <h3>You’re about to cancel tokens.</h3>
          <div className="text-grey-60">
            <span>Canceling will deposit delegated tokens in the</span>
            &nbsp;
            <span>
              <ViewAddressInBlockExplorer
                address={tokenStakingEscrowAddress}
                text="TokenStakingEscrow contract."
              />
            </span>
            <p>You can withdraw them via Release tokens.</p>
          </div>
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              name="confirmationText"
              type="text"
              onChange={formik.handleChange}
              value={formik.values.confirmationText}
              label={"Type CANCEL to confirm."}
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
          >
            delegate
          </Button>
          <Button className="btn btn-unstyled" onClick={onCancel}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
