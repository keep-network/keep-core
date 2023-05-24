import React from "react"
import List from "../List"
import * as Icons from "../Icons"
import {
  Modal,
  ModalOverlay,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
} from "./Modal"
import { LINK } from "../../constants/constants"
import Button from "../Button"
import { colors } from "../../constants/colors"
import { useShouldShowLegacyDappModal } from "../../hooks/useShowLegacyDappModal"

const styles = {
  header: {
    background: colors.yellow30,
    color: colors.yellowSecondary,
    borderBottomColor: "inherit",
  },
  body: {
    illustration: {
      marginTop: "1.5rem",
    },
    title: {
      marginTop: "2rem",
      marginBottom: "1.5rem",
    },
  },
  footer: {
    link: {
      textDecoration: "underline !important",
    },
  },
}

export const LegacyDashboardModal = ({ isOpen, onClose, size }) => {
  const { modalHasBeenClosed } = useShouldShowLegacyDappModal()

  const _onClose = () => {
    onClose()
    modalHasBeenClosed()
  }

  return (
    <Modal isOpen={isOpen} onClose={_onClose} size={size}>
      <ModalOverlay />
      <ModalContent>
        <ModalCloseButton />
        <ModalHeader style={styles.header}>
          Take note! This is a Legacy Dashboard
        </ModalHeader>
        <ModalBody>
          <Icons.LegacyDappIllustration style={styles.body.illustration} />
          <h3 style={styles.body.title}>
            The Keep Network Dashboard is in legacy mode.
          </h3>
          <div className="text-big text-grey-70 mb-1">
            <p className="mb-1">
              Only the following features are supported in legacy mode:
            </p>
            <List className="ml-3">
              <List.Item>Stake delegation/undelegation</List.Item>
              <List.Item>Token grant withdrawal</List.Item>
              <List.Item>Stake upgrade for Threshold Network</List.Item>
            </List>
            <p className="mb-0 mt-1">
              Visit{" "}
              <a
                href={LINK.thresholdDapp}
                rel="noopener noreferrer"
                target="_blank"
              >
                Threshold Dashboard
              </a>{" "}
              for the tBTC v2 dApp and staking in Threshold Network.
            </p>
          </div>
        </ModalBody>
        <ModalFooter>
          <Button className="btn btn-primary btn-lg mr-2" onClick={_onClose}>
            stay on page
          </Button>
          <a
            href={LINK.thresholdDapp}
            rel="noopener noreferrer"
            target="_blank"
            className="no-arrow text-grey-70"
          >
            Go to Threshold dashboard
          </a>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}
