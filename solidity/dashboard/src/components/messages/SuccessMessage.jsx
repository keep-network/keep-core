import * as Icons from "../Icons"
import { Message } from "../Message"
import React from "react"

const SuccessMessage = ({
  sticky,
  withTransactionHash,
  txHash,
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.Success
  const title = `Success!`
  const classes = {
    iconClassName: "success-icon green",
  }

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={title}
      classes={classes}
      withTransactionHash={withTransactionHash}
      txHash={txHash}
      messageId={messageId}
      onMessageClose={onMessageClose}
    />
  )
}

export default React.memo(SuccessMessage)
