import React from "react"
import withWeb3Context from "./WithWeb3Context"
import { isSameEthAddress } from "../utils/general.utils"
import { getKeepTokenContractDeployerAddress } from "../contracts"

export const ContractsDataContext = React.createContext({})
class ContractsDataContextProvider extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      tokenBalance: "",
      isKeepTokenContractDeployer: false,
      contractsDataIsFetching: true,
    }
  }

  componentDidMount() {
    this.getContractsInfo()
  }

  componentDidUpdate(prevProps) {
    if (
      prevProps.web3.yourAddress !== this.props.web3.yourAddress ||
      this.areContractsChanged(prevProps)
    ) {
      this.getContractsInfo()
    }
  }

  areContractsChanged = (prevProps) => {
    const { web3 } = this.props
    const tokenAddress = web3.token.options.address
    const stakingAddress = web3.stakingContract.options.address
    const grantAddress = web3.grantContract.options.address

    return (
      tokenAddress !== prevProps.web3.token.options.address ||
      stakingAddress !== prevProps.web3.stakingContract.options.address ||
      grantAddress !== prevProps.web3.grantContract.options.address
    )
  }

  getContractsInfo = async () => {
    const {
      web3: { web3, token, stakingContract, grantContract, yourAddress, utils },
    } = this.props
    if (
      !token.methods ||
      !stakingContract.methods ||
      !grantContract.methods ||
      !yourAddress
    ) {
      return
    }
    try {
      this.setState({ contractsDataIsFetching: true })
      const tokenBalance = new utils.BN(
        await token.methods.balanceOf(yourAddress).call()
      )
      const keepTokenContractDeployerAddress = await getKeepTokenContractDeployerAddress(
        web3
      )
      const isKeepTokenContractDeployer = isSameEthAddress(
        yourAddress,
        keepTokenContractDeployerAddress
      )

      this.setState({
        isKeepTokenContractDeployer,
        tokenBalance,
        contractsDataIsFetching: false,
      })
    } catch (error) {
      this.setState({ contractsDataIsFetching: false })
    }
  }

  refreshKeepTokenBalance = async () => {
    const {
      web3: { token, yourAddress, utils },
    } = this.props

    const tokenBalance = new utils.BN(
      await token.methods.balanceOf(yourAddress).call()
    )
    this.setState({
      tokenBalance,
    })
  }

  render() {
    return (
      <ContractsDataContext.Provider
        value={{
          refreshKeepTokenBalance: this.refreshKeepTokenBalance,
          ...this.state,
        }}
      >
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
