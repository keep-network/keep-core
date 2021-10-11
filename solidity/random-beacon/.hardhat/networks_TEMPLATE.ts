// This is a template file. The file should be renamed to `networks.ts` and updated
// with valid data. The properties will overwrite the ones defined in hardhat
// project config file if hardhat.config.ts file contains `localNetworksConfig` property
// pointing to this file.

import { LocalNetworksConfig } from "@keep-network/hardhat-local-networks-config"

const config: LocalNetworksConfig = {
  networks: {
    ropsten: {
      url: "url not set",
      from: "address not set",
      accounts: ["private key not set"],
      tags: ["tenderly"],
    },
    mainnet: {
      url: "url not set",
      from: "address not set",
      accounts: ["private key not set"],
      tags: ["tenderly"],
    },
  },
}

module.exports = config
