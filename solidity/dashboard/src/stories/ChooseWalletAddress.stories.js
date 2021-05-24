import React from "react"
import centered from "@storybook/addon-centered/react"
import ChooseWalletAddress from "../components/ChooseWalletAddress"

export default {
  title: "ChooseWalletAddress",
  component: ChooseWalletAddress,
  argTypes: {
    onSelectAccount: {
      action: "onSelectAccount clicked",
    },
    onNext: {
      action: "next clicked",
    },
    onPrev: {
      action: "prev clicked",
    },
  },
  decorators: [centered],
}

const Template = (args) => <ChooseWalletAddress {...args} />

export const Default = Template.bind({})
Default.args = {
  addresses: [
    "0xf978F05003a5bb9A8BE6F18102F0070bf7c67b1f",
    "0x100a6c90b1927df586501001201f3390e2C69Efe",
    "0x065993c332b02ab8674Ac033CaCDBccBe7bc9047",
    "0xD827401B2343E21505639125e5D7e9207CC3cF94",
    "0x058630E50d87466843b28954cD23889c6D09F667",
    "0x5777C7DdEd294654FbefC1Ed262fC8Ba4Ac40De1",
  ],
  withPagination: true,
  renderPrevBtn: true,
}
