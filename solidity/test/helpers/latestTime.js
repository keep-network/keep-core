const {web3} = require("@openzeppelin/test-environment")
// Returns the time of the last mined block in seconds

module.exports = async promise => {
  return (await web3.eth.getBlock('latest')).timestamp;
}
