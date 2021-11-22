import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import { ViewAddressInBlockExplorer } from "../../ViewInBlockExplorer"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"

export const ConfirmReleaseTokensFromGrant = withBaseModal(
  ({ escrowContractAddress, onConfirm, onClose }) => {
    const formik = useTypeTextToConfirmFormik("RELEASE", onConfirm)

    return (
      <>
        <ModalHeader>Are you sure?</ModalHeader>
        <ModalBody>
          <h3>You’re about to release tokens.</h3>
          <div className="text-grey-60 mt-1">
            You have deposited tokens in the&nbsp;
            <ViewAddressInBlockExplorer
              text="TokenStakingEscrow contract"
              address={escrowContractAddress}
            />
            . To withdraw all tokens it may be necessary to confirm more than
            one transaction.
          </div>
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              name="confirmationText"
              type="text"
              onChange={formik.handleChange}
              value={formik.values.confirmationText}
              label={"Type RELEASE to confirm."}
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
            release
          </Button>
          <Button className="btn btn-unstyled" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
