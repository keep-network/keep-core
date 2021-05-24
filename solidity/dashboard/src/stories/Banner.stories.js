import centered from "@storybook/addon-centered/react"
import Banner from "../components/Banner"
import * as Icons from "../components/Icons"

export default {
  title: "Banner",
  component: Banner,
  decorators: [centered],
}

export const Default = Banner.bind({})
Default.args = {
  inline: true,
  title: "banner title",
  icon: Icons.Time,
}
