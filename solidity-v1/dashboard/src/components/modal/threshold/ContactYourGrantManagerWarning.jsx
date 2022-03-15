import React from "react"
import { withBaseModal } from "../withBaseModal"
import { ModalBody, ModalFooter, ModalHeader } from "../Modal"
import Button from "../../Button"
import { colors } from "../../../constants/colors"
import { GRANT_MANAGER_EMAIL } from "../../../constants/constants"

const styles = {
  header: {
    backgroundColor: colors.yellow30,
    color: colors.yellowSecondary,
    borderBottom: `1px solid ${colors.yellowSecondary}`,
  },
}

export const ContactYourGrantManagerWarning = withBaseModal(
  ({
    header = "Contact your grant manager",
    bodyTitle = "To enable your stake in Threshold you have to contact your Grant Manager",

    onClose,
  }) => {
    return (
      <>
        <ModalHeader style={styles.header}>{header}</ModalHeader>
        <ModalBody>
          <h3>{bodyTitle}</h3>
          <p className="text-grerey-70 mt-2">
            To use this stake in Threshold, you will have to contact your grant
            manager&nbsp;
            <span className="text-secondary">(</span>
            <a href={`mailto:${GRANT_MANAGER_EMAIL}`}>{GRANT_MANAGER_EMAIL}</a>
            <span className="text-secondary">)</span>.
          </p>
          <p className="mt-2">Please do it immediately!</p>
        </ModalBody>
        <ModalFooter>
          <Button className="btn btn-unstyled text-link" onClick={onClose}>
            Cancel
          </Button>
        </ModalFooter>
      </>
    )
  }
)
