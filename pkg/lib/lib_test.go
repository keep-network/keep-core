package lib

import "github.com/ethereum/go-ethereum/crypto/sha3"

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

/* Similar code in Ecma-Script

This program computes your address given a public key. In your case it returns "0xd09d3103ccabfb769edc3e9b01500ca7241d470a" as address

const assert = require('assert');
const EC = require('elliptic').ec;
const keccak256 = require('js-sha3').keccak256;

async function main() {
  try {
    const ec = new EC('secp256k1');

    // Decode public key
    const key = ec.keyFromPublic('025f37d20e5b18909361e0ead7ed17c69b417bee70746c9e9c2bcb1394d921d4ae', 'hex');

    // Convert to uncompressed format
    const publicKey = key.getPublic().encode('hex').slice(2);

    // Now apply keccak
    const address = keccak256(Buffer.from(publicKey, 'hex')).slice(64 - 40);

    console.log(`Public Key: 0x${publicKey}`);
    console.log(`Address: 0x${address.toString()}`);
  } catch (err) {
    console.log(err);
  }
}

main();

-----------------------------------------------

https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go

func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := FromECDSAPub(&p)
	return common.BytesToAddress(Keccak256(pubBytes[1:])[12:])
}


r and s values:
https://bitcoin.stackexchange.com/questions/58853/how-do-you-figure-out-the-r-and-s-out-of-a-signature-using-python/58862
https://bitcoin.stackexchange.com/questions/2376/ecdsa-r-s-encoding-as-a-signature

*/
