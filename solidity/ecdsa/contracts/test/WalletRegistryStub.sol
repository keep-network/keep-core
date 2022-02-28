pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "@keep-network/random-beacon/contracts/ReimbursementPool.sol";
import "../WalletRegistry.sol";
import "../DKGValidator.sol";
import "../libraries/DKG.sol";
import "../libraries/Wallets.sol";

contract WalletRegistryStub is WalletRegistry {
    constructor(
        SortitionPool _sortitionPool,
        IWalletStaking _staking,
        DKGValidator _dkgValidator,
        address _walletOwner,
        ReimbursementPool reimbursementPool
    ) WalletRegistry(_sortitionPool, _staking, _dkgValidator, _walletOwner, reimbursementPool) {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }

    // TODO: Use governance update function once it's implemented
    function setMaliciousDkgResultSlashingAmount(
        uint96 newMaliciousDkgResultSlashingAmount
    ) external {
        maliciousDkgResultSlashingAmount = newMaliciousDkgResultSlashingAmount;
    }
}
