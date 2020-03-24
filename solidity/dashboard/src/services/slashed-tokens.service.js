import { contractService } from './contracts.service'
import { OPERATOR_CONTRACT_NAME, TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from '../contracts'

const fetchSlashedTokens = async (web3Context, ...args) => {
  const { yourAddress } = web3Context
  const operatorEventsSearchFilters = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[OPERATOR_CONTRACT_NAME],
  }

  const eventsSearchFilters = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME],
    filter: { operator: yourAddress }
  }
  const data = []

  const slashedTokensEvents = await contractService
    .getPastEvents(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      'TokensSlashed',
      eventsSearchFilters,
    )
  const seizedTokensEvents = await contractService
    .getPastEvents(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      'TokensSeized',
      eventsSearchFilters,
    )
  
  const unauthorizedSigningEvents = contractService
    .getPastEvents(
      web3Context,
      OPERATOR_CONTRACT_NAME,
      'UnauthorizedSigningReported',
      operatorEventsSearchFilters,
    )
  
  const relayEntryTimeoutEvents = contractService
    .getPastEvents(
      web3Context,
      OPERATOR_CONTRACT_NAME,
      'RelayEntryTimeoutReported',
      operatorEventsSearchFilters,
    )

  const punishmentEvents = [...unauthorizedSigningEvents, ...relayEntryTimeoutEvents]
  
  for(i = 0; i < punishmentEvents.length; i++) {
    const { event, transactionHash, blockNumber, returnValues: { groupIndex } } = unauthorizedSigningEvents[i]
    let punishmentData = {}
    if (slashedTokensEvents.hasOwnProperty(transactionHash)) {
      const { amount } = slashedTokensEvents[transactionHash]
      punishmentData = { amount, type: 'SLASHED', event }
    } else if(seizedTokensEvents.hasOwnProperty(transactionHash)) {
      const { amount } = slashedTokensEvents[transactionHash]
      punishmentData = { amount, type: 'SEIZED', event }
    }

    punishmentData.date = await eth.getBlock(blockNumber).timestamp
    punishmentData.groupPublicKey = await contractService.makeCall(web3Context, OPERATOR_CONTRACT_NAME, 'getGroupPublicKey', groupIndex)
    data.push(punishmentData)
  }

  return data
}

export const slashedTokensService = {
  fetchSlashedTokens,
}
