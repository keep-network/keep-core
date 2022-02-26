pragma solidity ^0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@keep-network/sortition-pools/contracts/SortitionPool.sol";
import "../WalletRegistry.sol";
import "../EcdsaDkgValidator.sol";
import "../libraries/EcdsaDkg.sol";
import "../libraries/Wallets.sol";

contract WalletRegistryStub is WalletRegistry {
    constructor(
        SortitionPool _sortitionPool,
        IStaking _staking,
        EcdsaDkgValidator _dkgValidator,
        IRandomBeacon _randomBeacon,
        address _walletOwner
    )
        WalletRegistry(
            _sortitionPool,
            _staking,
            _dkgValidator,
            _randomBeacon,
            _walletOwner
        )
    {}

    function getDkgData() external view returns (EcdsaDkg.Data memory) {
        return dkg;
    }
}
