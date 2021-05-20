import React from "react"
import { Message } from "../components/Message"
import * as Icons from "../components/Icons"

export default {
  title: "Message",
  component: Message,
  argTypes: {
    onMessageClose: {
      action: "onMessageClose function called",
    },
  },
  decorators: [
    (Story) => (
      <div style={{ width: "25rem" }}>
        <Story />
      </div>
    ),
  ],
}

const Template = (args) => <Message {...args} />

export const Default = Template.bind({})
Default.args = {
  icon: Icons.KeepBlackGreen,
  sticky: true,
  title: "Message title",
  content: "Message content",
  withTransactionHash: false,
  messageId: 1,
}
