import React from "react"
import AuthorizationHistory from "../components/AuthorizationHistory"
import centered from "@storybook/addon-centered/react"

const mockContracts = [
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
    contractName: "Keep Random Beacon Operator Contract",
    operatorContractAddress: "0x160C453639B469aCeC3E47e8c225296d3Fd5fA3b",
    isAuthorized: true,
  },
]

export default {
  title: "AuthorizationHistory",
  component: AuthorizationHistory,
  decorators: [centered],
}

const Template = (args) => <AuthorizationHistory {...args} />

export const Empty = Template.bind({})
Empty.args = { contracts: [] }

export const WithRows = Template.bind({})
WithRows.args = { contracts: mockContracts }
