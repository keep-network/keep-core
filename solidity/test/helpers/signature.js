const {web3} = require("@openzeppelin/test-environment")

const sign = async (dataToSign, address) => {
    // go-ethereum/crypto produces signature with v={0, 1} and we need to add
    // 27 to v-part (signature[64]) to conform wtih the on-chain signature 
    // validation code that accepts v={27, 28} as specified in the
    // Appendix F of the Ethereum Yellow Paper 
    // https://ethereum.github.io/yellowpaper/paper.pdf
    return '0x' + web3.utils.toBN(
        await web3.eth.sign(dataToSign, address)
    ).add(web3.utils.toBN(27)).toBuffer('be', 65).toString('hex')
}

module.exports = sign
