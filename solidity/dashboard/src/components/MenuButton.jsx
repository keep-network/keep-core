import React, { useContext } from 'react'
import { SideMenuContext } from './SideMenu'

export const MenuButton = (proops) => {
  const { isOpen, toggle } = useContext(SideMenuContext)

  return (
    <div className={`${isOpen ? 'active ' : ''}menu-button`} onClick={toggle} />
  )
}
