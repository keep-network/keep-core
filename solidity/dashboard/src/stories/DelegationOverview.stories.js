import React from "react"
import centered from "@storybook/addon-centered/react"
import DelegationOverview from "../components/DelegationOverview"
import { Provider } from "react-redux"
import store from "../store"
import BigNumber from "bignumber.js";
import moment from "moment";

const mockedDelegations = [
  {
    undelegatedAt: "0",
    amount: "20000000000000000000000000",
    beneficiary: "0xCDAfb5A23A1F1c6f80706Cc101BCcf4b9A1A3e3B",
    operatorAddress: "0xCDAfb5A23A1F1c6f80706Cc101BCcf4b9A1A3e3B",
    createdAt: "1621235278",
    authorizerAddress: "0xCDAfb5A23A1F1c6f80706Cc101BCcf4b9A1A3e3B",
    managedGrantContractInstance: null,
    isInInitializationPeriod: false,
    initializationOverAt: moment("2021-05-17T07:07:57.000Z"),
  },
  {
    undelegatedAt: "0",
    amount: "20000000000000000000000000",
    beneficiary: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
    operatorAddress: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
    createdAt: "1621235277",
    authorizerAddress: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
    managedGrantContractInstance: null,
    isInInitializationPeriod: false,
    initializationOverAt: moment("2021-05-17T07:07:57.000Z"),
  },
  {
    undelegatedAt: "0",
    amount: "20000000000000000000000000",
    beneficiary: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
    operatorAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
    createdAt: "1621235277",
    authorizerAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
    managedGrantContractInstance: null,
    isInInitializationPeriod: false,
    initializationOverAt: moment("2021-05-17T07:07:57.000Z"),
  },
  {
    undelegatedAt: "0",
    amount: "20000000000000000000000000",
    beneficiary: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
    operatorAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
    createdAt: "1621235276",
    authorizerAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
    managedGrantContractInstance: null,
    isInInitializationPeriod: false,
    initializationOverAt: moment("2021-05-17T07:07:57.000Z"),
  },
]

const mockedUndelegations = [
  {
    undelegatedAt: "1621338570",
    amount: "20000000000000000000000000",
    beneficiary: "0xf119557AC33585405467135eC9A343DCDb047517",
    operatorAddress: "0xf119557AC33585405467135eC9A343DCDb047517",
    createdAt: "1621235276",
    authorizerAddress: "0xf119557AC33585405467135eC9A343DCDb047517",
    managedGrantContractInstance: null,
    undelegationCompleteAt: moment("2021-06-01T11:49:30.000Z"),
    canRecoverStake: false,
  },
]

const mockedAvailableTopUps = [
  {
    operatorAddress: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
    availableTopUpAmount: new BigNumber("70000000000000000000000"),
    createdAt: 1621338556,
    readyToBeCommitted: true,
    isInUndelegation: false,
  },
]

export default {
  title: "DelegationOverview",
  component: DelegationOverview,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
    centered,
  ],
}

const Template = (args) => <DelegationOverview {...args} />

export const EmptyData = Template.bind({})
EmptyData.args = {
  delegations: [],
  undelegations: [],
  isFetching: false,
  topUps: [],
  areTopUpsFetching: false,
  undelegationPeriod: "1209600",
  initializationPeriod: "1",
  grants: [],
  selectedGrant: {},
  context: "wallet",
}

export const WithMockedData = Template.bind({})
WithMockedData.args = {
  delegations: mockedDelegations,
  undelegations: mockedUndelegations,
  isFetching: false,
  topUps: mockedAvailableTopUps,
  areTopUpsFetching: false,
  undelegationPeriod: "1209600",
  initializationPeriod: "1",
  keepTokenBalance: "388850139352913431577157930",
  grants: [],
  selectedGrant: {},
  context: "wallet",
}

export const FetchingData = Template.bind({})
FetchingData.args = {
  delegations: [],
  undelegations: [],
  isFetching: true,
  topUps: [],
  areTopUpsFetching: true,
  undelegationPeriod: "1209600",
  initializationPeriod: "1",
  grants: [],
  selectedGrant: {},
  context: "wallet",
}
