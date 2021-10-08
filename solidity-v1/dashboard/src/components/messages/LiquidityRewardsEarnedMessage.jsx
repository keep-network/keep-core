import React from "react"
import { Message } from "../Message"
import * as Icons from "../Icons"
import { Link } from "react-router-dom"
import { useDispatch, useSelector } from "react-redux"

const LiquidityRewardsEarnedMessage = ({
  sticky,
  messageId,
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
  const dispatch = useDispatch()

  const formattedTitle = (liquidityRewardPairNames) => {
    const mainText = "You've earned rewards"
    let pairNamesContent = ""
    if (liquidityRewardPairNames.length > 0) {
      pairNamesContent = " for "
      for (const [i, pairName] of liquidityRewardPairNames.entries()) {
        pairNamesContent = pairNamesContent.concat(pairName.replace("_", "+"))
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

  /**
   * Clears the pairs saved in the redux store.
   *
   * This function should be called whenever notification is being closed
   */
  const clearNotificationData = () => {
    dispatch({
      type: "notifications_data/liquidityRewardNotification/pairs_displayed_updated",
      payload: [],
    })
  }

  return (
    <Message
      sticky={sticky}
      icon={icon}
      title={formattedTitle(pairsDisplayed)}
      content={content}
      classes={classes}
      messageId={messageId}
      onMessageClose={(messageId) => {
        onMessageClose(messageId)
        setTimeout(() => {
          clearNotificationData()
        }, 300)
      }}
    />
  )
}

export default React.memo(LiquidityRewardsEarnedMessage)
