import React from "react"
import SlashedTokensList from "../components/SlashedTokensList"
import centered from "@storybook/addon-centered/react"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

const mockedSlashedTokens = [
  {
    amount: "300000000000000000000",
    event: null,
    groupPublicKey: "0xeF42ac774dD0d3519E7CBFD59F36e52038D4e255",
    date: null,
  },
  {
    amount: "200000000000000000000",
    event: "UnauthorizedSigningReported",
    groupPublicKey: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
    date: null,
  },
]

export default {
  title: "SlashedTokensList",
  component: SlashedTokensList,
  decorators: [
    (Story) => (
      <div style={{ width: "50rem" }}>
        <Story />
      </div>
    ),
    whiteBackground,
    centered,
  ],
}

const Template = (args) => <SlashedTokensList {...args} />

export const EmptyTable = Template.bind({})
EmptyTable.args = {}

export const WithMockedData = Template.bind({})
WithMockedData.args = {
  slashedTokens: mockedSlashedTokens,
}
