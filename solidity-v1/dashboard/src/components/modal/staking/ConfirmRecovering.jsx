import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../Modal"
import Button from "../../Button"
import { FormInputBase } from "../../FormInput"
import { ViewAddressInBlockExplorer } from "../../ViewInBlockExplorer"
import { useTypeTextToConfirmFormik } from "../../../hooks/useTypeTextToConfirmFormik"
import { withBaseModal } from "../withBaseModal"

export const ConfirmRecovering = withBaseModal(
  ({ tokenStakingEscrowAddress, onConfirm, onClose }) => {
    const formik = useTypeTextToConfirmFormik("CLAIM", onConfirm)

    return (
      <>
        <ModalHeader>Claiming tokens</ModalHeader>
        <ModalBody>
          <h3>You’re about to claim your undelegated tokens.</h3>
          <div className="text-grey-60 mt-1">
            <span>
              Because these are tokens from your grant, claiming will deposit
              the undelegated tokens in the
            </span>
            &nbsp;
            <span>
              <ViewAddressInBlockExplorer
                address={tokenStakingEscrowAddress}
                text="TokenStakingEscrow contract."
              />
            </span>
            <p>Withdraw them afterwards, via Release tokens</p>
          </div>
          <form onSubmit={formik.handleSubmit} className="mt-2">
            <FormInputBase
              name="confirmationText"
              type="text"
              onChange={formik.handleChange}
              value={formik.values.confirmationText}
              label={"Type CLAIM to confirm."}
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
            claim
          </Button>
          <Button className="btn btn-unstyled" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
