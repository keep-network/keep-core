import { renderHook } from "@testing-library/react-hooks"
import useExplorerModeConnect from "../../hooks/useAutoConnect"

const mockCurrentUrl = "/dashboard"
const mockHistoryPush = jest.fn()
const mockAddress = "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756"

jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useLocation: jest.fn(() => ({
    pathname: mockCurrentUrl,
  })),
  useHistory: jest.fn(() => ({
    pathname: mockCurrentUrl,
    push: mockHistoryPush,
  })),
}))

jest.mock("../../hooks/useWalletAddressFromUrl", () =>
  jest.fn(() => {
    return ""
  })
)

jest.mock("../../components/WithWeb3Context", () => ({
  useWeb3Context: jest.fn(() => ({
    yourAddress: mockAddress,
    connector: {
      name: "EXPLORER_MODE",
    },
    connectAppWithWallet: jest.fn(),
  })),
}))

describe("Current url without wallet address", () => {
  it("change url when connecting to explorer mode", () => {
    renderHook(() => useExplorerModeConnect())

    expect(mockHistoryPush).toHaveBeenCalledWith({
      pathname: `/${mockAddress}/dashboard`,
    })
  })
})
