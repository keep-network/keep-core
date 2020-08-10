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
