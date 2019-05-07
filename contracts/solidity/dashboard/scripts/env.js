const fs = require('fs');
const KeepToken = artifacts.require("./KeepToken.sol");
const StakingProxy = artifacts.require("./StakingProxy.sol");
const Staking = artifacts.require("./Staking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

module.exports = async function() {
  const keepToken = await KeepToken.deployed();
  const stakingProxy = await StakingProxy.deployed();
  const staking = await Staking.deployed();
  const tokenGrant = await TokenGrant.deployed();

  // Write deployed contract addresses into the .env file
  fs.writeFileSync(process.cwd() + "/.env",
  "REACT_APP_TOKEN_ADDRESS=" + keepToken.address +
  "\nREACT_APP_STAKING_PROXY_ADDRESS=" + stakingProxy.address +
  "\nREACT_APP_STAKING_ADDRESS=" + staking.address +
  "\nREACT_APP_TOKENGRANT_ADDRESS=" + tokenGrant.address);
};
