import React from 'react'
import { Col, Row } from 'react-bootstrap'
import VestingChart from './VestingChart'
import VestingDetails from './VestingDetails'
import TokenGrants from './TokenGrants'
import StakingDelegateTokenGrantForm from './StakingDelegateTokenGrantForm'
import FetchTokenGrantsCore from './FetchTokenGrantsCore'
import { withContractsDataContext } from './ContractsDataContextProvider'

class TokenGrantsTab extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      selectedGrantIndex: 0,
    }
  }

    selectTokenGrant = (selectedGrantIndex) => this.setState({ selectedGrantIndex })

    selectedGrant = () => this.props.data ? this.props.data[this.state.selectedGrantIndex] : {}

    render() {
      const { data, grantBalance } = this.props
      return (
        <>
          <h3>Tokens granted to you</h3>
          <Row>
            <Col xs={12} md={6}>
              <VestingDetails
                details={this.selectedGrant()}
              />
            </Col>
            <Col xs={12} md={6}>
              <VestingChart
                details={this.selectedGrant()}
              />
            </Col>
          </Row>
          <Row>
            <Col xs={12} md={12}>
              <TokenGrants
                data={data}
                selectTokenGrant={this.selectTokenGrant}
              />
            </Col>
          </Row>
          <Row>
            <Col xs={12} md={12}>
              <h3>Stake Delegation of Token Grants</h3>
              <p>
                            Keep network does not require token owners to perform the day-to-day operations of staking
                            with the private keys holding the tokens. This is achieved by stake delegation, where different
                            addresses hold different responsibilities and cold storage is supported to the highest extent practicable.
              </p>

              <StakingDelegateTokenGrantForm
                tokenBalance={grantBalance}
              />
            </Col>
          </Row>
        </>
      )
    }
}
const TokenGrantsTabWithContractDataContext = withContractsDataContext(TokenGrantsTab)

export default (props) => (
  <FetchTokenGrantsCore userGrants >
    <TokenGrantsTabWithContractDataContext />
  </FetchTokenGrantsCore>
)
