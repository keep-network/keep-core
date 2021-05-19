import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import PageWrapper from "../components/PageWrapper"
import PendingUndelegation from "../components/PendingUndelegation";
import RandomBeaconRewardsTable from "../components/RandomBeaconRewardsTable";

// storiesOf("PendingUndelegation", module).addDecorator(centered)

export default {
  title: "RandomBeaconRewardsTable",
  component: RandomBeaconRewardsTable,
}

const Template = (args) => <RandomBeaconRewardsTable {...args} />

// TODO WithMockedData
// export const WithMockedData = Template.bind({})
// WithTitle.args = { title: "PageWrapper title" }
