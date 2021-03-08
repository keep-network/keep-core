pragma solidity 0.5.17;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import "openzeppelin-solidity/contracts/token/ERC20/IERC20.sol";
import "openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";
import "openzeppelin-solidity/contracts/cryptography/MerkleProof.sol";

/// @title Token Distributor
/// @notice This contract can be used to distribute ERC20 tokens
/// @dev This contract is based on the Uniswap's Merkle Distributor
/// https://github.com/Uniswap/merkle-distributor with some modifications:
/// - added 'allocate()' function that will be called to allocate tokens for a
///   merkle root,
/// - added a possibility for recipient to redirect tokens withdrawal to another
///   address by providing a signature over that address.
contract TokenDistributor is Ownable {
    using SafeERC20 for IERC20;

    IERC20 public token;

    bytes32 public merkleRoot;

    mapping(uint256 => uint256) private claimedBitMap;

    event TokensAllocated(bytes32 merkleRoot, uint256 amount);

    event TokensClaimed(
        uint256 indexed index,
        address indexed recipient,
        address indexed destination,
        uint256 amount
    );

    constructor(IERC20 _token) public {
        token = _token;
    }

    // In the claim function, you need to provide Ethereum address and a signed
    // address of token recipient (the signature by Ethereum address from a
    // merkle tree). We'll validate the signature and see how many tokens should
    // be claimable by that address based on the information in Merkle tree.
    // TODO:Update docs
    function claim(
        address _recipient,
        address _destination,
        uint8 _v,
        bytes32 _r,
        bytes32 _s,
        uint256 _index,
        uint256 _amount,
        bytes32[] calldata _merkleProof
    ) external {
        require(_recipient != address(0), "recipient address cannot be zero");
        require(
            _destination != address(0),
            "destination address cannot be zero"
        );
        require(merkleRoot != "", "tokens were not allocated yet");
        require(!isClaimed(_index), "tokens already claimed");

        // Verify the signature over destination address.
        require(
            _recipient ==
                recoverSignerAddress(
                    abi.encodePacked(_destination),
                    _v,
                    _r,
                    _s
                ),
            "invalid signature of destination address"
        );

        // Verify the merkle proof.
        bytes32 node = keccak256(abi.encodePacked(_index, _recipient, _amount));
        require(
            MerkleProof.verify(_merkleProof, merkleRoot, node),
            "invalid proof"
        );

        // Mark it claimed and send the token.
        setClaimed(_index);

        token.safeTransfer(_destination, _amount);

        emit TokensClaimed(_index, _recipient, _destination, _amount);
    }

    /// Allocates amount of tokens for the merkle root.
    /// @param _merkleRoot The merkle root.
    /// @param _amount The amount of tokens allocated for the merkle root.
    function allocate(bytes32 _merkleRoot, uint256 _amount) public onlyOwner {
        require(merkleRoot == "", "tokens were already allocated");
        require(_merkleRoot != "", "merkle root cannot be empty");

        token.safeTransferFrom(msg.sender, address(this), _amount);

        merkleRoot = _merkleRoot;

        emit TokensAllocated(_merkleRoot, _amount);
    }

    function isClaimed(uint256 _index) public view returns (bool) {
        uint256 claimedWordIndex = _index / 256;
        uint256 claimedBitIndex = _index % 256;
        uint256 claimedWord = claimedBitMap[claimedWordIndex];
        uint256 mask = (1 << claimedBitIndex);
        return claimedWord & mask == mask;
    }

    function recoverSignerAddress(
        bytes memory _message,
        uint8 _v,
        bytes32 _r,
        bytes32 _s
    ) internal pure returns (address) {
        // Validate `s` and `v` values for a malleability concern described in EIP-2.
        // Only signatures with `s` value in the lower half of the secp256k1
        // curve's order and `v` value of 27 or 28 are considered valid.
        require(
            uint256(_s) <=
                0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF5D576E7357A4501DDFE92F46681B20A0,
            "Invalid signature 's' value"
        );
        require(_v == 27 || _v == 28, "Invalid signature 'v' value");

        bytes32 digest = keccak256(_message);
        bytes32 prefixedDigest =
            keccak256(
                abi.encodePacked("\x19Ethereum Signed Message:\n32", digest)
            );

        return ecrecover(prefixedDigest, _v, _r, _s);
    }

    function setClaimed(uint256 _index) private {
        uint256 claimedWordIndex = _index / 256;
        uint256 claimedBitIndex = _index % 256;
        claimedBitMap[claimedWordIndex] =
            claimedBitMap[claimedWordIndex] |
            (1 << claimedBitIndex);
    }

    // TODO: Tokens not claimed within a given timeout should go to a treasury
    // wallet address set on that contract.
}
