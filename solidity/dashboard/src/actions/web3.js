import {
  TOKEN_STAKING_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  KEEP_TOKEN_CONTRACT_NAME,
  OPERATOR_CONTRACT_NAME,
} from "../constants/constants"

const WEB3_SEND_TRANSACTION = "web3/send_transaction"

export const undelegateStake = (operator, meta) => {
  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: TOKEN_STAKING_CONTRACT_NAME,
      methodName: "undelegate",
      args: [operator],
    },
    meta,
  }
}

export const cancelStake = (operator, meta) => {
  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: TOKEN_STAKING_CONTRACT_NAME,
      methodName: "cancel",
      args: [operator],
    },
    meta,
  }
}

export const keepTokenApproveAndCall = (data, meta) => {
  const { tokenAddress, amount, extraData } = data

  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: KEEP_TOKEN_CONTRACT_NAME,
      methodName: "approveAndCall",
      args: [tokenAddress, amount, extraData],
    },
    meta,
  }
}

export const depositEthForOperator = (operatorAddress, amountInWei, meta) => {
  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: KEEP_BONDING_CONTRACT_NAME,
      methodName: "deposit",
      args: [operatorAddress],
      options: { value: amountInWei },
    },
    meta,
  }
}

export const withdrawUnbondedEth = (weiToWithdraw, operatorAddress, meta) => {
  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: KEEP_BONDING_CONTRACT_NAME,
      methodName: "withdraw",
      args: [weiToWithdraw, operatorAddress],
    },
    meta,
  }
}

export const withdrawUnbondedEthAsManagedGrantee = (
  weiToWithdraw,
  operatorAddress,
  managedGrantAddress,
  meta
) => {
  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: KEEP_BONDING_CONTRACT_NAME,
      methodName: "withdrawAsManagedGrantee",
      args: [weiToWithdraw, operatorAddress, managedGrantAddress],
    },
    meta,
  }
}

export const authorizeOperatorContract = (data, meta) => {
  const { operatorAddress, operatorContractAddress } = data

  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: TOKEN_STAKING_CONTRACT_NAME,
      methodName: "authorizeOperatorContract",
      args: [operatorAddress, operatorContractAddress],
    },
    meta,
  }
}

export const authorizeSortitionPoolContract = (data, meta) => {
  const { operatorAddress, sortitionPoolAddress } = data

  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: KEEP_BONDING_CONTRACT_NAME,
      methodName: "authorizeSortitionPoolContract",
      args: [operatorAddress, sortitionPoolAddress],
    },
    meta,
  }
}

export const deauthorizeSortitionPoolContract = (data, meta) => {
  const { operatorAddress, sortitionPoolAddress } = data

  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: KEEP_BONDING_CONTRACT_NAME,
      methodName: "deauthorizeSortitionPoolContract",
      args: [operatorAddress, sortitionPoolAddress],
    },
    meta,
  }
}

export const recoverStake = (operator, meta) => {
  return {
    type: WEB3_SEND_TRANSACTION,
    payload: {
      contractName: TOKEN_STAKING_CONTRACT_NAME,
      methodName: "recoverStake",
      args: [operator],
    },
    meta,
  }
}

export const releaseTokens = (data, meta) => {
  return {
    type: "token-grant/release_tokens",
    payload: data,
    meta,
  }
}
