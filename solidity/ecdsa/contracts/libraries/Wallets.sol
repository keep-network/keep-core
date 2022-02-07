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

    /// @notice Checks if a wallet with given ID was registered.
    /// @param walletID Wallet's ID.
    /// @return True if wallet was registered, false otherwise.
    function isWalletRegistered(Data storage self, bytes32 walletID)
        external
        view
        returns (bool)
    {
        return self.registry[walletID].membersIdsHash != bytes32(0);
    }
}
