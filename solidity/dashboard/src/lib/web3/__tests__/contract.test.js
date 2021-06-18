import {
  BaseContract,
  Web3LibWrapper,
  Web3jsWrapper,
  ContractFactory,
} from "../../web3"

describe("Test `BaseContract` wrapper", () => {
  let contract
  const mockedContractInstance = {}
  const deployedAtBlock = "100"
  const deploymentTxnHash = "0x0"
  const web3Wrapper = new Web3LibWrapper({})

  beforeEach(() => {
    contract = new BaseContract(
      mockedContractInstance,
      deploymentTxnHash,
      web3Wrapper
    )

    contract._makeCall = jest.fn()
    contract._sendTransaction = jest.fn()
    contract._getPastEvents = jest.fn()
  })

  it("should call the `makeCall` correctly", async () => {
    const mockedMethodName = "method"
    const mockedArg1 = "arg1"
    const mockedArg2 = 20
    const mockedResult = 1
    contract._makeCall.mockResolvedValue(mockedResult)

    const result = await contract.makeCall(
      mockedMethodName,
      mockedArg1,
      mockedArg2
    )

    expect(contract._makeCall).toHaveBeenCalledWith(
      mockedMethodName,
      mockedArg1,
      mockedArg2
    )
    expect(result).toEqual(mockedResult)
  })

  it("should call the `sendTransaction` correctly", async () => {
    const mockedMethodName = "method"
    const mockedPromiEvent = { on: jest.fn() }
    contract._sendTransaction.mockResolvedValue(mockedPromiEvent)

    const result = await contract.sendTransaction(mockedMethodName)

    expect(contract._sendTransaction).toHaveBeenCalledWith(mockedMethodName)
    expect(result).toEqual(mockedPromiEvent)
  })

  describe("`getPastEvents` test", () => {
    let mockedEventName
    let mockedResult
    let mockedFilter
    let spyOnGetDeploymentBlock

    beforeEach(() => {
      mockedEventName = "event"
      mockedResult = [{ transactionHash: "0x1" }, { transactionHash: "0x2" }]
      mockedFilter = { indexedParam1: "0x0", indexedParam2: 3 }
      contract._getPastEvents.mockResolvedValue(mockedResult)

      spyOnGetDeploymentBlock = jest
        .spyOn(contract, "_getDeploymentBlock")
        .mockResolvedValue(deployedAtBlock)
    })

    it("should call the `getPastEvents` correctly w/o `fromBlock` arg", async () => {
      const result = await contract.getPastEvents(mockedEventName, mockedFilter)
      expect(contract._getPastEvents).toHaveBeenCalledWith(
        mockedEventName,
        mockedFilter,
        deployedAtBlock
      )
      expect(spyOnGetDeploymentBlock).toHaveBeenCalled()
      expect(result).toEqual(mockedResult)
    })

    it("should call the `getPastEvents` correctly w/ `fromBlock` arg", async () => {
      const result = await contract.getPastEvents(
        mockedEventName,
        mockedFilter,
        deployedAtBlock
      )

      expect(contract._getPastEvents).toHaveBeenCalledWith(
        mockedEventName,
        mockedFilter,
        deployedAtBlock
      )
      expect(spyOnGetDeploymentBlock).not.toHaveBeenCalled()
      expect(result).toEqual(mockedResult)
    })
  })

  it("should fetch the deployment block via the web3 wrapper and save in `deployedAtBlock` field", async () => {
    const mockedTxObj = { blockNumber: 100 }
    const web3Wrapper = contract.web3
    const spyOnGetTx = jest
      .spyOn(web3Wrapper, "getTransaction")
      .mockResolvedValue(mockedTxObj)

    const result = await contract._getDeploymentBlock()
    const result2 = await contract._getDeploymentBlock()

    expect(spyOnGetTx).toHaveBeenCalledTimes(1)
    expect(spyOnGetTx).toHaveBeenCalledWith(contract.deploymentTxnHash)
    expect(result).toEqual(mockedTxObj.blockNumber.toString())
    expect(contract.deployedAtBlock).toEqual(result)
    expect(contract.deployedAtBlock).toEqual(result2)
  })
})

describe("Test Web3jsContract wrapper", () => {
  let contract
  const mockedMethodName = "mockedMethod"
  const mockedMethodCall = jest.fn()
  const mockedMethodSend = jest.fn()

  const mockedWebjsContractInstance = {
    methods: {
      [mockedMethodName]: jest.fn(() => ({
        call: mockedMethodCall,
        send: mockedMethodSend,
      })),
    },
    getPastEvents: jest.fn(),
    options: {
      defaultAccount: null,
      address: "0x0",
    },
  }
  const deploymentTxnHash = "0x0"
  const web3Wrapper = new Web3jsWrapper({})

  beforeEach(async () => {
    contract = await ContractFactory.createWeb3jsContract(
      mockedWebjsContractInstance,
      deploymentTxnHash,
      web3Wrapper
    )
  })

  it("should call a provided method", async () => {
    const mockedArg1 = "0x123"
    await contract.makeCall(mockedMethodName, mockedArg1)

    expect(
      mockedWebjsContractInstance.methods[mockedMethodName]
    ).toHaveBeenCalledWith(mockedArg1)
    expect(mockedMethodCall).toHaveBeenCalled()
  })

  it("should send transaction", async () => {
    const mockedArg1 = "0x123"
    const mockedArg2 = 1

    await contract.sendTransaction(mockedMethodName, mockedArg1, mockedArg2)

    expect(
      mockedWebjsContractInstance.methods[mockedMethodName]
    ).toHaveBeenCalledWith(mockedArg1, mockedArg2)
    expect(mockedMethodSend).toHaveBeenCalled()
  })

  it("should get past events", async () => {
    const eventName = "event"
    const filter = { indexedParam1: 1, indexedParam2: 2 }
    const fromBlock = 1

    await contract.getPastEvents(eventName, filter, fromBlock)

    expect(
      mockedWebjsContractInstance.getPastEvents
    ).toHaveBeenCalledWith(eventName, { fromBlock, filter })
  })

  it("should return the contract instance address", () => {
    const result = contract.address

    expect(result).toEqual(mockedWebjsContractInstance.options.address)
  })

  it("should set the default account for contract", () => {
    const mockedDefaultAccount = "0x123456789"

    contract.defaultAccount = mockedDefaultAccount

    expect(contract.defaultAccount).toEqual(mockedDefaultAccount)
  })
})
