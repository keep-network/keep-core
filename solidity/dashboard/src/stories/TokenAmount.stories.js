import React from "react"
import TokenAmount from "../components/TokenAmount"

export default {
  title: "TokenAmount",
  component: TokenAmount,
}

const Template = (args) => <TokenAmount {...args} />

export const Default = Template.bind({})
Default.args = {
  amount: "3000000000000000000",
  withMetricSuffix: true,
  amountClassName: "text-mint-100",
  symbolClassName: "text-mint-100",
}