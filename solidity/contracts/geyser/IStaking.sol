/**
 This code is copied from:
 https://github.com/ampleforth/token-geyser/tree/d8352f62a0432494c39416d090e68582e13b2b22/contracts
 */
pragma solidity 0.5.17;

/**
 * @title Staking interface, as defined by EIP-900.
 * @dev https://github.com/ethereum/EIPs/blob/master/EIPS/eip-900.md
 */
contract IStaking {
    event Staked(
        address indexed user,
        uint256 amount,
        uint256 total,
        bytes data
    );
    event Unstaked(
        address indexed user,
        uint256 amount,
        uint256 total,
        bytes data
    );

    function stake(uint256 amount, bytes calldata data) external;

    function stakeFor(
        address user,
        uint256 amount,
        bytes calldata data
    ) external;

    function unstake(uint256 amount, bytes calldata data) external;

    function token() external view returns (address);

    /**
     * @return False. This application does not support staking history.
     */
    function supportsHistory() external pure returns (bool) {
        return false;
    }

    function totalStakedFor(address addr) public view returns (uint256);

    function totalStaked() public view returns (uint256);
}
