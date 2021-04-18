const useFinalPath = (to, walletAddress) => {
  return walletAddress ? `/${walletAddress}${to}` : to
}

export default useFinalPath
