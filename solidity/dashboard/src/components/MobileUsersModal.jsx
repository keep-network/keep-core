import React from "react"
import Button from "./Button"
import * as Icons from "./Icons"

const MobileUsersModal = ({ closeModal }) => {
  return (
    <div className={"mobile-users-modal"}>
      <Icons.Dashboard className={"mobile-users-modal__icon"} />
      <h2 className={"mobile-users-modal__main-text"}>
        {"The dashboard shines on desktop."}
      </h2>
      <span className={"mobile-users-modal__secondary-text"}>
        {"Switch to a desktop for the best viewing experience of the dashboard"}
      </span>
      <a
        href="https://keep.network/"
        className="btn btn-lg btn-primary mobile-users-modal__button h2"
        rel="noopener noreferrer"
        target="_blank"
      >
        {"VIEW WEBSITE"}
      </a>
      <Button
        className={"btn btn-lg btn-secondary mobile-users-modal__button"}
        onClick={closeModal}
      >
        {"VIEW DASHBOARD"}
      </Button>
    </div>
  )
}

export default MobileUsersModal
