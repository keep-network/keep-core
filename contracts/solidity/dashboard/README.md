# KEEP Token Dashboard

A react web app to interact with Keep network staking and token grant contracts.
User has the ability to visualise their token and stake balances, token grant vesting schedules, stake their tokens and initiate/finish unstake, grant tokens via new vesting schedules.

### Prerequisites

* The dashboard requires addresses of the contracts to be listed in the `.env` file similar to this:

```
REACT_APP_TOKEN_ADDRESS=0x74e3fc764c2474f25369b9d021b7f92e8441a2dc
REACT_APP_STAKING_ADDRESS=0x8e4c131b37383e431b9cd0635d3cf9f3f628edae
REACT_APP_TOKENGRANT_ADDRESS=0xb9b7e0cb2edf5ea031c8b297a5a1fa20379b6a0a
```

This is done automatically as part of the `npm start` script, just make sure the contracts are deployed to the network the dashboard is using. You can also use the demo script in the parent folder to add demo stakes and token grants `cd ../ && npm run demo`

### Setup and run

* Run `npm start` from the project directory

* Open [http://localhost:3000](http://localhost:3000) to view the dashboard.
