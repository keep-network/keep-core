import React from "react"
import Button from "./Button"
import * as Icons from "./Icons"

const MobileUsersModal = ({ closeModal }) => {
  const renderDiscordLink = () => {
    return (
      <a
        target="_blank"
        rel="noopener noreferrer"
        href={"https://discord.gg/wYezN7v"}
        className={`text-link`}
      >
        Join our Discord
      </a>
    )
  }

  return (
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
      <span className={"mobile-users-modal__discord-info text-grey-60"}>
        {"Curious for more? "}
        {renderDiscordLink()}
      </span>
    </div>
  )
}

export default MobileUsersModal
