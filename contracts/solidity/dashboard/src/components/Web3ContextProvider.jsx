import React from 'react'
import { getWeb3, getWeb3SocketProvider } from '../utils'
import { Web3Context } from './WithWeb3Context'
import { getKeepToken, getTokenStaking, getTokenGrant } from '../contracts'
import { MessagesContext, messageType } from './Message'

export default class Web3ContextProvider extends React.Component {
    static contextType = MessagesContext

    constructor(props) {
        super(props)
        this.state = {
            web3: null,
            isFetching: false,
            yourAddress: '',
            networkType: '',
            token: { options: { address: '' } },
            stakingContract: { options: { address: '' } },
            grantContract: { options: { address: '' } },
            utils: {},
            eth: {},
            error: '',
        }
    }

    componentDidMount() {
       this.initialize()
    }

    initialize = async () => {
        const web3 = getWeb3()
        if(!web3) {
            return
        }
        this.setState({ web3 }, this.setData)
    }

    setData = async () => {
        this.connectAppWithAccount()
        this.initializeContracts()
        this.state.web3.eth.currentProvider.on('accountsChanged', this.accountHasBeenChanged)
    }

    connectAppWithAccount = async () => {
        const { web3 } = this.state
        this.setState({ isFetching: true })
        this.context.showMessage({ type: messageType.INFO, title: 'Please check web3 provider' })
        
        try{
            const [account] = await web3.currentProvider.enable()
            this.setState({
                yourAddress: account,
                networkType: await web3.eth.net.getNetworkType(),
                isFetching: false
            })
        } catch(error) {
            this.context.showMessage({ type: 'error', title: error.message })
            this.setState({ isFetching: false })
        }      
    }

    initializeContracts = async () => {
        const { web3 } = this.state
        try {
            const web3EventProvider = getWeb3SocketProvider()
            const [token, grantContract, stakingContract] = await this.getContracts(web3)
            const [eventToken, eventGrantContract, eventStakingContract] = await this.getContracts(web3EventProvider)
            this.setState({
                token,
                grantContract,
                stakingContract,
                defaultContract: stakingContract,
                utils: web3.utils,
                eth: web3.eth,
                eventToken,
                eventGrantContract,
                eventStakingContract
            })
        } catch(error) {
            this.setState({
                error: "Failed to load contracts. Please check if Metamask is enabled and connected to the correct network.",
            })
        }
    }

    getContracts = async (web3) => await Promise.all([
        getKeepToken(web3),
        getTokenGrant(web3),
        getTokenStaking(web3)
    ])

    accountHasBeenChanged = ([yourAddress]) => {
        if(!yourAddress) {
            this.setState({
                isFetching: false,
                yourAddress: '',
                token: { options: { address: '' } },
                stakingContract: { options: { address: '' } },
                grantContract: { options: { address: '' } },
            })
            return
        }
        this.setState({ yourAddress })
    }

    changeDefaultContract = (defaultContract) => {
        this.setState({ defaultContract })
    }

    render() {
        return (
            <Web3Context.Provider value={{ ...this.state, changeDefaultContract: this.changeDefaultContract, connectAppWithAccount: this.connectAppWithAccount }}>
                {this.props.children}
            </Web3Context.Provider>    
        )
    }
}