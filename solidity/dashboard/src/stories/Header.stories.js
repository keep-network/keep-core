import React from "react"
import Header from "../components/Header"

const mockedSublinks = [
  {
    title: "SubLink1",
    path: "/sublink1",
    exact: true,
  },
  {
    title: "SubLink2",
    path: "/sublink2",
    exact: true,
  },
]

export default {
  title: "Header",
  component: Header,
}

const Template = (args) => <Header {...args} />

export const Default = Template.bind({})
Default.args = {}

export const WithTitle = Template.bind({})
WithTitle.args = { title: "Header title" }

export const WithSublinks = Template.bind({})
WithSublinks.args = { title: "Header title", subLinks: mockedSublinks }

export const IsNewPage = Template.bind({})
IsNewPage.args = { title: "New Page", subLinks: mockedSublinks, newPage: true }
