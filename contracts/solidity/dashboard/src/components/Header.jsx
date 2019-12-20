import React from 'react'
import { Navbar, NavbarBrand } from 'react-bootstrap'
import * as Icons from './Icons'
import WithWeb3Context from './WithWeb3Context'
import { Web3Status } from './Web3Status'

const Header = ({ web3: { networkType, token } }) => {
  return (
    <>
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
            <span className="txt-primary">{ token ? token.options.address : '' }</span>
          </div>
          <div>
            <strong>Network: </strong>
            <span className="txt-primary">{ networkType }</span>
          </div>
        </div>
      </Navbar>
      <Web3Status/>
    </>
  )
}

export default WithWeb3Context(Header);
