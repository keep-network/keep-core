// Returns the number of the last mined block
const {web3} = require("@openzeppelin/test-environment")

module.exports = async promise => {
    return (await web3.eth.getBlockNumber());
}
