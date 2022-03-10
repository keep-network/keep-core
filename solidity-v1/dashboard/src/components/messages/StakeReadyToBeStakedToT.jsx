import { Message } from "../Message"
import React from "react"
import * as Icons from "../Icons"
import NavLink from "../NavLink"

const StakeReadyToBeStakedToT = ({
  sticky,
  numberOfStakes = 0,
  linkText = "Go to Applications →",
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.TTokenSymbol
  const content = <NavLink to={"/applications/threshold"}>{linkText}</NavLink>

  const inflectionOfAStakeWord = numberOfStakes === 1 ? "stake" : "stakes"
  const title = `You have ${numberOfStakes} ${inflectionOfAStakeWord} to stake on Threshold and earn rewards.`

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={title}
      content={content}
      messageId={messageId}
      onMessageClose={onMessageClose}
      classes={{
        iconClassName: "stake-ready-to-be-staked-to-t-message__icon",
      }}
    />
  )
}

export default React.memo(StakeReadyToBeStakedToT)
