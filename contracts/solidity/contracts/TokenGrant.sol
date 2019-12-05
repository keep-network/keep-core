pragma solidity ^0.5.4;

import "openzeppelin-solidity/contracts/token/ERC20/ERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";
import "openzeppelin-solidity/contracts/cryptography/ECDSA.sol";
import "solidity-bytes-utils/contracts/BytesLib.sol";
import "./utils/AddressArrayUtils.sol";
import "./TokenStaking.sol";


/**
 @dev Interface of sender contract for approveAndCall pattern.
*/
interface tokenSender {
    function approveAndCall(address _spender, uint256 _value, bytes calldata _extraData) external;
}

/**
 * @title TokenGrant
 * @dev A token grant contract for a specified standard ERC20 token.
 * Has additional functionality to stake/unstake token grants.
 * Tokens are granted to the grantee via vesting scheme and can be
 * withdrawn gradually based on the vesting schedule cliff and vesting duration.
 * Optionally grant can be revoked by the token grant manager.
 */
contract TokenGrant {
    using SafeMath for uint256;
    using SafeERC20 for ERC20;
    using BytesLib for bytes;
    using ECDSA for bytes32;
    using AddressArrayUtils for address[];

    event CreatedTokenGrant(uint256 id);
    event WithdrawnTokenGrant(uint256 amount);
    event RevokedTokenGrant(uint256 id);

    struct Grant {
        address grantManager; // Token grant manager.
        address grantee; // Address to which granted tokens are going to be withdrawn.
        bool revoked; // Whether the grant was revoked by the grant manager.
        bool revocable; // Whether grant manager can revoke the grant.
        uint256 amount; // Amount of tokens to be granted.
        uint256 duration; // Duration in seconds of the period in which the granted tokens will vest.
        uint256 start; // Timestamp at which vesting will start.
        uint256 cliff; // Duration in seconds of the cliff after which tokens will begin to vest.
        uint256 withdrawn; // Amount that was withdrawn to the grantee.
        uint256 staked; // Amount that was staked by the grantee.
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
    // involved in a grant as a grantee or as a grant manager.
    mapping(address => uint256[]) public grantIndices;

    // Token grants balances. Sum of all granted tokens to a grantee.
    // This includes granted tokens that are already vested and
    // available to be withdrawn to the grantee
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
    function balanceOf(address _owner) public view returns (uint256 balance) {
        return balances[_owner];
    }

    /**
     * @dev Gets the stake balance of the specified address.
     * @param _address The address to query the balance of.
     * @return An uint256 representing the amount staked by the passed address.
     */
    function stakeBalanceOf(address _address) public view returns (uint256 balance) {
        for (uint i = 0; i < grantIndices[_address].length; i++) {
            uint256 id = grantIndices[_address][i];
            balance += grants[id].staked;
        }
        return balance;
    }

    /**
     * @dev Gets grant by ID. Returns only basic grant data.
     * If you need vesting schedule for the grant you must call `getGrantVestingSchedule()`
     * This is to avoid Ethereum `Stack too deep` issue described here:
     * https://forum.ethereum.org/discussion/2400/error-stack-too-deep-try-removing-local-variables
     * @param _id ID of the token grant.
     * @return amount, withdrawn, staked, revoked.
     */
    function getGrant(uint256 _id) public view returns (uint256 amount, uint256 withdrawn, uint256 staked, bool revoked) {
        return (
            grants[_id].amount,
            grants[_id].withdrawn,
            grants[_id].staked,
            grants[_id].revoked
        );
    }

    /**
     * @dev Gets grant vesting schedule by grant ID.
     * @param _id ID of the token grant.
     * @return grantManager, duration, start, cliff
     */
    function getGrantVestingSchedule(uint256 _id) public view returns (address grantManager, uint256 duration, uint256 start, uint256 cliff) {
        return (
            grants[_id].grantManager,
            grants[_id].duration,
            grants[_id].start,
            grants[_id].cliff
        );
    }

    /**
     * @dev Gets grant ids of the specified address.
     * @param _granteeOrGrantManager The address to query.
     * @return An uint256 array of grant IDs.
     */
    function getGrants(address _granteeOrGrantManager) public view returns (uint256[] memory) {
        return grantIndices[_granteeOrGrantManager];
    }

    /**
     * @notice Receives approval of token transfer and creates a token grant with a vesting
     * schedule where balance withdrawn to the grantee gradually in a linear fashion until
     * start + duration. By then all of the balance will have vested.
     * @param _from The owner of the tokens who approved them to transfer.
     * @param _amount Approved amount for the transfer to create token grant.
     * @param _token Token contract address.
     * @param _extraData This byte array must have the following values concatenated:
     * grantee (20 bytes) Address of the grantee.
     * cliff (32 bytes) Duration in seconds of the cliff after which tokens will begin to vest.
     * start (32 bytes) Timestamp at which vesting will start.
     * revocable (1 byte) Whether the token grant is revocable or not (1 or 0).
     */
    function receiveApproval(address _from, uint256 _amount, address _token, bytes memory _extraData) public {
        require(ERC20(_token) == token, "Token contract must be the same one linked to this contract.");
        require(_amount <= token.balanceOf(_from), "Sender must have enough amount.");

        address _grantee = _extraData.toAddress(0);
        uint256 _duration = _extraData.toUint(20);
        uint256 _start = _extraData.toUint(52);
        uint256 _cliff = _extraData.toUint(84);
        
        require(_grantee != address(0), "Grantee address can't be zero.");
        require(_cliff <= _duration, "Vesting cliff duration must be less or equal total vesting duration.");

        bool _revocable;
        if (_extraData.slice(116, 1)[0] == 0x01) {
            _revocable = true;
        } 

        uint256 id = numGrants++;
        grants[id] = Grant(_from, _grantee, false, _revocable, _amount, _duration, _start, _start.add(_cliff), 0, 0);

        // Maintain a record to make it easier to query grants by grant manager.
        grantIndices[_from].push(id);

        // Maintain a record to make it easier to query grants by grantee.
        grantIndices[_grantee].push(id);

        token.safeTransferFrom(_from, address(this), _amount);

        // Maintain a record of the vested amount
        balances[_grantee] = balances[_grantee].add(_amount);
        emit CreatedTokenGrant(id);
    }

    /**
     * @notice Withdraws Token grant amount to grantee.
     * @dev Transfers vested tokens of the token grant to grantee.
     * @param _id Grant ID.
     */
    function withdraw(uint256 _id) public {
        uint256 amount = withdrawable(_id);
        require(amount > 0, "Grant available to withdraw amount should be greater than zero.");

        // Update withdrawn amount.
        grants[_id].withdrawn = grants[_id].withdrawn.add(amount);

        // Update grantee grants balance.
        balances[grants[_id].grantee] = balances[grants[_id].grantee].sub(amount);

        // Transfer tokens from this contract balance to the grantee token balance.
        token.safeTransfer(grants[_id].grantee, amount);

        emit WithdrawnTokenGrant(amount);
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
     * @notice Calculates withdrawable granted amount.
     * @dev Calculates the amount that has already vested but hasn't been withdrawn yet.
     * @param _id Grant ID.
     */
    function withdrawable(uint256 _id) public view returns (uint256) {
        return grantedAmount(_id).sub(grants[_id].withdrawn).sub(grants[_id].staked);
    }

    /**
     * @notice Allows the grant manager to revoke the grant. 
     * @dev Granted tokens that are already vested (releasable amount) remain so grantee can still withdraw them
     * the rest are returned to the token grant manager.
     * @param _id Grant ID.
     */
    function revoke(uint256 _id) public {

        require(grants[_id].grantManager == msg.sender, "Only grant manager can revoke.");
        require(grants[_id].revocable, "Grant must be revocable in the first place.");
        require(!grants[_id].revoked, "Grant must not be already revoked.");

        uint256 amount = withdrawable(_id);
        uint256 refund = grants[_id].amount.sub(amount);
        grants[_id].revoked = true;

        // Update grantee's grants balance.
        balances[grants[_id].grantee] = balances[grants[_id].grantee].sub(refund);

        // Transfer tokens from this contract balance to the token grant manager.
        token.safeTransfer(grants[_id].grantManager, refund);
        emit RevokedTokenGrant(_id);
    }

    /**
     * @notice Stake token grant.
     * @dev Stakable token grant amount is the amount of vested tokens minus what user already withdrawn from the grant
     * @param _id Grant Id.
     * @param _stakingContract Address of the staking contract.
     * @param _amount Amount to stake.
     * @param _extraData Data for stake delegation. This byte array must have the following values concatenated:
     * Magpie address (20 bytes) where the rewards for participation are sent, operator's ECDSA (65 bytes) signature of
     * the grantee address and ECDSA (65 bytes) signature of this contract address.
     */
    function stake(uint256 _id, address _stakingContract, uint256 _amount, bytes memory _extraData) public {
        require(!grants[_id].revocable, "Revocable grants can not be staked.");
        require(grants[_id].grantee == msg.sender, "Only grantee of the grant can stake it.");
        require(
            stakingContracts.contains(_stakingContract),
            "Provided staking contract is not authorized."
        );

        // Expecting 150 bytes _extraData for stake delegation
        // 20 bytes address + two 65 bytes ECDSA signatures
        require(_extraData.length == 150, "Stake delegation data must be provided.");
        address operator = keccak256(abi.encodePacked(address(this))).toEthSignedMessageHash().recover(_extraData.slice(20, 65));
        require(
            operator == keccak256(abi.encodePacked(msg.sender)).toEthSignedMessageHash().recover(_extraData.slice(85, 65)),
            "Signer of the grantee doesn't match signer of the grant contract."
        );

        // Calculate available amount. Amount of vested tokens minus what user already withdrawn and staked.
        uint256 available = grants[_id].amount.sub(grants[_id].withdrawn).sub(grants[_id].staked);
        require(_amount <= available, "Must have available granted amount to stake.");

        // Keep staking record.
        grantStakes[operator] = GrantStake(_id, _stakingContract, _amount);
        grants[_id].staked += _amount;

        // Staking contract expects 85 bytes _extraData for stake delegation
        // 20 bytes address + 65 bytes ECDSA signature
        tokenSender(address(token)).approveAndCall(_stakingContract, _amount, _extraData.slice(0, 85));
    }

    /**
     * @notice Initiate unstake of the token grant.
     * @param _operator Operator of the stake.
     */
    function initiateUnstake(address _operator) public {
        uint256 grantId = grantStakes[_operator].grantId;
        require(
            msg.sender == _operator || msg.sender == grants[grantId].grantee,
            "Only operator or grantee can initiate unstake."
        );

        TokenStaking(grantStakes[_operator].stakingContract).initiateUnstake(grantStakes[_operator].amount, _operator);
    }

    /**
     * @notice Finish unstake of the token grant.
     * @param _operator Operator of the stake.
     */
    function finishUnstake(address _operator) public {
        uint256 grantId = grantStakes[_operator].grantId;
        grants[grantId].staked = grants[grantId].staked.sub(grantStakes[_operator].amount);

        TokenStaking(grantStakes[_operator].stakingContract).finishUnstake(_operator);
        delete grantStakes[_operator];
    }
}
