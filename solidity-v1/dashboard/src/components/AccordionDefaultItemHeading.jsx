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
      className={`accordion__heading--with-plus-minus-sign ${className}`}
    >
      <AccordionItemButton {...accordionItemButtonProps}>
        <h3>{children}</h3>
        <div>
          <AccordionItemState>
            {({ expanded }) => (expanded ? "-" : "+")}
          </AccordionItemState>
        </div>
      </AccordionItemButton>
    </AccordionItemHeading>
  )
}

export default AccordionDefaultItemHeading
