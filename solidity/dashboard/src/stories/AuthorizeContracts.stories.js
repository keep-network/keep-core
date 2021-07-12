import React from "react"
import AuthorizeContracts from "../components/AuthorizeContracts"
import centered from "@storybook/addon-centered/react"

const mockedData = {
  operatorAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
  stakeAmount: "20000000000000000000000000",
  contracts: [
    {
      contractName: "BondedECDSAKeepFactory",
      operatorContractAddress: "0x3CA39f71A5977F4A4386e4e72687bc7C7eaaecF3",
      isAuthorized: false,
    },
    {
      contractName: "TBTCSystem",
      operatorContractAddress: "0xDD610C207e6bd70D656Bf1C046A6ff8de0720BC0",
      isAuthorized: false,
    },
  ],
}

export default {
  title: "AuthorizeContracts",
  component: AuthorizeContracts,
  argTypes: {
    onSelectOperator: {
      action: "operator selected",
    },
    onAuthorizeBtn: {
      action: "Authorize button clicked",
    },
    onDeauthorizeBtn: {
      action: "Deauthorize button clicked",
    },
  },
  decorators: [centered],
}

const Template = (args) => <AuthorizeContracts {...args} />

export const Empty = Template.bind({})
Empty.args = {
  data: [],
  selectedOperator: {},
  filterDropdownOptions: [],
}

export const WithMockedData = Template.bind({})
WithMockedData.args = {
  data: [mockedData],
  selectedOperator: {},
  filterDropdownOptions: [mockedData],
}
