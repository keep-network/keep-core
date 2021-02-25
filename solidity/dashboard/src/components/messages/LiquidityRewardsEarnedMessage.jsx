import React from "react"
import { Message } from "../Message"
import * as Icons from "../Icons"
import { Link } from "react-router-dom"
import { useSelector } from "react-redux"

const LiquidityRewardsEarnedMessage = ({
  sticky,
  messageId,
  messageType,
  onMessageClose,
}) => {
  const {
    liquidityRewardNotification: { pairsDisplayed },
  } = useSelector((state) => state.notificationsData)

  const icon = Icons.Rewards
  const content = <Link to={"/liquidity"}>View your balance</Link>
  const classes = {
    bannerDescription: "m-0",
    iconClassName: "reward-icon brand-violet",
  }

  const formattedTitle = (liquidityRewardPairNames) => {
    const mainText = "You've earned rewards"
    let pairNamesContent = ""
    if (liquidityRewardPairNames.length > 0) {
      pairNamesContent = " for "
      for (const [i, pairName] of liquidityRewardPairNames.entries()) {
        pairNamesContent = pairNamesContent.concat(pairName)
        if (
          i !== liquidityRewardPairNames.length - 1 &&
          i !== liquidityRewardPairNames.length - 2
        ) {
          pairNamesContent = pairNamesContent.concat(", ")
        } else if (i === liquidityRewardPairNames.length - 2) {
          pairNamesContent = pairNamesContent.concat(" and ")
        }
      }
    }
    pairNamesContent = pairNamesContent.concat("!")
    return mainText.concat(pairNamesContent)
  }

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={formattedTitle(pairsDisplayed)}
      content={content}
      classes={classes}
      messageId={messageId}
      messageType={messageType}
      onMessageClose={onMessageClose}
    />
  )
}

export default React.memo(LiquidityRewardsEarnedMessage)
