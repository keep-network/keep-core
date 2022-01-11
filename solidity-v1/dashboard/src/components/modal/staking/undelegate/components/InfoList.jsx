import React from "react"
import moment from "moment"
import List from "../../../../List"
import AddToCalendar from "../../../../AddToCalendar"
import { UNDELEGATE_STAKE_CALENDAR_EVENT } from "../../../../../constants/constants"

export const InfoList = ({ undelegationPeriod, undelegatedAt }) => {
  const undelegationCompletedAt = moment
    .unix(undelegatedAt)
    .add(undelegationPeriod, "seconds")
  const undelegationPeriodInRelativeTime = undelegationCompletedAt.fromNow(true)

  return (
    <List>
      <List.Content className="bullets bullets--violet-80 text-grey-60">
        <List.Item>
          Tokens will be available in {undelegationPeriodInRelativeTime}, on{" "}
          {undelegationCompletedAt.format("D MMM YYYY")}.
          {/* TODO: Add `add to calendar` button */}
          <AddToCalendar
            {...UNDELEGATE_STAKE_CALENDAR_EVENT}
            startsAt={undelegationCompletedAt.unix()}
            endsAt={undelegationCompletedAt.add(15, "minutes").unix()}
            className={"bullets__add-to-calendar"}
          />
        </List.Item>
        <List.Item>
          Withdraw your tokens and upgrade your KEEP to T using the portal on
          the Threshold dapp.
        </List.Item>
      </List.Content>
    </List>
  )
}
