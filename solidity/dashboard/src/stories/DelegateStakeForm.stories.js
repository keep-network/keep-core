import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import DelegateStakeForm from "../components/DelegateStakeForm"

storiesOf("DelegateStakeForm", module).addDecorator(centered)

export default {
  title: "DelegateStakeForm",
  component: DelegateStakeForm,
  decorators: [
    (Story) => (
      <section className="tile granted-page__overview__stake-form">
        <Story />
      </section>
    ),
  ],
  argTypes: {
    onSubmit: {
      action: "onSubmit clicked",
    },
  },
}

const Template = (args) => <DelegateStakeForm {...args} />

export const Default = Template.bind({})
Default.args = {
  minStake: "50000000000000000000000",
  availableToStake: "827900278705826863154315859",
}
