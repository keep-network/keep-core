import React from "react"
import { Message } from "../Message"
import * as Icons from "../Icons"

const LiquidityRewardsEarnedMessage = ({
  liquidityRewardPairName,
  sticky,
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.Rewards
  const title = `[${liquidityRewardPairName}] You've earned rewards!`
  const content = <a href={"/liquidity"}>View your balance</a>
  const classes = {
    bannerDescription: "m-0",
    iconClassName: "reward-icon brand-violet",
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

export default React.memo(LiquidityRewardsEarnedMessage)
