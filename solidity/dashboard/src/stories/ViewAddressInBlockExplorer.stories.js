import React from "react"
import { ViewAddressInBlockExplorer } from "../components/ViewInBlockExplorer"

export default {
  title: "ViewInViewAddressInBlockExplorerBlockExplorer",
  component: ViewAddressInBlockExplorer,
}

const Template = (args) => <ViewAddressInBlockExplorer {...args} />

export const Default = Template.bind({})
Default.args = {
  address: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
  text: "View address in block explorer",
  urlSuffix: "#code",
}
