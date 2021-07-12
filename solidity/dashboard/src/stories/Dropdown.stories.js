// TODO: Dropdown story
import centered from "@storybook/addon-centered/react"
import Dropdown from "../components/Dropdown"

export default {
  title: "Dropdown",
  component: Dropdown,
  decorators: [centered],
}

// export const Default = Dropdown.bind({})
// Default.args = {
//   label: "Dropdown label",
//   optrion: ["raz", "dwa", "trzy"],
//   noItemSelectedText: "no item selected",
//   withLabel: true,
// }
