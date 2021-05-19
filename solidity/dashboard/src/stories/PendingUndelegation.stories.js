import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import PageWrapper from "../components/PageWrapper"
import PendingUndelegation from "../components/PendingUndelegation";

const mockedEmptyTableData = {
  authorizerAddress: "0x0000000000000000000000000000000000000000",
  beneficiaryAddress: "0x0000000000000000000000000000000000000000",
  error: null,
  ownerAddress: "0x0000000000000000000000000000000000000000",
  slashedTokensError: null,
  stakedBalance: "0",
  delegationStatus: "UNDELEGATED",
}

// storiesOf("PendingUndelegation", module).addDecorator(centered)

export default {
  title: "PendingUndelegation",
  component: PendingUndelegation,
}

const Template = (args) => <PendingUndelegation {...args} />

// TODO WithMockedData
// export const WithMockedData = Template.bind({})
// WithTitle.args = { title: "PageWrapper title" }
