import React from "react"
import AuthorizeContracts from "../components/AuthorizeContracts"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("AuthorizeContracts", module).addDecorator(centered)

const mockFilterDropdownOptions = [
  {
    operatorAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
    stakeAmount: "20000000000000000000000000",
    contracts: [
      {
        contractName: "Keep Random Beacon Operator Contract",
        operatorContractAddress: "0x160C453639B469aCeC3E47e8c225296d3Fd5fA3b",
        isAuthorized: true,
      },
    ],
  },
]

export default {
  title: "AuthorizeContracts",
  component: AuthorizeContracts,
  argTypes: {
    onSelectOperator: {
      action: "operator selected",
    },
  },
}

const Template = (args) => <AuthorizeContracts {...args} />

export const Empty = Template.bind({})
Empty.args = {
  data: [],
  selectedOperator: {},
  filterDropdownOptions: mockFilterDropdownOptions,
}

// TODO: template with data
