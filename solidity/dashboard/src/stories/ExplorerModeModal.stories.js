import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import ExplorerModeModal from "../components/ExplorerModeModal"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"

storiesOf("ExplorerModeModal", module).addDecorator(centered)

export default {
  title: "ExplorerModeModal",
  component: ExplorerModeModal,
  argTypes: {
    connectAppWithWallet: {
      action: "connectAppWithWallet function called",
    },
    closeModal: {
      action: "closeModal clicked",
    },
  },
}

const Template = (args) => <ExplorerModeModal {...args} />

export const Default = Template.bind({})
Default.args = {
  address: "0x5777C7DdEd294654FbefC1Ed262fC8Ba4Ac40De1",
  connector: new ExplorerModeConnector(),
}
