import React from "react"
import {
  AccordionItemHeading,
  AccordionItemButton,
  AccordionItemState,
} from "react-accessible-accordion"

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
              <h3>{children}</h3>
              <div>{expanded ? "-" : "+"}</div>
            </AccordionItemButton>
          )
        }}
      </AccordionItemState>
    </AccordionItemHeading>
  )
}

export default AccordionDefaultItemHeading
