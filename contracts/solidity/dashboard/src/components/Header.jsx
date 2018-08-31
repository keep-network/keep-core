import React from 'react'
import { Navbar, NavbarBrand } from 'react-bootstrap'
import * as Icons from './Icons'

const Header = ({networkType}) => {
  return (
    <Navbar>
      <Navbar.Header>
        <NavbarBrand>
          <a href="/"><Icons.Keep height="61px" width="235px"/></a>
          <p>Token Dashboard</p>
        </NavbarBrand>
      </Navbar.Header>
      <div className="pull-right">
        <div>
          <strong>KEEP Token: </strong>
          <span className="txt-primary">{ process.env.REACT_APP_TOKEN_ADDRESS }</span>
        </div>
        <div>
          <strong>Network: </strong>
          <span className="txt-primary">{ networkType }</span>
        </div>
      </div>
    </Navbar>
  )
}

export default Header
