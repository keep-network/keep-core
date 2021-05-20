import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import StepNav from "../components/StepNav"

const mockedSteps = ["step1", "step2", "step3", "step4"]

storiesOf("StepNav", module).addDecorator(centered)

export default {
  title: "StepNav",
  component: StepNav,
}

const Template = (args) => <StepNav {...args} />

export const Default = Template.bind({})
Default.args = {
  steps: mockedSteps,
  activeStep: 1,
}
