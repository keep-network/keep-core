import { contractService } from './contracts.service'
import web3Utils from 'web3-utils'
import { TOKEN_STAKING_CONTRACT_NAME, TOKEN_GRANT_CONTRACT_NAME, OPERATOR_CONTRACT_NAME, KEEP_TOKEN_CONTRACT_NAME } from '../constants/constants'

export const fetchTokensPageData = async (web3Context) => {
  const { yourAddress, eth } = web3Context
  let tokenStakingBalance = web3Utils.toBN(0)
  let pendingUndelegationBalance = web3Utils.toBN(0)
  const delegations = []
  const undelegations = []

  const [
    ownedKeepBalance,
    tokenGrantsBalance,
    tokenGrantsStakeBalance,
    minimumStake,
    operatorsAddresses,
    undelegationPeriod,
    initializationPeriod,
  ] = await Promise.all([
    contractService.makeCall(web3Context, KEEP_TOKEN_CONTRACT_NAME, 'balanceOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'balanceOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'stakeBalanceOf', yourAddress),
    contractService.makeCall(web3Context, OPERATOR_CONTRACT_NAME, 'minimumStake'),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'operatorsOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'undelegationPeriod'),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'initializationPeriod'),
  ])

  const operatorsAddressesSet = new Set(operatorsAddresses)
  for (const operatorAddress of operatorsAddressesSet) {
    const {
      createdAt,
      undelegatedAt,
      amount,
    } = await contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'getDelegationInfo', operatorAddress)
    const beneficiary = await contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'magpieOf', operatorAddress)
    const operatorData = { undelegatedAt, amount, beneficiary, operatorAddress, createdAt }
    const balance = web3Utils.toBN(amount)

    if (!balance.isZero() && operatorData.undelegatedAt === '0') {
      const initializationOverAt = web3Utils.toBN(createdAt || 0).add(web3Utils.toBN(initializationPeriod))
      operatorData.isInInitializationPeriod = initializationOverAt.gte(web3Utils.toBN(await eth.getBlockNumber()))
      delegations.push(operatorData)
      tokenStakingBalance = tokenStakingBalance.add(balance)
    }
    if (operatorData.undelegatedAt !== '0') {
      operatorData.undelegationCompleteAt = web3Utils.toBN(undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
      operatorData.canRecoverStake = web3Utils.toBN(await eth.getBlockNumber()).gte(operatorData.undelegationCompleteAt)
      pendingUndelegationBalance = pendingUndelegationBalance.add(balance)
      undelegations.push(operatorData)
    }
  }

  return {
    ownedKeepBalance,
    undelegationPeriod,
    tokenStakingBalance: tokenStakingBalance.toString(),
    pendingUndelegationBalance: pendingUndelegationBalance.toString(),
    tokenGrantsBalance,
    tokenGrantsStakeBalance,
    minimumStake,
    delegations,
    undelegations,
  }
}

const delegateStake = async (web3Context, data, onTransactionHashCallback) => {
  const {
    authorizerAddress,
    beneficiaryAddress,
    operatorAddress,
    stakeTokens,
    context,
    selectedGrant,
  } = data
  const amount = web3Utils.toBN(stakeTokens).mul(web3Utils.toBN(10).pow(web3Utils.toBN(18))).toString()
  const delegation = '0x' + Buffer.concat([
    Buffer.from(beneficiaryAddress.substr(2), 'hex'),
    Buffer.from(operatorAddress.substr(2), 'hex'),
    Buffer.from(authorizerAddress.substr(2), 'hex'),
  ]).toString('hex')

  const { token, stakingContract, grantContract, yourAddress } = web3Context

  if (context === 'owned') {
    await token.methods
      .approveAndCall(stakingContract.options.address, amount, delegation)
      .send({ from: yourAddress })
      .on('transactionHash', onTransactionHashCallback)
  } else if (context === 'granted') {
    await grantContract.methods.stake(selectedGrant.id, stakingContract.options.address, amount, delegation)
      .send({ from: yourAddress })
      .on('transactionHash', onTransactionHashCallback)
  }
}

export const tokensPageService = {
  fetchTokensPageData,
  delegateStake,
}
