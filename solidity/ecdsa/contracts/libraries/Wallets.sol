// SPDX-License-Identifier: MIT
//
// ▓▓▌ ▓▓ ▐▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▄
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▌▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓    ▓▓▓▓▓▓▓▀    ▐▓▓▓▓▓▓    ▐▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▄▄▓▓▓▓▓▓▓▀      ▐▓▓▓▓▓▓▄▄▄▄         ▓▓▓▓▓▓▄▄▄▄         ▐▓▓▓▓▓▌   ▐▓▓▓▓▓▓
//   ▓▓▓▓▓▓▓▓▓▓▓▓▓▀        ▐▓▓▓▓▓▓▓▓▓▓         ▓▓▓▓▓▓▓▓▓▓         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
//   ▓▓▓▓▓▓▀▀▓▓▓▓▓▓▄       ▐▓▓▓▓▓▓▀▀▀▀         ▓▓▓▓▓▓▀▀▀▀         ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▀
//   ▓▓▓▓▓▓   ▀▓▓▓▓▓▓▄     ▐▓▓▓▓▓▓     ▓▓▓▓▓   ▓▓▓▓▓▓     ▓▓▓▓▓   ▐▓▓▓▓▓▌
// ▓▓▓▓▓▓▓▓▓▓ █▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
// ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓ ▐▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓▓▓
//
//                           Trust math, not hardware.

pragma solidity ^0.8.9;

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// TODO: This contract is just a Stub implementation that was used for gas
// comparisons for Wallets creation. It should be implemented according to the
// Wallets' actual use case.
library Wallets {
    using ECDSA for bytes32;

    struct Wallet {
        // TODO: Verify if we want to store the whole public key or having
        // just the public key hash is enough.
        // Keccak256 hash of group members identifiers array. Group members do not
        // include operators selected by the sortition pool that misbehaved during DKG.
        bytes32 membersIdsHash;
        bytes32 digestToSign;
    }

    struct Signature {
        uint8 v;
        bytes32 r;
        bytes32 s;
    }

    struct Data {
        // Mapping of keccak256 hashes of wallet public keys to wallet details.
        // Hash of public key is considered an unique wallet identifier (walletID).
        mapping(bytes32 => Wallet) registry;
    }

    event SignatureRequested(bytes32 indexed walletID, bytes32 indexed digest);

    event SignatureSubmitted(
        bytes32 indexed walletID,
        bytes32 indexed digest,
        Signature signature
    );

    /// @notice Registers a new wallet.
    /// @param _membersIdsHash Keccak256 hash of group members identifiers array.
    /// @param _publicKeyHash Keccak256 hash of group public key.
    /// @return walletID ID of the newly registered wallet.
    function addWallet(
        Data storage self,
        bytes32 _membersIdsHash,
        bytes32 _publicKeyHash
    ) internal returns (bytes32 walletID) {
        walletID = _publicKeyHash;

        // TODO: If we decide to store the group public key we should switch this
        // check to use the public key.
        require(
            self.registry[walletID].membersIdsHash == bytes32(0),
            "Wallet with given public key hash already exists"
        );

        self.registry[walletID].membersIdsHash = _membersIdsHash;
    }

    /// @notice Requests a new signature.
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param digest Digest to sign.
    function requestSignature(
        Data storage self,
        bytes32 walletID,
        bytes32 digest
    ) internal {
        // TODO: When we decide to introduce wallet termination it would be enough
        // to check if wallet is active instead checking if value is set.
        require(
            self.registry[walletID].membersIdsHash != bytes32(0),
            "Wallet with given public key hash doesn't exist"
        );

        // TODO: Implement; Compare with AbstractBondedECDSAKeep from V1.

        self.registry[walletID].digestToSign = digest;

        emit SignatureRequested(walletID, digest);
    }

    // TODO: Compare gas usage if we store the whole public key not just the hash
    /// @notice Submits a calculated signature for the digest that is currently
    ///         under signing.
    /// @dev Implementation assumes the walletID is a public key hash of the wallet.
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param signature Calculated signature.
    function submitSignature(
        Data storage self,
        bytes32 walletID,
        Signature calldata signature
    ) external {
        // TODO: Check if wallet is available
        Wallet storage wallet = self.registry[walletID];

        require(
            wallet.digestToSign != bytes32(0),
            "Signature was not requested"
        );

        // TODO: Implement; Compare with AbstractBondedECDSAKeep from V1.

        // Calculate address from the walletID as it is the same as the wallet's
        // public key hash.
        require(
            address(uint160(uint256(walletID))) ==
                wallet.digestToSign.recover(
                    signature.v,
                    signature.r,
                    signature.s
                ),
            "Invalid signature"
        );

        emit SignatureSubmitted(walletID, wallet.digestToSign, signature);

        delete wallet.digestToSign;
    }
}
