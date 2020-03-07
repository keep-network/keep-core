import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import { registryService } from './registry.service'
import { isSameEthAddress } from '../utils'

const fetchAuthorizationPageData = async (web3Context) => {
  const { yourAddress } = web3Context
  const approvedContractsInRegistry = await registryService.fetchAuthorizedOperatorContracts(web3Context)
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
    if (isSameEthAddress(authorizerOfOperator, yourAddress)) {
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

const fetchOperatorAuthorizedContracts = async (web3Context) => {
  const authorizedContracts = []
  const { yourAddress } = web3Context
  const approvedContractsInRegistry = await registryService.fetchAuthorizedOperatorContracts(web3Context)
  const authorizer = await await contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'authorizerOf', yourAddress)
  if (authorizer === '0x0000000000000000000000000000000000000000') {
    return { isOperator: false, contracts: authorizedContracts }
  }

  for (let i = 0; i < approvedContractsInRegistry.length; i++) {
    const contractAddress = approvedContractsInRegistry[i].contractAddress
    const isAuthorized = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      'isAuthorizedForOperator',
      yourAddress,
      contractAddress
    )
    if (isAuthorized) {
      authorizedContracts.push({ contractAddress, authorizer })
    }
  }

  return { isOperator: true, contracts: authorizedContracts }
}


export const authorizationService = {
  fetchAuthorizationPageData,
  fetchOperatorAuthorizedContracts,
}
