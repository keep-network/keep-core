import React from "react"
import Undelegations from "../components/Undelegations"
import moment from "moment";
import centered from "@storybook/addon-centered/react"

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
  {
    undelegatedAt: "1621338570",
    amount: "20000000000000000000000000",
    beneficiary: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
    operatorAddress: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
    createdAt: "1621235276",
    authorizerAddress: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
    managedGrantContractInstance: null,
    undelegationCompleteAt: moment("2021-06-01T11:49:30.000Z"),
    canRecoverStake: false,
  },
]

export default {
  title: "Undelegations",
  component: Undelegations,
  decorators: [centered]
}

const Template = (args) => <Undelegations {...args} />

export const EmptyTable = Template.bind({})
EmptyTable.args = {
  undelegations: [],
}

export const WithMockedData = Template.bind({})
WithMockedData.args = {
  undelegations: mockedUndelegations,
}
