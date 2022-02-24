import React from "react"
import { ModalBody, ModalFooter, ModalHeader } from "../Modal"
import * as Icons from "../../Icons"
import { colors } from "../../../constants/colors"
import { ViewInBlockExplorer } from "../../ViewInBlockExplorer"
import Button from "../../Button"
import OnlyIf from "../../OnlyIf"

const ViewTransactionSuccessModal = ({
  modalHeader,
  furtherDescription = "",
  txHash,
  onClose,
  renderAdditionalButtons = null,
}) => {
  return (
    <>
      <ModalHeader>{modalHeader}</ModalHeader>
      <ModalBody>
        <h3 className="flex row center mb-1">
          <Icons.OK color={colors.success} />
          <span className="ml-1">Success!</span>
        </h3>
        <h4 className="text-grey-70 mb-2">
          View your transaction&nbsp;
          <ViewInBlockExplorer
            type="tx"
            className="text-grey-70"
            id={txHash}
            text="here"
          />
          . {furtherDescription}
        </h4>
      </ModalBody>
      <ModalFooter>
        <OnlyIf condition={renderAdditionalButtons}>
          {renderAdditionalButtons}
        </OnlyIf>
        <Button className="btn btn-secondary btn-lg" onClick={onClose}>
          close
        </Button>
      </ModalFooter>
    </>
  )
}

export default ViewTransactionSuccessModal
