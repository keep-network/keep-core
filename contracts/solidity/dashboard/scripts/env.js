const fs = require('fs');
const KeepToken = artifacts.require("./KeepToken.sol");
const TokenStaking = artifacts.require("./TokenStaking.sol");
const TokenGrant = artifacts.require("./TokenGrant.sol");

module.exports = async function() {
  const keepToken = await KeepToken.deployed();
  const tokenStaking = await TokenStaking.deployed();
  const tokenGrant = await TokenGrant.deployed();

  // Write deployed contract addresses into the .env file
  fs.writeFileSync(process.cwd() + "/.env",
  "REACT_APP_TOKEN_ADDRESS=" + keepToken.address +
  "\nREACT_APP_STAKING_ADDRESS=" + tokenStaking.address +
  "\nREACT_APP_TOKENGRANT_ADDRESS=" + tokenGrant.address);
};
