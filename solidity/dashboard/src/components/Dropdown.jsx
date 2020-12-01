import React, { useState, useCallback } from "react"
import { isEmptyObj } from "../utils/general.utils"
import * as Icons from "./Icons"
import { useContext } from "react"

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

const DropdownContext = React.createContext({
  selectedItem: {},
  onSelect: () => {},
  comparePropertyName: null,
})

const useDropdownContext = () => {
  const dropdownContext = useContext(DropdownContext)

  if (!dropdownContext) {
    throw new Error("DropdownContext used outside of the Dropdown component.")
  }
  return dropdownContext
}

const CompoundDropdown = ({
  selectedItem,
  onSelect,
  comparePropertyName,
  children,
  className = "",
  rounded = false,
}) => {
  const [isOpen, setIsOpen] = useState(false)

  const onTriggerDropdown = useCallback(() => {
    setIsOpen((isOpen) => !isOpen)
  }, [setIsOpen])

  return (
    <DropdownContext.Provider
      value={{
        selectedItem,
        onSelect,
        comparePropertyName,
        onTriggerDropdown,
        isOpen,
      }}
    >
      <div
        className={`dropdown${isOpen ? "--open" : "--closed"} ${
          rounded ? "dropdown--rounded" : ""
        } ${className}`}
      >
        {children}
      </div>
    </DropdownContext.Provider>
  )
}

const DropdownOptions = ({
  children,
  selectedIcon = (
    <Icons.Success className="success-icon--black" width={12} height={12} />
  ),
}) => {
  return (
    <ul className="dropdown__options">
      {children.map((child) => React.cloneElement(child, { selectedIcon }))}
    </ul>
  )
}

const DropdownOption = ({ value, selectedIcon, children }) => {
  const { onSelect, selectedItem, comparePropertyName } = useDropdownContext()

  const onClick = useCallback(() => {
    onSelect(value)
  }, [value, onSelect])

  let isSelected
  if (!comparePropertyName) {
    isSelected = selectedItem === value
  } else {
    isSelected =
      !isEmptyObj(selectedItem) &&
      selectedItem[comparePropertyName] === value[comparePropertyName]
  }

  return (
    <li
      className={`dropdown__option${isSelected ? "--selected" : ""}`}
      onClick={onClick}
    >
      {children}
      {isSelected && selectedIcon && (
        <div className="dropdown__option__selected-icon">{selectedIcon}</div>
      )}
    </li>
  )
}

const DropdownTrigger = ({ children }) => {
  const { onTriggerDropdown } = useDropdownContext()
  return (
    <div className="dropdown__trigger" onClick={onTriggerDropdown}>
      <div className="dropdown__trigger__content">{children}</div>
      <div className="dropdown__trigger__arrow" />
    </div>
  )
}

CompoundDropdown.Trigger = DropdownTrigger
CompoundDropdown.Options = DropdownOptions
CompoundDropdown.Option = DropdownOption

export { CompoundDropdown }
