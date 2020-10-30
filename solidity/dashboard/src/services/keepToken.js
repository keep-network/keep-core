import { ContractsLoaded } from "../contracts"

const balanceOf = async (address) => {
  const { token } = await ContractsLoaded

  return await token.methods.balanceOf(address).call()
}

const keepToken = {
  balanceOf,
}

export default keepToken
