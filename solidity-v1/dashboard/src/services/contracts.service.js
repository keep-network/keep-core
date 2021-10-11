const makeCall = async (
  web3Context,
  contractName,
  contractMethodName,
  ...args
) => {
  return await web3Context[contractName].methods[contractMethodName](
    ...args
  ).call()
}

const sendTransaction = async (
  web3Context,
  contractName,
  contractMethodName,
  sendArgs,
  ...args
) => {
  return await web3Context[contractName].methods[contractMethodName](
    ...args
  ).send(sendArgs)
}

const getPastEvents = async (web3Context, contractName, eventName, ...args) => {
  return await web3Context[contractName].getPastEvents(eventName, ...args)
}
export const contractService = {
  makeCall,
  sendTransaction,
  getPastEvents,
}
