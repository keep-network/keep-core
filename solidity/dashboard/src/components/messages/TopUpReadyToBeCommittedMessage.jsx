import React from "react"
import { Link } from "react-router-dom"
import * as Icons from "../Icons"
import { Message } from "../Message"

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
  const content = <Link to={path}>Commit top up</Link>

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
