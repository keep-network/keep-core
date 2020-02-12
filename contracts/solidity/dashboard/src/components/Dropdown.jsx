import React, { useState } from 'react'

const Dropdown = ({ onSelect, options, valuePropertyName, labelPropertyName, selectedItem, labelPrefix }) => {
  const [isOpen, setIsOpen] = useState(false)


  const renderDropdownItem = (item) => <DropdownItem
    key={item[valuePropertyName]}
    value={item[valuePropertyName]}
    label={item[labelPropertyName]}
    isSelected={item[valuePropertyName == selectedItem[valuePropertyName]]}
    onChange={onChange}
    labelPrefix={labelPrefix}
  />

  const onChange = (e) => {
    console.log('event', e.target.value, options)
    const selectedItem = options.find((option) => option[valuePropertyName] == e.target.value)
    console.log('selectedItem', selectedItem)
    onSelect(selectedItem)
    setIsOpen(false)
  }

  return (
    <div className="select-wrapper">
      <div className={`select${isOpen ? ' open' : ''}`}>
        <div className="select-trigger" onClick={() => setIsOpen(!isOpen)}>
          <span>{selectedItem ? `${labelPrefix} ${selectedItem[labelPropertyName]}` : 'Select grant'}</span>
          <div className="arrow"/>
        </div>
        <ul className="options">
          {options.map(renderDropdownItem)}
        </ul>
      </div>
    </div>
  )
}

const DropdownItem = ({ value, label, labelPrefix, isSelected, onChange }) => {
  return (
    <li className={`option${isSelected ? ' selected' : ''}`} value={value} onClick={onChange}>
      {`${labelPrefix} ${label}`}
    </li>
  )
}

export default Dropdown
