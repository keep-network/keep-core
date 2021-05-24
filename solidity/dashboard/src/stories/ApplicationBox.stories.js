import React from "react"
import centered from "@storybook/addon-centered/react"
import ApplicationBox from "../components/ApplicationBox"
import { BrowserRouter } from "react-router-dom"
import Web3ContextProvider from "../components/Web3ContextProvider"

export default {
  title: "ApplicationBox",
  component: ApplicationBox,
  argTypes: {
    onClick: {
      action: "onClick clicked",
    },
  },
  decorators: [
    centered,
    (Story) => (
      <BrowserRouter>
        <Web3ContextProvider>
          <Story />
        </Web3ContextProvider>
      </BrowserRouter>
    ),
  ],
}

const Template = (args) => <ApplicationBox {...args} />

export const Default = Template.bind({})
Default.args = {
  // icon: Icons.KeepBlackGreen,
  name: "Test",
  websiteUrl: "https://google.com",
  websiteName: "Google",
  description: "google website",
  btnLink: "/liquidity",
}
