# KEEP Token Dashboard

A react web app to interact with Keep network staking and token grant contracts.
User has the ability to visualise their token and stake balances, token grant vesting schedules, stake their tokens and initiate/finish unstake, grant tokens via new vesting schedules.

## Development setup

#### MacOS

* Install [Ganache](http://truffleframework.com/ganache/) and [Brew](https://brew.sh/)
* Install node.js via brew `brew install node`
* Run Ganache and depoloy contracts demo setup to it:

```
cd keep-core/contracts/solidity
npm install
npm run demo
```

* Go to the dashboard directory `cd dashboard` and run `npm install`

* Run `npm start`

* Open [http://localhost:3000](http://localhost:3000) to view the dashboard.

* Use metamask with `localhost:8545` to use Ganache test network. Import your first Ganache test account into metamask and you should be able to see the demo data.
