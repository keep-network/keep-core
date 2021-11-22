import React from "react"
import { ModalHeader, ModalBody, ModalFooter } from "../../Modal"
import { InfoList } from "./components"
import * as Icons from "../../../Icons"
import { ViewInBlockExplorer } from "../../../ViewInBlockExplorer"
import Button from "../../../Button"
import { withBaseModal } from "../../withBaseModal"
import { colors } from "../../../../constants/colors"

const UndelegationInitiatedComponent = ({
  txHash,
  undelegationPeriod,
  undelegatedAt,
  onClose,
}) => {
  return (
    <>
      <ModalHeader>Undelegate</ModalHeader>
      <ModalBody>
        <h3 className="flex row center">
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
          .
        </h4>
        <InfoList
          undelegationPeriod={undelegationPeriod}
          undelegatedAt={undelegatedAt}
        />
      </ModalBody>
      <ModalFooter>
        <Button className="btn btn-secondary btn-lg" onClick={onClose}>
          close
        </Button>
      </ModalFooter>
    </>
  )
}

const UndelegationInitiatedWithBaseModal = withBaseModal(
  UndelegationInitiatedComponent
)

export const UndelegationInitiated = (props) => (
  <UndelegationInitiatedWithBaseModal size="sm" {...props} />
)
