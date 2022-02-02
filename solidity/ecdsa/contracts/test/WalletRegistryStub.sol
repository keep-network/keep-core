pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "../WalletRegistry.sol";
import "../DKGValidator.sol";
import "../libraries/DKG.sol";
import "../libraries/Wallets.sol";

contract WalletRegistryStub is WalletRegistry {
    constructor(
        SortitionPool _sortitionPool,
        IWalletStaking _staking,
        DKGValidator _dkgValidator,
        address _walletOwner
    ) WalletRegistry(_sortitionPool, _staking, _dkgValidator, _walletOwner) {}

    function getDkgData() external view returns (DKG.Data memory) {
        return dkg;
    }
}
