import React, { useState } from "react"
import { isEmptyObj } from "../utils/general.utils"
import * as Icons from "./Icons"

const Dropdown = ({
  label,
  onSelect,
  options,
  valuePropertyName,
  labelPropertyName,
  selectedItem,
  labelPrefix,
  noItemSelectedText,
  selectedItemComponent,
  renderOptionComponent,
  withLabel,
  isFilterDropdow,
  allItemsFilterText,
}) => {
  const [isOpen, setIsOpen] = useState(false)

  const renderDropdownItem = (item) => (
    <DropdownItem
      key={item[valuePropertyName]}
      value={item[valuePropertyName]}
      label={item[labelPropertyName]}
      isSelected={item[valuePropertyName] === selectedItem[valuePropertyName]}
      onChange={onChange}
      labelPrefix={labelPrefix}
      item={item}
      renderOptionComponent={renderOptionComponent}
    />
  )

  const onChange = (value) => {
    const selectedItem = options.find(
      (option) => option[valuePropertyName] === value
    )
    onSelect(selectedItem || {})
    setIsOpen(false)
  }

  const renderSelectedItem = () => {
    if (options && options.length === 0 && !isFilterDropdow) {
      return <span className="text-smaller">No items to select</span>
    } else if (
      isFilterDropdow &&
      ((options && options.length === 0) || isEmptyObj(selectedItem))
    ) {
      return (
        <div className="flex row center">
          <Icons.Filter className="mr-1" />
          <span className="text-smaller text-grey-60">
            {allItemsFilterText}
          </span>
        </div>
      )
    } else if (isEmptyObj(selectedItem)) {
      return <span>{noItemSelectedText}</span>
    } else if (isFilterDropdow) {
      return (
        <div className="flex row center">
          <Icons.Filter className="mr-1" />
          <div className="text-smaller text-grey-60">
            {selectedItemComponent
              ? selectedItemComponent
              : `${labelPrefix} ${selectedItem[labelPropertyName]}`}
          </div>
        </div>
      )
    } else if (selectedItemComponent) {
      return selectedItemComponent
    } else {
      return <span>{`${labelPrefix} ${selectedItem[labelPropertyName]}`}</span>
    }
  }

  return (
    <React.Fragment>
      {withLabel && (
        <label className="text-small text-grey-60">
          {label}&nbsp;
          <span className="text-smaller">{options.length} available</span>
        </label>
      )}
      <div
        className={`select${isFilterDropdow ? " filter" : ""}${
          isOpen ? " open" : ""
        }`}
      >
        <div className="select-trigger" onClick={() => setIsOpen(!isOpen)}>
          {renderSelectedItem()}
          <div className="arrow" />
        </div>
        <ul className="options">
          {isFilterDropdow && (
            <li
              className={`option${isEmptyObj(selectedItem) ? " selected" : ""}`}
              onClick={() => onChange({})}
            >
              {allItemsFilterText}
            </li>
          )}
          {options.map(renderDropdownItem)}
        </ul>
      </div>
    </React.Fragment>
  )
}

const DropdownItem = React.memo(
  ({
    value,
    label,
    labelPrefix,
    isSelected,
    onChange,
    renderOptionComponent,
    item,
  }) => {
    return (
      <li
        className={`option${isSelected ? " selected" : ""}`}
        onClick={() => onChange(value)}
      >
        {renderOptionComponent
          ? renderOptionComponent(item)
          : `${labelPrefix} ${label}`}
      </li>
    )
  },
  (prevProps, nextProps) => prevProps.isSelected === nextProps.isSelected
)

const dropdownPropsAreEqual = (prevProps, nextProps) => {
  return (
    prevProps.selectedItem[prevProps.valuePropertyName] ===
      nextProps.selectedItem[prevProps.valuePropertyName] &&
    prevProps.options === nextProps.options
  )
}

Dropdown.defaultProps = {
  noItemSelectedText: "Select Item",
  label: "Select Item",
  withLabel: true,
  isFilterDropdow: false,
  allItemsFilterText: "All items",
}

export default React.memo(Dropdown, dropdownPropsAreEqual)
