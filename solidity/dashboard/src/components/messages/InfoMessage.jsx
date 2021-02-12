import * as Icons from "../Icons"
import { Message } from "../Message"
import React from "react"

// TODO: use info message somewhere or delete this component completely
const InfoMessage = ({ sticky, content, messageId, onMessageClose }) => {
  const icon = Icons.MoreInfo
  const title = `Info`

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

export default React.memo(InfoMessage)
