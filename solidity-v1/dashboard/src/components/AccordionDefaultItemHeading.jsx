import React from "react"
import {
  AccordionItemHeading,
  AccordionItemButton,
  AccordionItemState,
} from "react-accessible-accordion"
import * as Icons from "../components/Icons"

const AccordionDefaultItemHeading = ({
  children,
  className = "",
  ...accordionItemButtonProps
}) => {
  return (
    <AccordionItemHeading
      className={`accordion__heading accordion__heading--default ${className}`}
    >
      <AccordionItemState>
        {({ expanded }) => {
          return (
            <AccordionItemButton
              {...accordionItemButtonProps}
              className={`accordion__button ${
                expanded ? "accordion__button--expanded" : ""
              }`}
            >
              {children}
              {expanded ? <Icons.Subtract /> : <Icons.Add />}
            </AccordionItemButton>
          )
        }}
      </AccordionItemState>
    </AccordionItemHeading>
  )
}

export default AccordionDefaultItemHeading
