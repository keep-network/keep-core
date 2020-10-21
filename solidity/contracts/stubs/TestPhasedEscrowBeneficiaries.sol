pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

// Simple beneficiary that does nothing when notified that it has received
// tokens.
contract TestSimpleBeneficiary {
    function __escrowSentTokens(uint256 amount) external {}
}

contract TestLPTokenWrapper {
    using SafeMath for uint256;
    using SafeERC20 for IERC20;

    IERC20 public token;

    constructor(IERC20 _token) public {
        token = _token;
    }

    uint256 private _totalSupply;
    mapping(address => uint256) private _balances;

    function totalSupply() public view returns (uint256) {
        return _totalSupply;
    }

    function balanceOf(address account) public view returns (uint256) {
        return _balances[account];
    }

    function stake(uint256 amount) public {
        _totalSupply = _totalSupply.add(amount);
        _balances[msg.sender] = _balances[msg.sender].add(amount);
        token.safeTransferFrom(msg.sender, address(this), amount);
    }
}

contract TestCurveRewards is TestLPTokenWrapper {
    event Staked(address indexed user, uint256 amount);

    constructor(IERC20 _token) public TestLPTokenWrapper(_token) {}

    // stake visibility is public as overriding LPTokenWrapper's stake() function
    function stake(uint256 amount) public {
        require(amount > 0, "Cannot stake 0");
        super.stake(amount);
        emit Staked(msg.sender, amount);
    }
}
