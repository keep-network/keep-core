pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "../api/IWalletOwner.sol";
import "../WalletRegistry.sol";
import "../EcdsaDkgValidator.sol";
import "../libraries/EcdsaDkg.sol";
import "../libraries/Wallets.sol";

contract WalletRegistryStub is WalletRegistry {
    constructor(
        SortitionPool _sortitionPool,
        IWalletStaking _staking,
        EcdsaDkgValidator _dkgValidator,
        IRandomBeacon _randomBeacon
    ) WalletRegistry(_sortitionPool, _staking, _dkgValidator, _randomBeacon) {}

    function getDkgData() external view returns (EcdsaDkg.Data memory) {
        return dkg;
    }

    // TODO: Use governance update function once it's implemented
    function setMaliciousDkgResultSlashingAmount(
        uint96 newMaliciousDkgResultSlashingAmount
    ) external {
        maliciousDkgResultSlashingAmount = newMaliciousDkgResultSlashingAmount;
    }
}
