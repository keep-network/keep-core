import React from "react"
import Button from "../Button"
import * as Icons from "../Icons"
import { LINK } from "../../constants/constants"
import {
  Modal,
  ModalOverlay,
  ModalCloseButton,
  ModalContent,
  ModalBody,
} from "./Modal"

const MobileUsersModal = ({ onClose }) => {
  return (
    <Modal isOpen onClose={onClose} size="sm">
      <ModalOverlay />
      <ModalContent>
        <ModalCloseButton />
        <ModalBody>
          <div className={"mobile-users-modal"}>
            <Icons.Dashboard className={"mobile-users-modal__icon"} />
            <div className={"mobile-users-modal__text-container"}>
              <h2 className={"mobile-users-modal__main-text"}>
                {"The dashboard shines on desktop."}
              </h2>
              <span className={"mobile-users-modal__secondary-text"}>
                {
                  "Switch to a desktop for the best viewing experience of the dashboard"
                }
              </span>
            </div>
            <a
              href={LINK.keepWebsite}
              className="btn btn-lg btn-primary mobile-users-modal__button h2"
              rel="noopener noreferrer"
              target="_blank"
            >
              {"VIEW WEBSITE"}
            </a>
            <Button
              className={"btn btn-lg btn-secondary mobile-users-modal__button"}
              onClick={onClose}
            >
              {"VIEW DASHBOARD"}
            </Button>
            <span className={"mobile-users-modal__discord-info text-grey-60"}>
              {"Curious for more? "}
              {
                <a
                  target="_blank"
                  rel="noopener noreferrer"
                  href={LINK.discord}
                  className={`text-link`}
                >
                  Join our Discord
                </a>
              }
            </span>
          </div>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

export default MobileUsersModal
