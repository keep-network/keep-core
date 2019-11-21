import React from 'react';
import { getWeb3 } from '../utils';
import { Web3Context } from './WithWeb3Context';
import { getKeepToken, getTokenStaking, getTokenGrant } from '../contracts';


export default class Web3ContextProvider extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            yourAddress: '',
            networkType: '',
            token: { options: { address: '' } },
            stakingContract: { options: { address: '' } },
            grantContract: { options: { address: '' } },
            utils: {},
            eth: {},
            error: '',
            dataIsReady: false,
        }
    }

    async componentDidMount() {
        try {
            await this.initialize();
        } catch(e) {
            console.log('onit', e)
        }
    }

    initialize = async () => {
        try {
            const web3 = await getWeb3();
            if (!web3) {
                this.setState({
                    error: "No network detected. Do you have MetaMask installed?",
                })
                return
            }

            window.ethereum.on('accountsChanged', this.accountHasBeenChanged)
            const contracts = await this.getContracts(web3);
            if (!contracts) {
                this.setState({
                    error: "Failed to load contracts. Please check if Metamask is enabled and connected to the correct network.",
                })
                return
            }

            this.setState({
                ...contracts,
                yourAddress: (await web3.eth.getAccounts())[0],
                networkType: await web3.eth.net.getNetworkType(),
                defaultContract: contracts.stakingContract,
                utils: web3.utils,
                eth: web3.eth,
                dataIsReady: true,
            })
        } catch(e) {
            console.log('errrrror', e);
        }
    }

    getContracts = async (web3) => {

        try {
          const token = await getKeepToken(web3)
          const stakingContract = await getTokenStaking(web3)
          const grantContract = await getTokenGrant(web3)
          return {
            token: token,
            stakingContract: stakingContract,
            grantContract: grantContract
          }
        } catch (e) {
            console.log('error get contracts', e);
          return null
        }
    }

    accountHasBeenChanged = ([yourAddress]) => {
        console.log('new address', yourAddress);
        this.setState({ yourAddress })
    }

    changeDefaultContract = (defaultContract) => {
        this.setState({ defaultContract })
    };

      render() {
          const { dataIsReady } = this.state;
          console.log('state', this.state);
          return (
              <Web3Context.Provider value={{ ...this.state, changeDefaultContract: this.changeDefaultContract }}>
                  {!dataIsReady ? 'Loading.....': this.props.children}
              </Web3Context.Provider>    
          );
      }
}