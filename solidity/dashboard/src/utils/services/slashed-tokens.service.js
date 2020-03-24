import { contractService } from './contracts.service'
import { OPERATOR_CONTRACT_NAME } from '../constants/constants'
import moment from 'moment'
import web3Utils from 'web3-utils'

const fetchSlashedTokens = async (web3Context, ...args) => {
  const { yourAddress, eth } = web3Context

  // return [
  //   {
  //     groupPublicKey: '0x2A489EacBf4de172B4018D2b4a405F05C400f530',
  //     date: moment(),
  //     amount: 200000000000000000000,
  //     typeOfPunishment: 1,
  //     id: 1,
  //   },
  //   {
  //     groupPublicKey: '0x2A489EacBf4de172B4018D2b4a405F05C400f530',
  //     date: moment(),
  //     amount: 200000000000000000000,
  //     typeOfPunishment: 0,
  //     id: 2,
  //   },
  // ]

  const acitveGroupsCount = await contractService.makeCall(web3Context, OPERATOR_CONTRACT_NAME, 'numberOfGroups')
  const expiredGroupsCount = await contractService.makeCall(web3Context, OPERATOR_CONTRACT_NAME, 'getFirstActiveGroupIndex')

  const groupIndexToPublicKey = {}

  if (web3Utils.toBN(acitveGroupsCount).add(web3Utils.toBN(expiredGroupsCount)).isZero()) {
    return []
  }

  for (let groupIndex = acitveGroupsCount; groupIndex < acitveGroupsCount + expiredGroupsCount; groupIndex++) {
    const groupPublicKey = await contractService.makeCall(web3Context, OPERATOR_CONTRACT_NAME, 'getGroupPublicKey', groupIndex)
    const memebrIndices = await contractService
      .makeCall(web3Context, OPERATOR_CONTRACT_NAME, 'getGroupMemberIndices', groupPublicKey, yourAddress)

    if (memebrIndices) {
      groupIndexToPublicKey[groupIndex] = groupPublicKey
    }
  }

  if (Object.keys(groupIndexToPublicKey).length === 0) {
    return []
  }

  const events = contractService.getPastEvents(web3Context, OPERATOR_CONTRACT_NAME, 'GroupPunishment',
    { fromBlock: '0', filter: { groupIndex: Object.keys(groupIndexToPublicKey) } })

  return Promise.all(events
    .map(async (event) => {
      const { eventId, blockNumber, returnValues: { groupIndex, seizedTokensPerMember, typeOfPunishment } } = event
      const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp

      return {
        groupPublicKey: groupIndexToPublicKey[groupIndex],
        date: moment.unix(withdrawnAt),
        amount: seizedTokensPerMember,
        typeOfPunishment,
        id: eventId,
      }
    }).sort((eventA, eventB) => eventA.date - eventB.date)
  )
}

export const slashedTokensService = {
  fetchSlashedTokens,
}
