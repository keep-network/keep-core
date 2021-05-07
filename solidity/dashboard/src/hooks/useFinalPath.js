import useWalletAddressFromUrl from "./useWalletAddressFromUrl"

const useFinalPath = (to) => {
  const walletAddress = useWalletAddressFromUrl()
  return walletAddress ? `/${walletAddress}${to}` : to
}

export default useFinalPath
