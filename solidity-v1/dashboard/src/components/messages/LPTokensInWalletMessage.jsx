import React from "react"
import * as Icons from "../Icons"
import { Message } from "../Message"
import NavLink from "../NavLink"

const LPTokensInWalletMessage = ({
  liquidityRewardPairName,
  sticky,
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.Wallet
  const title = `[${liquidityRewardPairName}] Your wallet has LP Tokens!`
  const content = (
    <NavLink to={"/liquidity"}>Deposit them and earn rewards</NavLink>
  )
  const classes = {
    bannerDescription: "m-0",
    iconClassName: "wallet-icon grey-50",
  }

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={title}
      content={content}
      classes={classes}
      messageId={messageId}
      onMessageClose={onMessageClose}
    />
  )
}

export default React.memo(LPTokensInWalletMessage)
