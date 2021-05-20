import React from "react"
import ProgressBar from "../components/ProgressBar"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import TokenAmount from "../components/TokenAmount"
import { colors } from "../constants/colors"

storiesOf("ProgressBar", module).addDecorator(centered)

export default {
  title: "ProgressBar",
  component: ProgressBar,
}

const Template = (args) => <ProgressBar {...args} />

export const Inline = Template.bind({})
Inline.args = {
  value: 10,
  total: 100,
  color: colors.mint80,
  bgColor: colors.mint20,
  children: <ProgressBar.Inline height={8} />,
}

export const Circular = Template.bind({})
Circular.args = {
  value: 10,
  total: 100,
  color: colors.mint80,
  bgColor: colors.mint20,
  children: (
    <div className="circular-progress-bar-percentage-label-wrapper">
      <ProgressBar.Circular radius={82} barWidth={16} />
      <ProgressBar.PercentageLabel text="Progress" />
    </div>
  ),
}

export const CircularWithLegend = Template.bind({})
CircularWithLegend.args = {
  value: 10,
  total: 100,
  color: colors.mint80,
  bgColor: colors.mint20,
  children: (
    <div>
      <div className="circular-progress-bar-percentage-label-wrapper">
        <ProgressBar.Circular radius={82} barWidth={16} />
        <ProgressBar.PercentageLabel text="Unlocked" />
      </div>
      <ProgressBar.Legend
        leftValueLabel="Locked"
        valueLabel="Unlocked"
        renderValuePattern={
          <TokenAmount
            withMetricSuffix
            withSymbol={false}
            amountClassName=""
            symbolClassName=""
          />
        }
      />
    </div>


  )
}
