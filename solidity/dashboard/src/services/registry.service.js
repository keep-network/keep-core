import { contractService } from "./contracts.service"
import { REGISTRY_CONTRACT_NAME } from "../constants/constants"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"

const fetchAuthorizedOperatorContracts = async (web3Context) => {
  const options = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[REGISTRY_CONTRACT_NAME],
  }
  const events = await contractService.getPastEvents(
    web3Context,
    REGISTRY_CONTRACT_NAME,
    "OperatorContractApproved",
    options
  )

  const contracts = {}
  for (let i = 0; i < events.length; i++) {
    const {
      blockNumber,
      returnValues: { operatorContract },
    } = events[i]
    if (contracts.hasOwnProperty(operatorContract)) {
      continue
    }
    const isAuthorized = await contractService.makeCall(
      web3Context,
      REGISTRY_CONTRACT_NAME,
      "isApprovedOperatorContract",
      operatorContract
    )
    if (isAuthorized) {
      contracts[operatorContract] = {
        contractAddress: operatorContract,
        blockNumber,
      }
    }
  }
  return Object.keys(contracts).map((contract) => contracts[contract])
}

export const registryService = {
  fetchAuthorizedOperatorContracts,
}
