import {
  resolveWeb3Deferred,
  getContractDeploymentBlockNumber,
} from "../../contracts"
import { KEEP_TOKEN_CONTRACT_NAME } from "../../constants/constants"

jest.mock("@keep-network/keep-core/artifacts/KeepToken.json", () => ({
  networks: { "1": { addrress: "0x0", transactionHash: "0x00" } },
}))

describe("Test `getContractDeploymentBlockNumber` function", () => {
  const mockedWeb3 = {
    eth: {
      getTransaction: jest.fn().mockReturnValue({ blockNumber: 50 }),
    },
  }

  beforeAll(() => {
    resolveWeb3Deferred(mockedWeb3)
  })

  it("should fetch the contract deployment block from a node", async () => {
    const result = await getContractDeploymentBlockNumber(
      KEEP_TOKEN_CONTRACT_NAME
    )

    expect(mockedWeb3.eth.getTransaction).toHaveBeenCalledWith("0x00")
    expect(result).toEqual("50")
  })

  it("should return the block number from cache if it exists by contract name", async () => {
    const result = await getContractDeploymentBlockNumber(
      KEEP_TOKEN_CONTRACT_NAME
    )

    const result2 = await getContractDeploymentBlockNumber(
      KEEP_TOKEN_CONTRACT_NAME
    )

    expect(mockedWeb3.eth.getTransaction).toBeCalledTimes(1)
    expect(result).toEqual("50")
    expect(result).toEqual(result2)
  })
})
