import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import { registryService } from './registry.service'
import web3Utils from 'web3-utils'

const fetchAuthorizationPageData = async (web3Context) => {
  const { yourAddress } = web3Context
  const approvedContractsInRegistry = await registryService.fetchAuthorizedOperatorContracts(web3Context)
  console.log('aappp', approvedContractsInRegistry)
  const stakedEvents = await contractService.getPastEvents(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'Staked', { fromBlock: '0' })
  const visitedOperators = {}
  const authorizerOperators = []
  const data = {}

  // Fetch all authorizer operators
  for (let i = 0; i < stakedEvents.length; i++) {
    const { returnValues: { from: operatorAddress } } = stakedEvents[i]

    if (visitedOperators.hasOwnProperty(operatorAddress)) {
      continue
    }
    visitedOperators[operatorAddress] = operatorAddress
    const authorizerOfOperator = await contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'authorizerOf', operatorAddress)
    if (web3Utils.toChecksumAddress(yourAddress) === web3Utils.toChecksumAddress(authorizerOfOperator)) {
      authorizerOperators.push(operatorAddress)
    }
  }

  for (let i = 0; i < authorizerOperators.length; i++) {
    const operator = authorizerOperators[i]
    data[operator] = {
      contractsToAuthorize: [],
      authorizedContracts: [],
    }

    for (let j = 0; j < approvedContractsInRegistry.length; j++) {
      const contractAddress = approvedContractsInRegistry[j].contractAddress
      const isAuthorized = await contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'isAuthorizedForOperator', operator, contractAddress)
      if (isAuthorized) {
        data[operator].authorizedContracts.push(approvedContractsInRegistry[j])
      } else {
        data[operator].contractsToAuthorize.push(approvedContractsInRegistry[j])
      }
    }
  }

  return data
}

export const authorizationPageService = {
  fetchAuthorizationPageData,
}
