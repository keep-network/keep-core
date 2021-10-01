import * as Icons from "../Icons"
import { Message } from "../Message"
import React from "react"

const ErrorMessage = ({ sticky, content, messageId, onMessageClose }) => {
  const icon = Icons.Warning
  const title = `Error`

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

export default React.memo(ErrorMessage)
