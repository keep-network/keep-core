export const bls = {
  // data generated using master secret key 123
  
  // compressed altbn128 public key for secret key 123
  groupPubKey: "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd",
  
  // initial beacon entry from KeepRandomBeaconServiceImplV1.sol
  previousEntry: web3.utils.toBN('67739255176204957841308900500337009958963301977368411094653342041505945563546'),

  // compressed group signature over previousEntry
  groupSignature: web3.utils.toBN('61135065625162045940297002520113043735613273581525878142189953979726299940463'),
  
  // compressed group signature over groupSignature
  nextGroupSignature: web3.utils.toBN('12042949674602040267884054775326766869329495529174676183189893197609078382612'),
  
  // compressed group signature over nextGroupSignature 
  nextNextGroupSignature: web3.utils.toBN('6429987415739497838560918343065324227487415866120830222155997935020021330721'),
};
