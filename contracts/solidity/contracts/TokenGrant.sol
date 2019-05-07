pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";


/**
 * @title TokenGrant
 * @dev A token grant contract for a specified standard ERC20 token.
 * Tokens are granted to the beneficiary via vesting scheme and can be
 * released gradually based on the vesting schedule cliff and vesting duration.
 * Optionally grant can be revoked by the token grant creator.
 */
contract TokenGrant {
    using SafeMath for uint256;
    using SafeERC20 for ERC20;

    ERC20 public token;

    event CreatedTokenGrant(uint256 id);
    event ReleasedTokenGrant(uint256 amount);
    event RevokedTokenGrant(uint256 id);

    struct Grant {
        address owner; // Creator of token grant.
        address beneficiary; // Address to which granted tokens are going to be released.
        bool staked; // Whether the grant is staked.
        bool revoked; // Whether the grant was revoked by the creator.
        bool revocable; // Whether creator of grant can revoke it.
        uint256 amount; // Amount of tokens to be granted.
        uint256 duration; // Duration in seconds of the period in which the granted tokens will vest.
        uint256 start; // Timestamp at which vesting will start.
        uint256 cliff; // Duration in seconds of the cliff after which tokens will begin to vest.
        uint256 released; // Amount that was released to the beneficiary.
    }

    uint256 public numGrants;

    // Token grants.
    mapping(uint256 => Grant) public grants;

    // Mapping of token grant IDs per particular address
    // involved in a grant as a beneficiary or as a creator.
    mapping(address => uint256[]) public grantIndices;

    // Token grants balances. Sum of all granted tokens to a beneficiary.
    // This includes granted tokens that are already vested and
    // available to be released to the beneficiary
    mapping(address => uint256) public balances;


    /**
     * @dev Creates a token grant contract for a provided Standard ERC20 token.
     * @param _tokenAddress address of a token that will be linked to this contract.
     */
    constructor(address _tokenAddress) public {
        require(_tokenAddress != address(0x0), "Token address can't be zero.");
        token = ERC20(_tokenAddress);
    }

    /**
     * @dev Gets the amount of granted tokens to the specified address.
     * @param _owner The address to query the grants balance of.
     * @return An uint256 representing the grants balance owned by the passed address.
     */
    function balanceOf(address _owner) public view returns (uint256 balance) {
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
    function getGrant(uint256 _id) public view returns (uint256, uint256, bool, bool) {
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
    function getGrantVestingSchedule(uint256 _id) public view returns (address, uint256, uint256, uint256) {
        return (
            grants[_id].owner,
            grants[_id].duration,
            grants[_id].start,
            grants[_id].cliff
        );
    }

    /**
     * @dev Gets grant ids of the specified address.
     * @param _beneficiaryOrCreator The address to query.
     * @return An uint256 array of grant IDs.
     */
    function getGrants(address _beneficiaryOrCreator) public view returns (uint256[] memory) {
        return grantIndices[_beneficiaryOrCreator];
    }

    /**
     * @notice Creates a token grant with a vesting schedule where balance released to the
     * beneficiary gradually in a linear fashion until start + duration. By then all
     * of the balance will have vested. You must approve the amount you want to grant
     * by calling approve() method of the token contract first.
     * @dev Transfers token amount from sender to this token grant contract
     * Sender should approve the amount first by calling approve() on the token contract.
     * @param _amount to be granted.
     * @param _beneficiary address to which granted tokens are going to be released.
     * @param _cliff duration in seconds of the cliff after which tokens will begin to vest.
     * @param _duration duration in seconds of the period in which the tokens will vest.
     * @param _start timestamp at which vesting will start.
     * @param _revocable whether the token grant is revocable or not.
     */
    function grant(
        uint256 _amount,
        address _beneficiary,
        uint256 _duration,
        uint256 _start,
        uint256 _cliff,
        bool _revocable
    ) public returns (uint256) {
        require(_beneficiary != address(0), "Beneficiary address can't be zero.");
        require(_cliff <= _duration, "Vesting cliff duration must be less or equal total vesting duration.");
        require(_amount <= token.balanceOf(msg.sender), "Sender must have enough amount.");

        uint256 id = numGrants++;
        grants[id] = Grant(msg.sender, _beneficiary, false, false, _revocable, _amount, _duration, _start, _start.add(_cliff), 0);
        
        // Maintain a record to make it easier to query grants by creator.
        grantIndices[msg.sender].push(id);

        // Maintain a record to make it easier to query grants by beneficiary.
        grantIndices[_beneficiary].push(id);

        token.transferFrom(msg.sender, address(this), _amount);

        // Maintain a record of the vested amount 
        balances[_beneficiary] = balances[_beneficiary].add(_amount);
        emit CreatedTokenGrant(id);
        return id;
    }

    /**
     * @notice Releases Token grant amount to beneficiary.
     * @dev Transfers vested tokens of the token grant to beneficiary.
     * @param _id Grant ID.
     */
    function release(uint256 _id) public {
        require(!grants[_id].staked, "Grant must not be staked.");
        uint256 unreleased = unreleasedAmount(_id);
        require(unreleased > 0, "Grant unreleased amount should be greater than zero.");

        // Update released amount.
        grants[_id].released = grants[_id].released.add(unreleased);

        // Update beneficiary grants balance.
        balances[grants[_id].beneficiary] = balances[grants[_id].beneficiary].sub(unreleased);

        // Transfer tokens from this contract balance to the beneficiary token balance.
        token.safeTransfer(grants[_id].beneficiary, unreleased);

        emit ReleasedTokenGrant(unreleased);
    }
    
    /**
     * @notice Calculates and returns vested grant amount.
     * @dev Calculates token grant amount that has already vested, 
     * including any tokens that have already been withdrawn by the beneficiary as well 
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
     * @dev Granted tokens that are already vested (releasable amount) remain so beneficiary can still release them
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

        // Update beneficiary's grants balance.
        balances[grants[_id].beneficiary] = balances[grants[_id].beneficiary].sub(refund);

        // Transfer tokens from this contract balance to the creator of the token grant.
        token.safeTransfer(grants[_id].owner, refund);
        emit RevokedTokenGrant(_id);
    }

    /**
     * @dev Allows authorized address/contract to transfer token grant to a new address.
     * @param _id Id of the grant to transfer the amount from.
     * @param _beneficiary Beneficiary of the new grant to receive the amount.
     * @param _amount Amount to transfer.
     * @return Returns id of a new token grant with the transferred amount
     */
    function transfer(uint256 _id, address _beneficiary, uint256 _amount) public returns (uint256) {

        // TODO: allowance/approval authorization for transfers
        // i.e. require(authorizedManager == msg.sender, "Only address authorized by grant beneficiary can invoke transfer.");

        uint256 available = grants[_id].amount.sub(grants[_id].released);
        require(available >= _amount, "Must have available granted amount to transfer.");

        uint256 unreleased = unreleasedAmount(_id);

        // Remove amount from the source grant.
        grants[_id].amount = grants[_id].amount.sub(_amount);

        // If there are enough unreleased tokens to cover the required amount.
        if (unreleased >= _amount) {

            // Calculate the remaining unreleased amount.
            uint256 remaining = unreleased.sub(_amount);

            // Adjust start time to maintain remaining unreleased amount.
            grants[_id].start = now.sub(remaining.mul(grants[_id].duration).div(unreleased + remaining));

        } else {
            // Reset vesting schedule so no unreleased amount is available.
            grants[_id].start = now;
            grants[_id].amount = grants[_id].amount.sub(grants[_id].released);
            grants[_id].released = 0;
        }

        // Destination token grant.
        uint256 newGrantId = numGrants++;
        grants[newGrantId] = Grant(
            grants[_id].owner,
            _beneficiary,
            grants[_id].staked,
            false,
            false,
            _amount,
            grants[_id].duration,
            // Adjust start time to keep unreleased amount the same as in the source grant.
            now.sub(unreleased.mul(grants[_id].duration).div(_amount)),
            grants[_id].cliff,
            0
        );

        emit CreatedTokenGrant(newGrantId);
        return newGrantId;
    }
}
