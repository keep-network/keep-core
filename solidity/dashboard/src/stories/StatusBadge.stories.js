import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import StatusBadge, { BADGE_STATUS } from "../components/StatusBadge"
import { COMPLETE_STATUS, PENDING_STATUS } from "../constants/constants"

/**
 * StakeDropChart is dropped for now, so we are not displaying story for it
 */

storiesOf("StatusBadge", module).addDecorator(centered)

export default {
  title: "StatusBadge",
  component: StatusBadge,
}

const Template = (args) => <StatusBadge {...args} />

export const Active = Template.bind({})
Active.args = {
  status: BADGE_STATUS.ACTIVE,
  text: "text",
  onlyIcon: false,
}

export const Disable = Template.bind({})
Disable.args = {
  status: BADGE_STATUS.DISABLED,
  text: "text",
  onlyIcon: false,
}

export const Pending = Template.bind({})
Pending.args = {
  status: BADGE_STATUS[PENDING_STATUS],
  text: "text",
  onlyIcon: false,
}

export const Complete = Template.bind({})
Complete.args = {
  status: BADGE_STATUS[COMPLETE_STATUS],
  text: "text",
  onlyIcon: false,
}
