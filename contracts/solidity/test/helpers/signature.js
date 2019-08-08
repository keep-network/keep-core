export const sign = async (dataToSign, address) => {
    return '0x' + Buffer.from(web3.utils.toBN(
        await web3.eth.sign(dataToSign, address)
    ).add(web3.utils.toBN(27)).toBuffer()).toString('hex')
}
