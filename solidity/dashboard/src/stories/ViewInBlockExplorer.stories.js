import React from "react"
import { ViewInBlockExplorer } from "../components/ViewInBlockExplorer"

export default {
  title: "ViewInBlockExplorer",
  component: ViewInBlockExplorer,
}

const Template = (args) => <ViewInBlockExplorer {...args} />

export const Default = Template.bind({})
Default.args = {
  type: "address",
  id: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
  hashParam: "",
  text: "View in Block Explorer",
}
