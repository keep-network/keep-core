import React from "react"
import { ModalBody } from "../Modal"
import { ModalHeader } from "../Modal"
import { withBaseModal } from "../withBaseModal"

export const DelegationAlreadyExists = withBaseModal(({ operatorAddress }) => {
  return (
    <>
      <ModalHeader>Delegation already exists</ModalHeader>
      <ModalBody>
        Delegate tokens for a different operator address or top-up the existing
        delegation for <strong>{operatorAddress}</strong>
        &nbsp;operartor via <strong>ADD KEEP</strong> button under&nbsp;
        <strong>Delegations</strong> table.
      </ModalBody>
    </>
  )
})
