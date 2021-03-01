import * as Icons from "../Icons"
import { Message } from "../Message"
import React from "react"

const WalletMessage = ({ sticky, messageId, onMessageClose }) => {
  const icon = Icons.Wallet
  const title = `Waiting for the transaction confirmation...`
  const classes = {
    iconClassName: "wallet-icon grey-50",
  }

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={title}
      classes={classes}
      messageId={messageId}
      onMessageClose={onMessageClose}
    />
  )
}

export default React.memo(WalletMessage)
