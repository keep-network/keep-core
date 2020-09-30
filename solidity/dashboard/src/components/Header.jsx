import React, { useState, useCallback } from "react"
import { NavLink } from "react-router-dom"
import { isEmptyArray } from "../utils/array.utils"
import { useContext } from "react"

const Header = () => {
  const { title, subLinks } = useHeaderContext()

  return (
    <header className="header">
      <div className="header__content">
        <h1 className="header__title">{title}</h1>
        {/* TODO add a web3 status component */}
      </div>
      {!isEmptyArray(subLinks) && (
        <nav className="header__sub-nav">
          <ul className="sub-nav__list">{subLinks.map(renderSubNavItem)}</ul>
        </nav>
      )}
    </header>
  )
}

const SubNavItem = ({ title, path }) => {
  return (
    <li className="sub-nav__item-wrapper">
      <NavLink
        to={path}
        className="sub-nav__item"
        activeClassName="sub-nav__item--active"
        exact={true}
      >
        {title}
      </NavLink>
    </li>
  )
}

const renderSubNavItem = (item, index) => (
  <SubNavItem key={`${index}-${item.path}`} {...item} />
)

const HeaderContext = React.createContext({
  title: "",
  subLinks: [],
  updateHeaderData: () => {},
})

const HeaderContextProvider = React.memo((props) => {
  const [title, setTitle] = useState("")
  const [subLinks, setSubLinks] = useState([])
  const updateHeaderData = useCallback((title, subLinks) => {
    setTitle(title)
    setSubLinks(subLinks)
  }, [])

  return (
    <HeaderContext.Provider value={{ title, subLinks, updateHeaderData }}>
      {props.children}
    </HeaderContext.Provider>
  )
})

const useHeaderContext = () => useContext(HeaderContext)

const useUpdateHeaderData = () => useHeaderContext().updateHeaderData

export default Header

export { HeaderContextProvider, useHeaderContext, useUpdateHeaderData }
