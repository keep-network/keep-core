import {bls} from './helpers/data'

const BLS = artifacts.require('./cryptography/BLS.sol');
const AltBn128 = artifacts.require('./cryptography/AltBn128.sol');
const AltBn128Stub = artifacts.require('./stubs/AltBn128Stub.sol');

contract('TestBLSRoundtrip', function(accounts) {
  let blsLibrary, altBn128Library, altBn128StubLibrary, from;
  
  before(async () => {
    from = accounts[1]
    blsLibrary = await BLS.new();
    altBn128Library = await AltBn128.new();
    
    await AltBn128Stub.link("AltBn128", altBn128Library.address);
    
    altBn128StubLibrary = await AltBn128Stub.new();
  });

  it("should be able to sign a message and verify it", async function() {
    let g1point = await altBn128StubLibrary.sign(bls.secretKey, {from: from});
    let g1pointX = web3.utils.padLeft(web3.utils.toHex(g1point[0].toString()), 64)
    let g1pointY = web3.utils.padLeft(web3.utils.toHex(g1point[1].toString()), 64)
    let signature = '0x' + Buffer.concat([
      Buffer.from(g1pointX.substr(2), 'hex'),
      Buffer.from(g1pointY.substr(2), 'hex')
    ]).toString('hex');

    let message = await altBn128StubLibrary.g1HashToPoint(from);

    let actual = await blsLibrary.verify(bls.groupPubKey, message, signature)
    assert.isTrue(actual, "Should be able to verify valid BLS signature.");
  })

});
