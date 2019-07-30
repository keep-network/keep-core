export const sign = async (dataToSign, address) => {
    return web3.utils.toBN(
        await web3.eth.sign(web3.utils.soliditySha3(dataToSign), address)
    ).add(web3.utils.toBN(27)).toString()
}