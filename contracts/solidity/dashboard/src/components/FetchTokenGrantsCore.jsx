import React from 'react'
import moment from 'moment'
import withWeb3Context from './WithWeb3Context'
import Loadable from './Loadable'
import { displayAmount } from '../utils'

class FetchTokenGrantsCore extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      grantedToYou: [],
      isFetching: true,
    }
  }

  componentDidUpdate(prevProps) {
    if (prevProps.web3.yourAddress !== this.props.web3.yourAddress) {
      this.getData()
    }
  }

  async componentDidMount() {
    await this.getData()
  }

    getData = async () => {
      try {
        const { web3: { grantContract, yourAddress } } = this.props
        if (!grantContract.options.address) {
          return
        }
        const grantIndexes = await grantContract.methods.getGrants(yourAddress).call()
        const grantedToYou = (await this.getGrants(grantIndexes))
          .filter(this.checkIsGrantedToYou)
          .map(this.mapToChartData)

        this.setState({
          grantedToYou: await Promise.all(grantedToYou),
          isFetching: false,
        })
      } catch (error) {
        this.setState({ isFetching: false, error: 'Cannot load data' })
      }
    }

    getGrants = (grantIndexes) => {
      return Promise.all(grantIndexes.map(this.mapToGrant))
    }

    mapToGrant = async (grantIndex) => {
      const { web3: { grantContract } } = this.props
      return { grantIndex, grant: await grantContract.methods.grants(grantIndex).call() }
    }

    checkIsGrantedToYou = ({ grant }) => {
      const { web3: { utils, yourAddress }, userGrants } = this.props
      return utils.toChecksumAddress(grant[userGrants ? 'grantee' : 'grantManager']) === utils.toChecksumAddress(yourAddress)
    }

    mapToChartData = async ({ grantIndex, grant }) => {
      const { web3: { utils, grantContract } } = this.props
      const grantedAmount = await grantContract.methods.grantedAmount(grantIndex).call()
      return {
        id: grantIndex,
        grantManager: utils.toChecksumAddress(grant.grantManager),
        grantee: utils.toChecksumAddress(grant.grantee),
        revoked: grant.revoked,
        revocable: grant.revocable,
        amount: grant.amount,
        grantedAmount: grantedAmount,
        end: new utils.BN(grant.start).add(new utils.BN(grant.duration)),
        start: grant.start,
        cliff: grant.cliff,
        withdrawn: grant.withdrawn,
        staked: grant.staked,
        decimals: 18,
        symbol: 'KEEP',
        formatted: {
          amount: displayAmount(grant.amount, 18, 3),
          end: moment(new utils.BN(grant.start).add(new utils.BN(grant.duration)).toNumber() * 1000).format('MMMM Do YYYY, h:mm:ss a'),
          start: moment(grant.start * 1000).format('MMMM Do YYYY, h:mm:ss a'),
          cliff: moment(grant.cliff * 1000).format('MMMM Do YYYY, h:mm:ss a'),
          withdrawn: grant.withdrawn,
        },
      }
    }

    render() {
      const { isFetching } = this.state

      return isFetching ? <Loadable /> : React.cloneElement(this.props.children, { data: this.state.grantedToYou })
    }
}

export default withWeb3Context(FetchTokenGrantsCore)
