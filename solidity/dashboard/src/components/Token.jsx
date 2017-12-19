import React, { Component } from 'react';
import Header from './Header';
import { Col, Grid, Row } from 'react-bootstrap'
import Network from '../network'
import { getTokenVesting } from '../contracts'

class Token extends Component {

  constructor() {
    super()
    this.state = { canRevoke: false }
  }

  async componentWillReceiveProps(nextProps) {
    const { owner, revoked } = nextProps.details
    const accounts = await Network.getAccounts()

    const isOwner = accounts[0]
      ? owner === accounts[0].toLowerCase()
      : undefined

    this.setState({ accounts, canRevoke: isOwner && ! revoked })
  }

  render() {
    const { address } = this.props.details;
    return
      <div className="token">
        <Header />
        <p>Keep token at { address }</p>
        <Grid>
          <Row>
            <Col xs={12} md={6}>
            </Col>  
            <Col xs={12} md={6}>
            </Col>
          </Row>
        </Grid>
      </div>
  }
}

export default Token
