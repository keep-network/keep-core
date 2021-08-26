import { Message } from "../Message"
import React from "react"
import * as Icons from "../Icons"
import NavLink from "../NavLink"

const CovPoolsClaimTokensTimeLeftMessage = ({
  sticky,
  title = "",
  linkText = "Claim your tokens",
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.Warning
  const content = <NavLink to={"/coverage-pools/deposit"}>{linkText}</NavLink>

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

export default React.memo(CovPoolsClaimTokensTimeLeftMessage)
