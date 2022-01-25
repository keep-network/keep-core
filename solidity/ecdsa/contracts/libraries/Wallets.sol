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

library Wallets {
    using ECDSA for bytes32;

    struct Wallet {
        bytes32 publicKeyHash;
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
        mapping(uint256 => Wallet) registry;
        // Holds keccak256 hashes of group public keys in the order of registration.
        uint256 walletCounter;
    }

    event SignatureRequested(uint256 indexed walletID, bytes32 indexed digest);

    event SignatureSubmitted(
        uint256 indexed walletID,
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
    ) internal returns (uint256 walletID) {
        walletID = ++self.walletCounter;

        self.registry[walletID].publicKeyHash = _publicKeyHash;
        self.registry[walletID].membersIdsHash = _membersIdsHash;
    }

    /// @notice Requests a new signature.
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param digest Digest to sign.
    function requestSignature(
        Data storage self,
        uint256 walletID,
        bytes32 digest
    ) internal {
        // TODO: When we decide to introduce wallet termination it would be enough
        // to check if wallet is active instead checking if public key is set.
        require(
            self.registry[walletID].publicKeyHash != bytes32(0),
            "Wallet with given ID doesn't exist"
        );

        // TODO: Implement; Compare with AbstractBondedECDSAKeep from V1.

        self.registry[walletID].digestToSign = digest;

        emit SignatureRequested(walletID, digest);
    }

    // TODO: Compare gas usage if we store the whole public key not just the hash
    /// @notice Submits a calculated signature for the digest that is currently
    ///         under signing.
    /// @param walletID ID of a wallet that should calculate a signature.
    /// @param publicKey Group public key.
    /// @param signature Calculated signature.
    function submitSignature(
        Data storage self,
        uint256 walletID,
        bytes calldata publicKey,
        Signature calldata signature
    ) external {
        // TODO: Check if wallet is available
        Wallet storage wallet = self.registry[walletID];
        require(
            wallet.publicKeyHash == keccak256(publicKey),
            "Invalid public key"
        );

        // TODO: Implement; Compare with AbstractBondedECDSAKeep from V1.

        require(
            address(uint160(uint256(wallet.publicKeyHash))) ==
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
