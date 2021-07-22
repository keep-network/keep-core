// import React from "react"
import centered from "@storybook/addon-centered/react"
import PendingUndelegation from "../components/PendingUndelegation"

// const mockedEmptyTableData = {
//   authorizerAddress: "0x0000000000000000000000000000000000000000",
//   beneficiaryAddress: "0x0000000000000000000000000000000000000000",
//   error: null,
//   ownerAddress: "0x0000000000000000000000000000000000000000",
//   slashedTokensError: null,
//   stakedBalance: "0",
//   delegationStatus: "UNDELEGATED",
// }

// storiesOf("PendingUndelegation", module).addDecorator(centered)

export default {
  title: "PendingUndelegation",
  component: PendingUndelegation,
  decorators: [centered],
}

// const Template = (args) => <PendingUndelegation {...args} />

// TODO WithMockedData
// export const WithMockedData = Template.bind({})
// WithTitle.args = { title: "PageWrapper title" }
