#!/usr/local/bin/node

// Code to generate a checksumed address for a contract - required for using
// a contant address for a contract in solidity.  Runs in node.js.

//
// to test
//
// $ ./cs.js test
//
// to run with data
//
// $ ./cs.js 0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359
//

//
// Cut/pasete from some webpage on EIP-55 encoding.
//
// This works to encode a number into an Address
// This implements EIP-55 checksum encoding of numbers so that they will
// be recoginized as addresses in Ethereum/Solidity.
//
// Example:
//  convert:
//      GenRequestIDSequence = GenRequestID(0x9fbda871d559710256a2502a2517b794b482db40);
//  to:
//		GenRequestIDSequence = GenRequestID(0x9FBDa871d559710256a2502A2517b794B482Db40);
// Notice the change in "case" of some of the a..f hex valeus.
//
// TODO: Should add test that checks that returned value "isAddress" - find code.
// TODO: add into build process for .m4 file so addresses get encoded properly.
// TODO: add isAddress into build process to check them.
//

//
// Required install 
// 	$ npm install keccak
//

const createKeccakHash = require('keccak')

function toChecksumAddress (address) {
	address = address.toLowerCase().replace('0x', '')
	var hash = createKeccakHash('keccak256').update(address).digest('hex')
	var ret = '0x'

	for (var i = 0; i < address.length; i++) {
		if (parseInt(hash[i], 16) >= 8) {
			ret += address[i].toUpperCase()
		} else {
			ret += address[i]
		}
	}

	return ret
}

// -------------------------------------------------------------------- test --------------------------------------------------------------------------------

var argv = process.argv;
// console.log ( argv );

if ( argv.length > 2 && argv[2] === "test" ) {

    if ( toChecksumAddress('0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359') === '0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359' ) {
        console.log ( "PASS\n" );
    } else {
        console.log ( "FAIL\n" );
    }

} else if ( argv.length > 2 ) {
    var output = toChecksumAddress(argv[2]);
    console.log ( output );
}

