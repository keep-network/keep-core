import React from "react"
import centered from "@storybook/addon-centered/react"
import CreateTokenGrantForm from "../components/CreateTokenGrantForm"
import Tile from "../components/Tile"

export default {
  title: "CreateTokenGrantForm",
  component: CreateTokenGrantForm,
  argTypes: {
    submitAction: {
      action: "onSUbmitAction clicked",
    },
  },
  decorators: [
    (Story) => (
      <Tile title="Create Grant" className="rewards-history tile flex column">
        <Story />
      </Tile>
    ),
    centered,
  ],
}

const Template = (args) => <CreateTokenGrantForm {...args} />

export const Default = Template.bind({})
Default.args = {
  keepBalance: "20000000000000000000000000",
}
