import { contracts } from "../contracts"

export const commitTopUp = async (operator, onTransactionHashCallback) => {
  await contracts.stakingContract.methods
    .commitTopUp(operator)
    .send()
    .on("transactionHash", onTransactionHashCallback)
}
