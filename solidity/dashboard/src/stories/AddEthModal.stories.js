import React from "react"
import AddETHModal from "../components/AddETHModal"
import { Provider } from "react-redux"
import store from "../store"
import { ModalContextProvider } from "../components/Modal"

export default {
  title: "AddEthModal",
  component: AddETHModal,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <ModalContextProvider>
          <Story />
        </ModalContextProvider>
      </Provider>
    ),
  ],
}

const Template = (args) => <AddETHModal {...args} />

export const Primary = Template.bind({})
Primary.args = { operatorAddress: "0x5777C7DdEd294654FbefC1Ed262fC8Ba4Ac40De1" }
