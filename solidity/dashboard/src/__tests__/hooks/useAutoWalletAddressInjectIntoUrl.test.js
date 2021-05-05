import { renderHook } from "@testing-library/react-hooks"
import useWalletAddressFromUrl from "../../hooks/useWalletAddressFromUrl"

const mockCurrentUrl = "/0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756/dashboard"
const mockAddress = "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756"

jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useRouteMatch: jest
    .fn(() => ({
      url: `/${mockAddress}/overview`,
      params: {
        address: mockAddress,
        actualPath: "overview",
      },
    }))
    .mockImplementationOnce(() => ({
      url: `/${mockAddress}/rewards/tbtc`,
      params: {
        address: mockAddress,
        actualPath: "rewards",
      },
    }))
    .mockImplementationOnce(() => ({
      url: `/234234/rewards/tbtc`,
      params: {
        address: "234234",
        actualPath: "rewards",
      },
    })),
  useLocation: jest.fn(() => ({
    pathname: mockCurrentUrl,
  })),
}))

jest.mock("../../components/Routing", () => ({
  pages: [
    {
      route: {
        path: "/overview",
      },
    },
    {
      route: {
        path: "/liquidity",
      },
    },
  ],
}))

describe("Test useWalletAddressFromUrl hook", () => {
  it("returns empty string if path does not exist", () => {
    const result = renderHook(() => useWalletAddressFromUrl())

    expect(result.result.current).toBe("")
  })

  it("returns empty string if wallet address is not a proper eth address", () => {
    const result = renderHook(() => useWalletAddressFromUrl())

    expect(result.result.current).toBe("")
  })

  it("returns wallet address from url properly", () => {
    const result = renderHook(() => useWalletAddressFromUrl())

    expect(result.result.current).toBe(mockAddress)
  })
})
