import { Web3LibWrapper, Web3jsWrapper } from ".."
import { ContractFactory } from "../contract"

describe("Test Web3 lib wrapper", () => {
  const mockedWeb3Lib = {}
  let web3Wrapper

  beforeEach(() => {
    web3Wrapper = new Web3LibWrapper(mockedWeb3Lib)
    web3Wrapper._getTransaction = jest.fn()
    web3Wrapper._createContractInstance = jest.fn()
  })

  it("should create web3 lib wrapper correctly", () => {
    expect(web3Wrapper.lib).toEqual(mockedWeb3Lib)
  })

  it("should call `getTransaction` correctly", async () => {
    const mockedResult = { blockNumber: 1 }
    const hash = "0x0"
    web3Wrapper._getTransaction.mockResolvedValue(mockedResult)

    const result = await web3Wrapper.getTransaction(hash)

    expect(web3Wrapper._getTransaction).toHaveBeenCalledWith(hash)
    expect(result).toEqual(mockedResult)
  })

  it("should call `createContractInstance` correctly", async () => {
    const mockedResult = {
      methods: { mockedMethod: { call: () => {}, send: () => {} } },
    }
    const abi = []
    const address = "0x123"
    const deploymentTxHash = "0x123456789"
    const deployetAtBlock = 100

    web3Wrapper._createContractInstance.mockResolvedValue(mockedResult)

    const result = await web3Wrapper.createContractInstance(
      abi,
      address,
      deploymentTxHash,
      deployetAtBlock
    )

    expect(web3Wrapper._createContractInstance).toHaveBeenCalledWith(
      abi,
      address,
      deploymentTxHash,
      deployetAtBlock
    )
    expect(result).toEqual(mockedResult)
  })
})

describe("Test Web3.js lib wrapper", () => {
  const mockedWeb3Lib = {
    eth: {
      getTransaction: jest.fn(),
      Contract: jest.fn(),
      defaultAccount: null,
    },
    setProvider: jest.fn(),
  }

  let web3Wrapper

  beforeEach(() => {
    web3Wrapper = new Web3jsWrapper(mockedWeb3Lib)
  })

  it("should return transaction data by transaction hash", async () => {
    const hash = "0x0"
    const mockedResult = { blockNumber: 1 }
    mockedWeb3Lib.eth.getTransaction.mockResolvedValue(mockedResult)

    const result = await web3Wrapper.getTransaction(hash)

    expect(mockedWeb3Lib.eth.getTransaction).toHaveBeenCalledWith(hash)
    expect(result).toEqual(mockedResult)
  })

  it("should create web3 contract instance", () => {
    const abi = []
    const address = "0x0"
    const deploymentTxHash = "0x0123"

    const mockedWeb3jsContractInstance = {
      methods: {},
      options: { defaultAccount: null },
    }

    mockedWeb3Lib.eth.Contract.mockImplementation(
      () => mockedWeb3jsContractInstance
    )
    const spyOnContractFactory = jest
      .spyOn(ContractFactory, "createWeb3jsContract")
      .mockReturnValue({})

    web3Wrapper.createContractInstance(abi, address, deploymentTxHash)

    expect(mockedWeb3Lib.eth.Contract).toHaveBeenCalledWith(abi, address)
    expect(spyOnContractFactory).toHaveBeenCalledWith(
      mockedWeb3jsContractInstance,
      deploymentTxHash,
      web3Wrapper,
      null
    )
  })

  it("should set the default account correctly", () => {
    const acc = "0x0123456789"

    web3Wrapper.defaultAccount = acc

    expect(web3Wrapper.defaultAccount).toEqual(acc)
    expect(mockedWeb3Lib.eth.defaultAccount).toEqual(acc)
  })

  it("should set the new provider correctly", () => {
    const provider = {}

    web3Wrapper.setProvider(provider)

    expect(mockedWeb3Lib.setProvider).toHaveBeenCalledWith(provider)
  })
})
