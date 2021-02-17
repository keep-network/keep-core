import React from "react"
import { Message } from "../Message"
import * as Icons from "../Icons"
import { Link } from "react-router-dom"

const LiquidityRewardsEarnedMessage = ({
  liquidityRewardPairName,
  sticky,
  messageId,
  onMessageClose,
}) => {
  const icon = Icons.Rewards
  const title = `[${liquidityRewardPairName}] You've earned rewards!`
  const content = <Link to={"/liquidity"}>View your balance</Link>
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
