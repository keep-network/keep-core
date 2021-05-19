import React from "react"
import { Message } from "../components/Message"
import * as Icons from "../components/Icons"
import Modal, {ModalContextProvider} from "../components/Modal"
import {storiesOf} from "@storybook/react";
import centered from "@storybook/addon-centered/react";
import {Provider} from "react-redux";
import store from "../store"

//TODO: Not sure how to write story for this yet

// storiesOf("Modal", module).addDecorator(centered)
//
// export default {
//   title: "Modal",
//   component: Modal,
//   argTypes: {
//     closeModal: {
//       action: "modal closed",
//     },
//   },
//   decorators: [
//     (Story) => (
//       <Provider store={store}>
//         <ModalContextProvider>
//           <Story />
//         </ModalContextProvider>
//       </Provider>
//     ),
//   ],
// }
//
// const Template = (args) => <Modal {...args} />
//
// export const Default = Template.bind({})
// Default.args = {
//   title: "Modal title",
//   children: <div>Supertest</div>,
//   // isOpen: true,
//   isFullScreen: false,
// }