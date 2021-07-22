import React from "react"
import centered from "@storybook/addon-centered/react"
import * as Icons from "../components/Icons"
import DoubleIcon from "../components/DoubleIcon"
import Card from "../components/Card"
import CardContainer from "../components/CardContainer"

export default {
  title: "DoubleIcon",
  component: DoubleIcon,
  decorators: [centered],
}

const Template = (args) => <DoubleIcon {...args} />

export const KEEP_TBTC = Template.bind({})
KEEP_TBTC.args = { MainIcon: Icons.KeepBlackGreen, SecondaryIcon: Icons.TBTC }
KEEP_TBTC.decorators = [
  (Story) => (
    <CardContainer>
      <Card className={`liquidity__card tile keep-tbtc`}>
        <div
          className={"liquidity__card-title"}
          style={{ justifyContent: "center" }}
        >
          <Story />
        </div>
      </Card>
    </CardContainer>
  ),
]

export const FirstIconWithTransparentBackground = Template.bind({})
FirstIconWithTransparentBackground.args = {
  MainIcon: Icons.TBTC,
  SecondaryIcon: Icons.EthToken,
  className: "liquidity__double-icon-container",
}
FirstIconWithTransparentBackground.decorators = [
  (Story) => (
    <CardContainer>
      <Card className={`liquidity__card tile tbtc-eth`}>
        <div
          className={"liquidity__card-title"}
          style={{ justifyContent: "center" }}
        >
          <Story />
        </div>
      </Card>
    </CardContainer>
  ),
]
