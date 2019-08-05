import mineBlocks from './mineBlocks';

export default async function expireGroup(
    operatorContract,
    groupIndex,
) {
    let groupRegistrationBlock = await operatorContract.getGroupRegistrationBlockHeight(groupIndex);
    let groupActiveTime = await operatorContract.groupActiveTime()

    let currentBlock = await web3.eth.getBlockNumber();

    // If current block is larger than group registration block by group active time then
    // it is not necessary to mine any blocks cause the group is already expired
    if (currentBlock - groupRegistrationBlock <= groupActiveTime) {
      await mineBlocks(groupActiveTime - (currentBlock - groupRegistrationBlock));
    }
}