import web3Utils from "web3-utils"
import {
  TBTC_TOKEN_CONTRACT_NAME,
  TBTC_SYSTEM_CONTRACT_NAME,
} from "../constants/constants"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  createDepositContractInstance,
  createBondedECDSAKeepContractInstance,
  ContractsLoaded,
  Web3Loaded,
} from "../contracts"
import { isSameEthAddress } from "../utils/general.utils"
import { isEmptyArray } from "../utils/array.utils"

const fetchTBTCRewards = async (beneficiaryAddress) => {
  const { tbtcTokenContract, tbtcSystemContract } = await ContractsLoaded

  if (!beneficiaryAddress) {
    return []
  }

  const transferEventSearchFilter = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_TOKEN_CONTRACT_NAME],
    filter: { to: web3Utils.toChecksumAddress(beneficiaryAddress) },
  }

  const transferEventToBeneficiary = await tbtcTokenContract.getPastEvents(
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

  const depositCreatedEvents = await tbtcSystemContract.getPastEvents(
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
  beneficairyAddress,
  depositId
) => {
  const web3 = await Web3Loaded
  const { stakingContract } = await ContractsLoaded
  const depositConract = createDepositContractInstance(web3, depositId)

  const keepAddress = await depositConract.methods.keepAddress().call()
  const bondedECDSAKeepContract = createBondedECDSAKeepContractInstance(
    web3,
    keepAddress
  )

  const bondedMembers = new Set(
    await bondedECDSAKeepContract.methods.getMembers().call()
  )

  const beneficiaryOperators = []
  for (const operator of bondedMembers) {
    const beneficiaryOfOperator = await stakingContract.methods
      .beneficiaryOf(operator)
      .call()
    if (isSameEthAddress(beneficiaryOfOperator, beneficairyAddress))
      beneficiaryOperators.push(operator)
  }

  return beneficiaryOperators
}

export const tbtcRewardsService = {
  fetchTBTCRewards,
  fetchBeneficiaryOperatorsFromDeposit,
}
