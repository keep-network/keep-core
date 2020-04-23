usePlugin("@nomiclabs/buidler-truffle5")
usePlugin("@nomiclabs/buidler-etherscan")

module.exports = {
  defaultNetwork: "local",
  networks: {
    local: {
      url: "http://localhost:8545",
    },
    ropsten: {
      url: "https://ropsten.infura.io/v3/59fb36a36fa4474b890c13dd30038be5",
      chainId: 3,
      from: "0x923C5Dbf353e99394A21Aa7B67F3327Ca111C67D",

      // Gas limit for every transaction.
      gas: "auto",
      // Gas price.
      gasPrice: "auto",
      // A number used to multiply the results of gas estimation to give it some
      // slack due to the uncertainty of the estimation process
      gasMultiplier: 1,

      accounts: [process.env.CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY],
    },
  },
  solc: { version: "0.5.17" },
  etherscan: {
    url: "https://api-ropsten.etherscan.io/api",
    apiKey: "YOUR_API_KEY",
  },
}
