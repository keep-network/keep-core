import React from "react"
import { withBaseModal } from "../withBaseModal"
import { ModalBody, ModalFooter, ModalHeader } from "../Modal"
import Button from "../../Button"

export const AuthorizedButNotStakedToTWarning = withBaseModal(({ onClose }) => {
  return (
    <>
      <ModalHeader>Take note</ModalHeader>
      <ModalBody>
        <h3>This stake is not yet staked on Threshold</h3>
        The stake amount is not yet confirmed. This stake is not staked on
        Threshold until it is confirmed. On the Applications page, click on
        “Stake” to initiate the stake confirmation transaction.
      </ModalBody>
      <ModalFooter>
        <Button className={"btn btn-secondary btn-lg"} onClick={onClose}>
          close
        </Button>
      </ModalFooter>
    </>
  )
})
