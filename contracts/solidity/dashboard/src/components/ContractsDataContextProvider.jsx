import React from 'react'
import WithWeb3Context from './WithWeb3Context'
import { displayAmount } from '../utils'

export const ContractsDataContext = React.createContext({})

class ContractsDataContextProvider extends React.Component {
    
    constructor(props) {
        super(props);
        this.state = {
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
        this.getContractsInfo();
    }

    componentDidUpdate(prevProps) {
        if(prevProps.web3.yourAddress !== this.props.web3.yourAddress)
            this.getContractsInfo();
    }

    getContractsInfo = async () => {
        const { web3: { token, stakingContract, grantContract, yourAddress, changeDefaultContract, utils } } = this.props;
        if(!token.options.address || !stakingContract.options.address || !grantContract.options.address)
            return;
        try {
            this.setState({ contractsDataIsFetching: true })
            const tokenBalance = new utils.BN(await token.methods.balanceOf(yourAddress).call());
            const stakeOwner = await stakingContract.methods.ownerOf(yourAddress).call();
            const grantBalance = await grantContract.methods.balanceOf(yourAddress).call()
            const grantStakeBalance = displayAmount(await grantContract.methods.stakeBalanceOf(yourAddress).call(), 18, 3)
            
            let isTokenHolder = tokenBalance.gt(new utils.BN(0));
            let isOperator = stakeOwner !== "0x0000000000000000000000000000000000000000" && utils.toChecksumAddress(yourAddress) !== utils.toChecksumAddress(stakeOwner)
        
            // Check if your account is an operator for a staked Token Grant.
            let stakedGrant
            let isOperatorOfStakedTokenGrant
            let stakedGrantByOperator = await grantContract.methods.grantStakes(yourAddress).call()
        
            if (stakedGrantByOperator.stakingContract === stakingContract.address) {
                isOperatorOfStakedTokenGrant = true
                stakedGrant = await grantContract.methods.grants(stakedGrantByOperator.grantId.toString()).call()
                changeDefaultContract(grantContract);
            }
            
            this.setState({
                isOperator,
                isTokenHolder,
                isOperatorOfStakedTokenGrant,
                tokenBalance,
                grantBalance,
                grantStakeBalance,
                stakedGrant,
                stakeOwner,
                contractsDataIsFetching: false
            })
        } catch(error) {
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

export default WithWeb3Context(ContractsDataContextProvider);

export const withContractsDataContext = (Component) => {
    return (props) => (
      <ContractsDataContext.Consumer>
        {data =>  <Component {...props} {...data} />}
      </ContractsDataContext.Consumer>
    )
}