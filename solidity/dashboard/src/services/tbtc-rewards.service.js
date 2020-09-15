import { contractService } from "./contracts.service"
import {
  TBTC_TOKEN_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TBTC_SYSTEM_CONTRACT_NAME,
} from "../constants/constants"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  createDepositContractInstance,
  createBondedECDSAKeepContractInstance,
} from "../contracts"
import web3Utils from "web3-utils"
import { isSameEthAddress } from "../utils/general.utils"
import { isEmptyArray } from "../utils/array.utils"

const fetchTBTCRewards = async (web3Context, beneficiaryAddress) => {
  const transferEventSearchFilter = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_TOKEN_CONTRACT_NAME],
    filter: { to: web3Utils.toChecksumAddress(beneficiaryAddress) },
  }

  const transferEventToBeneficiary = await contractService.getPastEvents(
    web3Context,
    TBTC_TOKEN_CONTRACT_NAME,
    "Transfer",
    transferEventSearchFilter
  )

  const depositCreatedFilterParam = isEmptyArray(transferEventToBeneficiary)
    ? {}
    : {
        _depositContractAddress: transferEventToBeneficiary.map(
          (_) => _.returnValues.from
        ),
      }
  const depositCreatedSearchFilter = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_SYSTEM_CONTRACT_NAME],
    filter: depositCreatedFilterParam,
  }

  const depositCreatedEvents = await contractService.getPastEvents(
    web3Context,
    TBTC_SYSTEM_CONTRACT_NAME,
    "Created",
    depositCreatedSearchFilter
  )

  const data = transferEventToBeneficiary
    .filter(({ returnValues: { from } }) =>
      depositCreatedEvents.some(
        ({ returnValues: { _depositContractAddress } }) =>
          isSameEthAddress(_depositContractAddress, from)
      )
    )
    .map(({ transactionHash, returnValues: { from, value } }) => ({
      depositTokenId: from,
      amount: value,
      transactionHash,
    }))

  return data
}

const fetchBeneficiaryOperatorsFromDeposit = async (
  web3Context,
  beneficairyAddress,
  depositId
) => {
  const { web3 } = web3Context
  const depositConract = createDepositContractInstance(web3, depositId)

  const keepAddress = await depositConract.methods.getKeepAddress().call()
  const bondedECDSAKeepContract = createBondedECDSAKeepContractInstance(
    web3,
    keepAddress
  )

  const bondedMembers = new Set(
    await bondedECDSAKeepContract.methods.getMembers().call()
  )

  const beneficiaryOperators = []
  for (const operator of bondedMembers) {
    const beneficiaryOfOperator = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "beneficiaryOf",
      operator
    )
    if (isSameEthAddress(beneficiaryOfOperator, beneficairyAddress))
      beneficiaryOperators.push(operator)
  }

  return beneficiaryOperators
}

export const tbtcRewardsService = {
  fetchTBTCRewards,
  fetchBeneficiaryOperatorsFromDeposit,
}
