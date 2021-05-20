import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import {SeeAllButton} from "../components/SeeAllButton";
import SelectedRewardDropdown from "../components/SelectedRewardDropdown"

const mockedGroupReward = {
  reward: 20,
  groupPublicKey: "0xeF42ac774dD0d3519E7CBFD59F36e52038D4e255",
}

// storiesOf("SelectedRewardDropdown", module).addDecorator(centered)

export default {
  title: "SelectedRewardDropdown",
  component: SelectedRewardDropdown,
  decorators: [
    (Story) => (
      <div style={{ width: "20rem" }}>
        <Story />
      </div>
    ),
  ],
}

const Template = (args) => <SelectedRewardDropdown {...args} />

export const Default = Template.bind({})
Default.args = { groupReward: mockedGroupReward }
