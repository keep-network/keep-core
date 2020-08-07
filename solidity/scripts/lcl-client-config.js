/*
This script is used to update client configuration file with latest deployed contracts
addresses.

Example:
KEEP_CORE_CONFIG_FILE_PATH=~go/src/github.com/keep-network/keep-core/config.toml \
    truffle exec scripts/lcl-client-config.js --network local
*/
const fs = require('fs')
const toml = require('toml')
const tomlify = require('tomlify-j0.4')

const KeepRandomBeaconOperator = artifacts.require('KeepRandomBeaconOperator')
const TokenStaking = artifacts.require('TokenStaking')
const KeepRandomBeaconService = artifacts.require('KeepRandomBeaconService')

module.exports = async function () {
    try {
        const configFilePath = process.env.KEEP_CORE_CONFIG_FILE_PATH

        try {
            await KeepRandomBeaconOperator.deployed()
            await TokenStaking.deployed()
            await KeepRandomBeaconService.deployed()

        } catch (err) {
            console.error('failed to get deployed contracts', err)
            process.exit(1)
        }

        try {
            const fileContent = toml.parse(fs.readFileSync(configFilePath, 'utf8'))

            fileContent.ethereum.ContractAddresses.KeepRandomBeaconOperator = KeepRandomBeaconOperator.address
            fileContent.ethereum.ContractAddresses.TokenStaking = TokenStaking.address
            fileContent.ethereum.ContractAddresses.KeepRandomBeaconService = KeepRandomBeaconService.address

            /*
            tomlify.toToml() writes our Seed/Port values as a float.  The added precision renders our config
            file unreadable by the keep-client as it interprets 3919.0 as a string when it expects an int.
            Here we format the default rendering to write the config file with Seed/Port values as needed.
            */
            let formattedConfigFile = tomlify.toToml(fileContent, {
                space: 2,
                replace: (key, value) => {
                    let result
                    try {
                      result =
                        // We expect the config file to contain arrays, in such case key for
                        // each entry is its' index number. We verify if the key is a string
                        // so we can run the following match check.
                        typeof key === "string" &&
                        // Find keys that match exactly `Port`, `MiningCheckInterval`,
                        // `MaxGasPrice` or end with `MetricsTick`.
                        key.match(
                          /(^Port|^MiningCheckInterval|^MaxGasPrice|MetricsTick)$/
                        )
                          ? value.toFixed(0) // convert float to integer
                          : false // do nothing
                    } catch (err) {
                      console.error(
                        `tomlify replace failed for key ${key} and value ${value} with error: [${err}]`
                      )
                      process.exit(1)
                    }
          
                    return result
                    },
            });

            fs.writeFileSync(configFilePath, formattedConfigFile, (err) => {
                if (err) throw err
            })

            console.log(`keep-core config written to ${configFilePath}`)
        } catch (err) {
            console.error('failed to update keep-core client config', err)
            process.exit(1)
        }

    } catch (err) {
        console.error(err)
        process.exit(1)
    }
    process.exit(0)
}
