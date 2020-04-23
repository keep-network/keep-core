import React, { useState } from "react"
import { isEmptyObj } from "../utils/general.utils"

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
    onSelect(selectedItem)
    setIsOpen(false)
  }

  const renderSelectedItem = () => {
    if (options && options.length === 0) {
      return <span className="text-smaller">No items to select</span>
    } else if (isEmptyObj(selectedItem)) {
      return <span>{noItemSelectedText}</span>
    } else if (selectedItemComponent) {
      return selectedItemComponent
    } else {
      return <span>{`${labelPrefix} ${selectedItem[labelPropertyName]}`}</span>
    }
  }

  return (
    <React.Fragment>
      <label className="text-small text-grey-60">{label}</label>
      <div className="select-wrapper">
        <div className={`select${isOpen ? " open" : ""}`}>
          <div className="select-trigger" onClick={() => setIsOpen(!isOpen)}>
            {renderSelectedItem()}
            <div className="arrow" />
          </div>
          <ul className="options">{options.map(renderDropdownItem)}</ul>
        </div>
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
}

export default React.memo(Dropdown, dropdownPropsAreEqual)
