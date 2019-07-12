pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./utils/AddressArrayUtils.sol";


/**
 @dev Interface of sender contract for approveAndCall pattern.
*/
interface tokenSender {
    function approveAndCall(address _spender, uint256 _value, bytes calldata _extraData) external;
}

/**
 @dev Staking contract interface.
*/
interface tokenStaking {
    function initiateUnstake(uint256 _value, address _operator) external;
    function finishUnstake(address _operator) external;
}

/**
 * @title TokenGrant
 * @dev A token grant contract for a specified standard ERC20 token.
 * Has additional functionality to stake/unstake token grants.
 * Tokens are granted to the grantee via vesting scheme and can be
 * released gradually based on the vesting schedule cliff and vesting duration.
 * Optionally grant can be revoked by the token grant creator.
 */
contract TokenGrant {
    using SafeMath for uint256;
    using SafeERC20 for ERC20;
    using BytesLib for bytes;
    using ECDSA for bytes32;
    using AddressArrayUtils for address[];

    event CreatedTokenGrant(uint256 id);
    event ReleasedTokenGrant(uint256 amount);
    event RevokedTokenGrant(uint256 id);

    struct Grant {
        address owner; // Creator of token grant.
        address grantee; // Address to which granted tokens are going to be released.
        bool staked; // Whether the grant is staked.
        bool revoked; // Whether the grant was revoked by the creator.
        bool revocable; // Whether creator of grant can revoke it.
        uint256 amount; // Amount of tokens to be granted.
        uint256 duration; // Duration in seconds of the period in which the granted tokens will vest.
        uint256 start; // Timestamp at which vesting will start.
        uint256 cliff; // Duration in seconds of the cliff after which tokens will begin to vest.
        uint256 released; // Amount that was released to the grantee.
    }

    struct GrantStake {
        uint256 grantId; // Id of the grant.
        address stakingContract; // Staking contract.
        uint256 amount; // Amount of staked tokens.
    }

    uint256 public numGrants;

    ERC20 public token;

    address[] public stakingContracts;

    // Token grants.
    mapping(uint256 => Grant) public grants;

    // Token grants stakes.
    mapping(address => GrantStake) public grantStakes;

    // Mapping of token grant IDs per particular address
    // involved in a grant as a grantee or as a creator.
    mapping(address => uint256[]) public grantIndices;

    // Token grants balances. Sum of all granted tokens to a grantee.
    // This includes granted tokens that are already vested and
    // available to be released to the grantee
    mapping(address => uint256) public balances;

    /**
     * @dev Creates a token grant contract for a provided Standard ERC20 token.
     * @param _tokenAddress address of a token that will be linked to this contract.
     * @param _stakingContract Address of a staking contract that will be linked to this contract.
     */
    constructor(address _tokenAddress, address _stakingContract) public {
        require(_tokenAddress != address(0x0), "Token address can't be zero.");
        token = ERC20(_tokenAddress);
        stakingContracts.push(_stakingContract);
    }

    /**
     * @dev Gets the amount of granted tokens to the specified address.
     * @param _owner The address to query the grants balance of.
     * @return An uint256 representing the grants balance owned by the passed address.
     */
    function totalBalanceOf(address _owner) public view returns (uint256 balance) {
        return balances[_owner];
    }

    /**
     * @dev Gets grant by ID. Returns only basic grant data.
     * If you need vesting schedule for the grant you must call `getGrantVestingSchedule()`
     * This is to avoid Ethereum `Stack too deep` issue described here:
     * https://forum.ethereum.org/discussion/2400/error-stack-too-deep-try-removing-local-variables
     * @param _id ID of the token grant.
     * @return amount, released, staked, revoked.
     */
    function getGrant(uint256 _id) public view returns (uint256 amount, uint256 released, bool staked, bool revoked) {
        return (
            grants[_id].amount,
            grants[_id].released,
            grants[_id].staked,
            grants[_id].revoked
        );
    }

    /**
     * @dev Gets grant vesting schedule by grant ID.
     * @param _id ID of the token grant.
     * @return  owner, duration, start, cliff
     */
    function getGrantVestingSchedule(uint256 _id) public view returns (address owner, uint256 duration, uint256 start, uint256 cliff) {
        return (
            grants[_id].owner,
            grants[_id].duration,
            grants[_id].start,
            grants[_id].cliff
        );
    }

    /**
     * @dev Gets grant ids of the specified address.
     * @param _granteeOrCreator The address to query.
     * @return An uint256 array of grant IDs.
     */
    function getGrants(address _granteeOrCreator) public view returns (uint256[] memory) {
        return grantIndices[_granteeOrCreator];
    }

    /**
     * @notice Creates a token grant with a vesting schedule where balance released to the
     * grantee gradually in a linear fashion until start + duration. By then all
     * of the balance will have vested. You must approve the amount you want to grant
     * by calling approve() method of the token contract first.
     * @dev Transfers token amount from sender to this token grant contract
     * Sender should approve the amount first by calling approve() on the token contract.
     * @param _amount to be granted.
     * @param _grantee address to which granted tokens are going to be released.
     * @param _cliff duration in seconds of the cliff after which tokens will begin to vest.
     * @param _duration duration in seconds of the period in which the tokens will vest.
     * @param _start timestamp at which vesting will start.
     * @param _revocable whether the token grant is revocable or not.
     */
    function grant(
        uint256 _amount,
        address _grantee,
        uint256 _duration,
        uint256 _start,
        uint256 _cliff,
        bool _revocable
    ) public returns (uint256) {
        require(_grantee != address(0), "Grantee address can't be zero.");
        require(_cliff <= _duration, "Vesting cliff duration must be less or equal total vesting duration.");
        require(_amount <= token.balanceOf(msg.sender), "Sender must have enough amount.");

        uint256 id = numGrants++;
        grants[id] = Grant(msg.sender, _grantee, false, false, _revocable, _amount, _duration, _start, _start.add(_cliff), 0);
        
        // Maintain a record to make it easier to query grants by creator.
        grantIndices[msg.sender].push(id);

        // Maintain a record to make it easier to query grants by grantee.
        grantIndices[_grantee].push(id);

        token.safeTransferFrom(msg.sender, address(this), _amount);

        // Maintain a record of the vested amount 
        balances[_grantee] = balances[_grantee].add(_amount);
        emit CreatedTokenGrant(id);
        return id;
    }

    /**
     * @notice Releases Token grant amount to grantee.
     * @dev Transfers vested tokens of the token grant to grantee.
     * @param _id Grant ID.
     */
    function release(uint256 _id) public {
        require(!grants[_id].staked, "Grant must not be staked.");
        uint256 unreleased = unreleasedAmount(_id);
        require(unreleased > 0, "Grant unreleased amount should be greater than zero.");

        // Update released amount.
        grants[_id].released = grants[_id].released.add(unreleased);

        // Update grantee grants balance.
        balances[grants[_id].grantee] = balances[grants[_id].grantee].sub(unreleased);

        // Transfer tokens from this contract balance to the grantee token balance.
        token.safeTransfer(grants[_id].grantee, unreleased);

        emit ReleasedTokenGrant(unreleased);
    }
    
    /**
     * @notice Calculates and returns vested grant amount.
     * @dev Calculates token grant amount that has already vested, 
     * including any tokens that have already been withdrawn by the grantee as well 
     * as any tokens that are available to withdraw but have not yet been withdrawn.
     * @param _id Grant ID.
     */
    function grantedAmount(uint256 _id) public view returns (uint256) {
        uint256 balance = grants[_id].amount;

        if (now < grants[_id].cliff) {
            return 0; // Cliff period is not over.
        } else if (now >= grants[_id].start.add(grants[_id].duration) || grants[_id].revoked) {
            return balance; // Vesting period is finished.
        } else {
            return balance.mul(now.sub(grants[_id].start)).div(grants[_id].duration);
        }
    }

    /**
     * @notice Calculates unreleased granted amount.
     * @dev Calculates the amount that has already vested but hasn't been released yet.
     * @param _id Grant ID.
     */
    function unreleasedAmount(uint256 _id) public view returns (uint256) {
        uint256 released = grants[_id].released;
        return grantedAmount(_id).sub(released);
    }

    /**
     * @notice Allows the creator of the token grant to revoke it. 
     * @dev Granted tokens that are already vested (releasable amount) remain so grantee can still release them
     * the rest are returned to the token grant creator.
     * @param _id Grant ID.
     */
    function revoke(uint256 _id) public {

        require(grants[_id].owner == msg.sender, "Only grant creator can revoke.");
        require(grants[_id].revocable, "Grant must be revocable in the first place.");
        require(!grants[_id].revoked, "Grant must not be already revoked.");
        require(!grants[_id].staked, "Grant must not be staked for staking.");

        uint256 unreleased = unreleasedAmount(_id);
        uint256 refund = grants[_id].amount.sub(unreleased);
        grants[_id].revoked = true;

        // Update grantee's grants balance.
        balances[grants[_id].grantee] = balances[grants[_id].grantee].sub(refund);

        // Transfer tokens from this contract balance to the creator of the token grant.
        token.safeTransfer(grants[_id].owner, refund);
        emit RevokedTokenGrant(_id);
    }

    /**
     * @notice Stake token grant.
     * @dev Stakable token grant amount is the amount of vested tokens minus what user already released from the grant
     * @param _id Grant Id.
     * @param _stakingContract Address of the staking contract.
     * @param _amount Amount to stake.
     * @param _extraData Data for stake delegation. This byte array must have the following values concatenated:
     * Magpie address (20 bytes) where the rewards for participation are sent, operator's ECDSA (65 bytes) signature of
     * the grantee address and ECDSA signature of this contract address.
     */
    function stake(uint256 _id, address _stakingContract, uint256 _amount, bytes memory _extraData) public {
        require(!grants[_id].revocable, "Revocable grants can not be staked.");
        require(grants[_id].grantee == msg.sender, "Only grantee of the grant can stake it.");
        require(
            stakingContracts.contains(_stakingContract),
            "Provided staking contract is not authorized."
        );

        require(_extraData.length == 150, "Stake delegation data must be provided.");
        address operator = keccak256(abi.encodePacked(address(this))).toEthSignedMessageHash().recover(_extraData.slice(20, 65));
        require(
            operator == keccak256(abi.encodePacked(msg.sender)).toEthSignedMessageHash().recover(_extraData.slice(85, 65)),
            "Signer of the grantee doesn't match signer of the grant contract."
        );

        // Calculate available amount. Amount of vested tokens minus what user already released.
        uint256 available = grants[_id].amount.sub(grants[_id].released);
        require(_amount <= available, "Must have available granted amount to stake.");

        // Keep staking record.
        grantStakes[operator] = GrantStake(_id, _stakingContract, _amount);
    
        tokenSender(address(token)).approveAndCall(_stakingContract, _amount, _extraData.slice(0, 85));
    }

    /**
     * @notice Initiate unstake of the token grant.
     * @param _amount Amount to unstake.
     * @param _operator Operator of the stake.
     */
    function initiateUnstake(uint256 _amount, address _operator) public {
        require(
            msg.sender == _operator || msg.sender == grants[grantStakes[_operator].grantId].grantee,
            "Only operator or grantee can initiate unstake."
        );

        tokenStaking(grantStakes[_operator].stakingContract).initiateUnstake(_amount, _operator);
    }

    /**
     * @notice Finish unstake of the token grant.
     * @param _operator Operator of the stake.
     */
    function finishUnstake(address _operator) public {

        tokenStaking(grantStakes[_operator].stakingContract).finishUnstake(_operator);
        delete grantStakes[_operator];
    }
}
