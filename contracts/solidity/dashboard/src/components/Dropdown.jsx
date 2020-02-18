import React, { useState } from 'react'

const Dropdown = ({ onSelect, options, valuePropertyName, labelPropertyName, selectedItem, labelPrefix, noItemSelectedText }) => {
  const [isOpen, setIsOpen] = useState(false)

  const renderDropdownItem = (item) =>
    <DropdownItem
      key={item[valuePropertyName]}
      value={item[valuePropertyName]}
      label={item[labelPropertyName]}
      isSelected={item[valuePropertyName] == selectedItem[valuePropertyName]}
      onChange={onChange}
      labelPrefix={labelPrefix}
    />

  const onChange = (e) => {
    const selectedItem = options.find((option) => option[valuePropertyName] == e.target.value)
    onSelect(selectedItem)
    setIsOpen(false)
  }

  return (
    <div className="select-wrapper">
      <div className={`select${isOpen ? ' open' : ''}`}>
        <div className="select-trigger" onClick={() => setIsOpen(!isOpen)}>
          <span>{selectedItem ? `${labelPrefix} ${selectedItem[labelPropertyName]}` : noItemSelectedText}</span>
          <div className="arrow"/>
        </div>
        <ul className="options">
          {options.map(renderDropdownItem)}
        </ul>
      </div>
    </div>
  )
}

const DropdownItem = React.memo(({ value, label, labelPrefix, isSelected, onChange }) => {
  return (
    <li className={`option${isSelected ? ' selected' : ''}`} value={value} onClick={onChange}>
      {`${labelPrefix} ${label}`}
    </li>
  )
}, (prevProps, nextProps) => prevProps.isSelected === nextProps.isSelected)

const dropdownPropsAreEqual = (prevProps, nextProps) => {
  return prevProps.selectedItem[prevProps.valuePropertyName] === nextProps.selectedItem[prevProps.valuePropertyName]
    && prevProps.options === nextProps.options
}

Dropdown.defaultProps = {
  noItemSelectedText: 'Select Item',
}

export default React.memo(Dropdown, dropdownPropsAreEqual)
