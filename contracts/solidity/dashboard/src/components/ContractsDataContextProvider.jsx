import React from 'react'
import withWeb3Context from './WithWeb3Context'
import { displayAmount } from '../utils'

export const ContractsDataContext = React.createContext({})

class ContractsDataContextProvider extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      isBeneficiary: false,
      isTokenHolder: false,
      isOperator: false,
      isOperatorOfStakedTokenGrant: false,
      tokenBalance: '',
      stakedGrant: '',
      stakeOwner: '',
      grantBalance: '',
      grantStakeBalance: '',
      contractsDataIsFetching: true,
    }
  }

  componentDidMount() {
    this.getContractsInfo()
  }

  componentDidUpdate(prevProps) {
    if (prevProps.web3.yourAddress !== this.props.web3.yourAddress || this.areContractsChanged(prevProps)) {
      this.getContractsInfo()
    }
  }

    areContractsChanged = (prevProps) => {
      const { web3 } = this.props
      const tokenAddress = web3.token.options.address
      const stakingAddress = web3.stakingContract.options.address
      const grantAddress = web3.grantContract.options.address

      return tokenAddress !== prevProps.web3.token.options.address ||
            stakingAddress !== prevProps.web3.stakingContract.options.address ||
            grantAddress !== prevProps.web3.grantContract.options.address
    }

    getContractsInfo = async () => {
      const { web3: { token, stakingContract, grantContract, yourAddress, changeDefaultContract, utils } } = this.props
      if (!token.options.address || !stakingContract.options.address || !grantContract.options.address || !yourAddress) {
        return
      }
      try {
        this.setState({ contractsDataIsFetching: true })
        const tokenBalance = new utils.BN(await token.methods.balanceOf(yourAddress).call())
        const stakeOwner = await stakingContract.methods.ownerOf(yourAddress).call()
        const grantBalance = await grantContract.methods.balanceOf(yourAddress).call()
        const grantStakeBalance = displayAmount(await grantContract.methods.stakeBalanceOf(yourAddress).call(), 18, 3)

        const isTokenHolder = tokenBalance.gt(new utils.BN(0))
        const isOperator = stakeOwner !== '0x0000000000000000000000000000000000000000' && utils.toChecksumAddress(yourAddress) !== utils.toChecksumAddress(stakeOwner)

        // Check if your account is an operator for a staked Token Grant.
        let stakedGrant
        let isOperatorOfStakedTokenGrant
        const stakedGrantByOperator = await grantContract.methods.grantStakes(yourAddress).call()

        if (stakedGrantByOperator.stakingContract === stakingContract.address) {
          isOperatorOfStakedTokenGrant = true
          stakedGrant = await grantContract.methods.grants(stakedGrantByOperator.grantId.toString()).call()
          changeDefaultContract(grantContract)
        }

        const operatorsOfMagpie = await stakingContract.methods.operatorsOfMagpie(yourAddress).call()
        const isBeneficiary = operatorsOfMagpie && operatorsOfMagpie.length > 0

        this.setState({
          isOperator,
          isTokenHolder,
          isOperatorOfStakedTokenGrant,
          tokenBalance,
          grantBalance,
          grantStakeBalance,
          stakedGrant,
          stakeOwner,
          contractsDataIsFetching: false,
          isBeneficiary,
        })
      } catch (error) {
        this.setState({ contractsDataIsFetching: false })
      }
    }

    render() {
      return (
        <ContractsDataContext.Provider value={{ ...this.state }}>
          {this.props.children}
        </ContractsDataContext.Provider>
      )
    }
}

export default withWeb3Context(ContractsDataContextProvider)

export const withContractsDataContext = (Component) => {
  return (props) => (
    <ContractsDataContext.Consumer>
      {(data) => <Component {...props} {...data} />}
    </ContractsDataContext.Consumer>
  )
}
