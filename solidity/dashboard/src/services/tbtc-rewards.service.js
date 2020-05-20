import { contractService } from "./contracts.service"
import { TBTC_TOKEN_CONTRACT_NAME } from "../constants/constants"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import web3Utils from "web3-utils"

const fetchTBTCReawrds = async (web3Context, beneficiaryAddress) => {
  const searchFilter = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_TOKEN_CONTRACT_NAME],
    filter: { to: web3Utils.toChecksumAddress(beneficiaryAddress) },
  }

  const transferEventToBeneficiary = (
    await contractService.getPastEvents(
      web3Context,
      TBTC_TOKEN_CONTRACT_NAME,
      "Transfer",
      searchFilter
    )
  ).map(({ fromBlock, returnValues: { from, amount } }) => ({
    depositTokenId: from,
    amount,
    date: fromBlock,
  }))

  return transferEventToBeneficiary
}

export const tbtcRewardsService = { fetchTBTCReawrds }
