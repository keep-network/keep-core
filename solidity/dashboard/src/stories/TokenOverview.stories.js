import React from "react"
import Web3ContextProvider from "../components/Web3ContextProvider"
import { BrowserRouter } from "react-router-dom"
import TokenOverview from "../components/TokenOverview"

export default {
  title: "TokenOverview",
  component: TokenOverview,
  decorators: [
    (Story) => (
      <BrowserRouter>
        <Web3ContextProvider>
          <Story />
        </Web3ContextProvider>
      </BrowserRouter>
    ),
  ],
}

const Template = (args) => <TokenOverview {...args} />

export const Default = Template.bind({})
Default.args = {
  totalKeepTokenBalance: "927900278705826863154315859",
  totalGrantedTokenBalance: "0",
  totalGrantedStakedBalance: "0",
  totalOwnedStakedBalance: "100000000000000000000000000",
  isFetching: false,
}

export const IsFetching = Template.bind({})
IsFetching.args = {
  totalKeepTokenBalance: "927900278705826863154315859",
  totalGrantedTokenBalance: "0",
  totalGrantedStakedBalance: "0",
  totalOwnedStakedBalance: "100000000000000000000000000",
  isFetching: true,
}
