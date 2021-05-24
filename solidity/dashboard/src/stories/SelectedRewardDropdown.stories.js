import React from "react"
import centered from "@storybook/addon-centered/react"
import SelectedRewardDropdown from "../components/SelectedRewardDropdown"

const mockedGroupReward = {
  reward: 20,
  groupPublicKey: "0xeF42ac774dD0d3519E7CBFD59F36e52038D4e255",
}

export default {
  title: "SelectedRewardDropdown",
  component: SelectedRewardDropdown,
  decorators: [
    (Story) => (
      <div style={{ width: "20rem" }}>
        <Story />
      </div>
    ),
    centered,
  ],
}

const Template = (args) => <SelectedRewardDropdown {...args} />

export const Default = Template.bind({})
Default.args = { groupReward: mockedGroupReward }
