import * as Icons from "../Icons"
import { Message } from "../Message"
import React from "react"

const PendingActionMessage = ({
  sticky,
  withTransactionHash,
  txHash,
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.Time
  const title = `Pending transaction`

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={title}
      withTransactionHash={withTransactionHash}
      txHash={txHash}
      messageId={messageId}
      onMessageClose={onMessageClose}
    />
  )
}

export default React.memo(PendingActionMessage)
