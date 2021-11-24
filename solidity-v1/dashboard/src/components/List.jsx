import React, { useContext } from "react"
import { isEmptyArray } from "../utils/array.utils"

const ListContext = React.createContext({
  items: [],
  renderItem: () => <></>,
})

const useListContext = () => {
  const context = useContext(ListContext)

  if (!context) {
    throw new Error("ListContext used outside of List component")
  }

  return context
}

const renderListItem = (item, index) => (
  <List.DefaultItem key={index} {...item} />
)

const List = ({
  items = [],
  renderItem = renderListItem,
  className = "",
  children,
  ...restProps
}) => {
  return (
    <ListContext.Provider value={{ items, renderItem }}>
      <section className={`list ${className}`} {...restProps}>
        {children}
      </section>
    </ListContext.Provider>
  )
}

List.Title = ({ children, className = "" }) => {
  return <h4 className={`list__title ${className}`}>{children}</h4>
}

const ListContent = ({ children, className = "" }) => {
  const { items, renderItem } = useListContext()

  return (
    <ul className={`list__content ${className}`}>
      {isEmptyArray(items) ? children : items.map(renderItem)}
    </ul>
  )
}

List.Content = ListContent

List.DefaultItem = ({ className = "", icon: IconComponent, label }) => {
  return (
    <List.Item className={className}>
      {IconComponent && (
        <>
          <IconComponent />
          &nbsp;
        </>
      )}
      <span>{label}</span>
    </List.Item>
  )
}

List.Item = ({ children, className = "" }) => {
  return <li className={`list__item ${className}`}>{children}</li>
}

export default List
