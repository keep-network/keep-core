import React from "react"
import * as Icons from "../Icons"
import { Message } from "../Message"
import NavLink from "../NavLink"

const TopUpReadyToBeCommittedMessage = ({
  sticky,
  messageId,
  onMessageClose,
  grantId = null,
}) => {
  const icon = Icons.KeepGreenOutline
  const title = "Top-up is ready to be committed"
  const path = grantId
    ? {
        pathname: "/delegations/granted",
        hash: `${grantId}`,
      }
    : "/delegations/wallet"
  const content = <NavLink to={path}>Commit top up</NavLink>

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={title}
      content={content}
      messageId={messageId}
      onMessageClose={onMessageClose}
    />
  )
}

export default React.memo(TopUpReadyToBeCommittedMessage)
