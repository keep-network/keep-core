import React from 'react'
import { Navbar, NavbarBrand } from 'react-bootstrap'
import * as Icons from './Icons';

const Header = () =>  {
  return ( 
    <Navbar>
      <Navbar.Header>
        <NavbarBrand>
          <Icons.Keep height="61px" width="235px"/>
        </NavbarBrand>
        <Navbar.Toggle />
      </Navbar.Header>
      <Navbar.Collapse>
      </Navbar.Collapse>
    </Navbar>
  )
}

export default Header