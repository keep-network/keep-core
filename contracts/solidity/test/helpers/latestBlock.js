// Returns the number of the last mined block
export default async promise => {
    return (await web3.eth.getBlockNumber());
}
