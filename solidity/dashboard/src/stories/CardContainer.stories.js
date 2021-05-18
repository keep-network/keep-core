import React from "react"
import CardContainer from "../components/CardContainer"
import { Default as DefaultCardStory } from "../stories/Card.stories"

export default {
  title: "CardContainer",
  component: CardContainer,
}

const TemplateEvenNumber = (args) => (
  <CardContainer {...args}>
    {[...Array(4)].map((e, i) => {
      return <DefaultCardStory {...DefaultCardStory.args} key={i} />
    })}
  </CardContainer>
)

const TemplateOddNumber = (args) => (
  <CardContainer {...args}>
    {[...Array(5)].map((e, i) => {
      return <DefaultCardStory {...DefaultCardStory.args} key={i} />
    })}
  </CardContainer>
)

export const WithEvenNumberOfCards = TemplateEvenNumber.bind({})
TemplateEvenNumber.args = {}

export const WithOddNumberOfCards = TemplateOddNumber.bind({})
TemplateOddNumber.args = {}
