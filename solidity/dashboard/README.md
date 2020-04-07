# KEEP Token Dashboard

A react web app to interact with Keep network staking and token grant contracts.
User has the ability to visualize their token and stake balances, token grant unlocking schedules, stake delegate/undelegate their tokens, grant tokens via new unlocking schedules.

## Development setup

### MacOS

* Install [Ganache](http://truffleframework.com/ganache/) and [Brew](https://brew.sh/)
* Install node.js via brew `brew install node`. If you experience errors during `npm install` please try previous versions of node i.e. `brew install node@11`
* Run Ganache and make sure it is using the correct port (8545):
  * Under settings, change port number to 8545

* Deploy contracts demo setup to it:

```
cd keep-core/contracts/solidity
npm install
npm run demo
```

* Go to the dashboard directory `cd dashboard` and run `npm install`

* Run `npm start`

* Open [http://localhost:3000](http://localhost:3000) to view the dashboard.

* Use metamask with `localhost:8545` to use Ganache test network. Import your first Ganache test account into metamask and you should be able to see the demo data.

### Windows

* Install [Ganache](https://github.com/trufflesuite/ganache/releases)
* Install node.js and npm. (https://nodejs.org/en/)
* Run Ganache and make sure it is using the correct port (8545):
  * Under settings, change port number to 8545.
* Start powershell as Administrator.

* Deploy contracts demo setup:

```
cd keep-core/contracts/solidity
npm install
```
* run package.json demo scripts individually:
```
truffle migrate --reset
truffle exec ./scripts/delegate-tokens.js
```
* Go to the dashboard directory `cd dashboard` and run `npm install`

* Run `npm start`

* Open [http://localhost:3000](http://localhost:3000) to view the dashboard.

* Use metamask with `localhost:8545` to use Ganache test network. Import your first Ganache test account into metamask and you should be able to see the demo data.

### Work with contracts deployed locally
* Go to the `solidity` directory then:
  * create symbolic link `artifacts/` to `build/contracts` with `ln -s build/contracts artifacts`.
  * Create a global symlink for a `@keep-network/keep-core` dependency with `npm link`
* Go to the `dashboard` directory and tell the application to use the global symlink with `npm link @keep-network/keep-core`
* When you don’t want to use the local version of `@keep-network/keep-core` anymore, delete the symlink with `npm uninstall --no-save some-dep && npm install`

## Internal testnet

A new version of staking dApp is automatically deployed to `keep-dev` internal testnet after each `master` merge. dApp can be accessed by VPN at `http://keep-dapp-staking.default.svc.cluster.local/` and requires an initial setup in MetaMask before the first use. All the setup described below has to be done only one time. 

### MetaMask extension setup

MetaMask is a web browser extension allowing to interact with Ethereum-enabled distributed applications (dApps). MetaMask is available for Chrome, Firefox, and Opera desktop browsers. To install MetaMask, please go through the extension installation process individual for your web browser used.

Before the MetaMask can be used for the first time, it requires an initial setup when the user creates a new Wallet. It includes setting up a wallet password, accepting terms of use, creating and confirming backup phrase.

### MetaMask network configuration

Before MetaMask can be used with `keep-dev` testnet for the first time, it needs to know what `keep-dev` is. This process includes setting up a new network:

1. Make sure you are connected to `keep-dev` testnet via VPN
2. Expand the list of networks and click on `Custom RPC`
3. Set `Network Name` to `keep-dev`
4. Set `New RPC URL` to `http://eth-tx-node.default.svc.cluster.local:8545`
5. Set `ChainID` to `1101`
6. Click `Save`

### MetaMask KEEP token owner account import
On `keep-dev`, account `0x0f0977c4161a371b5e5ee6a8f43eb798cd1ae1db` is the owner of contracts including KEEP ERC20 token contract. This account can be used to create token grants and delegate stake to operators. Grantees of tokens can also stake-delegate their grants.

To use this account in the dApp, it needs to be imported from [the JSON file](https://github.com/keep-network/keep-core/blob/master/private-testnet/keyfiles/UTC--2019-03-27T19-05-16.429364100Z--0f0977c4161a371b5e5ee6a8f43eb798cd1ae1db) secured by a [password](https://github.com/keep-network/keep-core/blob/master/private-testnet/eth-account-password.txt).

1. Download the account JSON file
2. Expand the list of accounts and click on `Import Account`
3. Select `JSON file` type
4. Click `Browse` and point MetaMask to the previously downloaded account JSON file
5. Copy-paste the password from the referenced password file
6. Click `Import`
