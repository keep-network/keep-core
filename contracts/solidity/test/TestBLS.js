const BLS = artifacts.require('./cryptography/BLS.sol');

contract('TestBLS', function() {

  let bls;
  before(async () => {
    bls = await BLS.new();
  });

  it("should be able to verify threshold BLS recovered/reconstructed signature", async function() {

    // Data generated using client Go code with master secret key 123 and message "Hello!"
    let result = await bls.verify(
      "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd",
      "0x48656c6c6f21",
      "0x884b130ed81751b63d0f5882483d4a24a7640bdf371f23b78dbeb520c84e3a85"
    )
    assert.isTrue(result, "Should be able to verify valid BLS signature.");
  });
  
  it("should use reasonable amount of gas", async () => {

    // Data generated using client Go code with master secret key 123 and message "Hello!"
    let gasEstimate = await bls.verify.estimateGas(
      "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd",
      "0x48656c6c6f21",
      "0x884b130ed81751b63d0f5882483d4a24a7640bdf371f23b78dbeb520c84e3a85"
    )

    // make sure no change will make the verification more expensive than it's now
    assert.isBelow(gasEstimate, 378257, "BLS verification is too expensive")
  })

  it("should be able to verify BLS aggregated signature", async function() {

    // Data generated using client Go code with multiple random signers signing the same message "Hello!"
    let result = await bls.verify(
      "0x2460893c494f57366f41e0e4062437fa34ccc2b356943244893241cd7d2aacce2142936a42f66c68261a663bef68d398e9e32e4c73009b482567943e1ffa51d42672228197922e5ba100cb8f9a83efa1b0c67f80b5960a2185952a0a43392b821144bf3de620433cf9bb3dbf2a2d08d1f7b9eacfd9dfecdd581f6a006b7ed7f8",
      "0x48656c6c6f21",
      "0x880e16fb15a6c6757c1dd22139504ccefe00fb6c3928f7ea026871fa80a68e46"
    )
    assert.isTrue(result, "Should be able to verify valid BLS signature.");
  });

  it("should fail to verify non valid BLS signature", async function() {

    let result = await bls.verify(
      "0x2460893c494f57366f41e0e4062437fa34ccc2b356943244893241cd7d2aacce2142936a42f66c68261a663bef68d398e9e32e4c73009b482567943e1ffa51d42672228197922e5ba100cb8f9a83efa1b0c67f80b5960a2185952a0a43392b821144bf3de620433cf9bb3dbf2a2d08d1f7b9eacfd9dfecdd581f6a006b7ed7f8",
      "0x48656c6c6f21",
      "0x884b130ed81751b63d0f5882483d4a24a7640bdf371f23b78dbeb520c84e3a85"      
    )
    assert.isFalse(result, "Should return false for failed verification.");
  });

  it("should fail to verify BLS signature without valid message", async function() {

    let result = await bls.verify(
      "0x2460893c494f57366f41e0e4062437fa34ccc2b356943244893241cd7d2aacce2142936a42f66c68261a663bef68d398e9e32e4c73009b482567943e1ffa51d42672228197922e5ba100cb8f9a83efa1b0c67f80b5960a2185952a0a43392b821144bf3de620433cf9bb3dbf2a2d08d1f7b9eacfd9dfecdd581f6a006b7ed7f8",
      "0x123456789",
      "0x880e16fb15a6c6757c1dd22139504ccefe00fb6c3928f7ea026871fa80a68e46"
    )
    assert.isFalse(result, "Should return false for failed verification.");
  });

  it("should fail to verify BLS signature without valid public key", async function() {

    let result = await bls.verify(
      "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd",
      "0x48656c6c6f21",
      "0x880e16fb15a6c6757c1dd22139504ccefe00fb6c3928f7ea026871fa80a68e46"
    )
    assert.isFalse(result, "Should return false for failed verification.");
  });
});
