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

import "@openzeppelin/contracts/access/Ownable.sol";

library Wallets {
    struct Wallet {
        // Keccak256 hash of group members identifiers array. Group members do not
        // include operators selected by the sortition pool that misbehaved during DKG.
        bytes32 membersIdsHash;
        // Uncompressed ECDSA public key stored as X and Y coordinates (32 bytes each).
        bytes32 publicKeyX;
        bytes32 publicKeyY;
    }

    struct Data {
        // Mapping of keccak256 hashes of wallet public keys to wallet details.
        // Hash of public key is considered an unique wallet identifier.
        mapping(bytes32 => Wallet) registry;
    }

    /// @notice Registers a new wallet.
    /// @dev Uses a public key hash as a unique identifier of a wallet.
    /// @param membersIdsHash Keccak256 hash of group members identifiers array
    /// @param publicKey Uncompressed public key
    /// @return walletID Wallet's ID
    /// @return publicKeyX Wallet's public key's X coordinate
    /// @return publicKeyY Wallet's public key's Y coordinate
    function addWallet(
        Data storage self,
        bytes32 membersIdsHash,
        bytes calldata publicKey
    )
        internal
        returns (
            bytes32 walletID,
            bytes32 publicKeyX,
            bytes32 publicKeyY
        )
    {
        walletID = keccak256(publicKey);

        require(
            self.registry[walletID].publicKeyX == bytes32(0),
            "Wallet with the given public key already exists"
        );
        require(publicKey.length == 64, "Invalid length of the public key");

        publicKeyX = bytes32(publicKey[:32]);
        publicKeyY = bytes32(publicKey[32:]);

        self.registry[walletID].membersIdsHash = membersIdsHash;
        self.registry[walletID].publicKeyX = publicKeyX;
        self.registry[walletID].publicKeyY = publicKeyY;
    }

    /// @notice Deletes wallet with the given ID from the registry. Reverts
    ///         if wallet with the given ID has not been registered or if it
    ///         has already been closed.
    function deleteWallet(Data storage self, bytes32 walletID) internal {
        require(
            isWalletRegistered(self, walletID),
            "Wallet with the given ID has not been registered"
        );

        delete self.registry[walletID];
    }

    /// @notice Checks if a wallet with the given ID is registered.
    /// @param walletID Wallet's ID
    /// @return True if a wallet is registered, false otherwise
    function isWalletRegistered(Data storage self, bytes32 walletID)
        internal
        view
        returns (bool)
    {
        return self.registry[walletID].publicKeyX != bytes32(0);
    }

    /// @notice Returns Keccak256 hash of the wallet signing group members
    ///         identifiers array. Group members do not include operators
    ///         selected by the sortition pool that misbehaved during DKG.
    ///         Reverts if wallet with the given ID is not registered.
    /// @param walletID ID of the wallet
    /// @return Wallet signing group members hash
    function getWalletMembersIdsHash(Data storage self, bytes32 walletID)
        internal
        view
        returns (bytes32)
    {
        require(
            isWalletRegistered(self, walletID),
            "Wallet with the given ID has not been registered"
        );

        return self.registry[walletID].membersIdsHash;
    }

    /// @notice Gets public key of a wallet with the given wallet ID.
    ///         The public key is returned as X and Y coordinates.
    ///         Reverts if wallet with the given ID is not registered.
    /// @param walletID ID of the wallet
    /// @return x Public key X coordinate
    /// @return y Public key Y coordinate
    function getWalletPublicKeyCoordinates(Data storage self, bytes32 walletID)
        internal
        view
        returns (bytes32 x, bytes32 y)
    {
        require(
            isWalletRegistered(self, walletID),
            "Wallet with the given ID has not been registered"
        );

        Wallet storage wallet = self.registry[walletID];

        return (wallet.publicKeyX, wallet.publicKeyY);
    }

    /// @notice Gets public key of a wallet with the given wallet ID.
    ///         The public key is returned in an uncompressed format as a 64-byte
    ///         concatenation of X and Y coordinates.
    ///         Reverts if wallet with the given ID is not registered.
    /// @param walletID ID of the wallet
    /// @return Uncompressed public key of the wallet
    function getWalletPublicKey(Data storage self, bytes32 walletID)
        internal
        view
        returns (bytes memory)
    {
        (bytes32 x, bytes32 y) = getWalletPublicKeyCoordinates(self, walletID);
        return bytes.concat(x, y);
    }
}
