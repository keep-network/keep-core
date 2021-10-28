const blsData = {
  // data generated using master secret key 123
  secretKey: 123,

  // altbn128 public key for secret key 123
  groupPubKey:
    "0x1f1954b33144db2b5c90da089e8bde287ec7089d5d6433f3b6becaefdb678b1b2a9de38d14bef2cf9afc3c698a4211fa7ada7b4f036a2dfef0dc122b423259d01659dc18b57722ecf6a4beb4d04dfe780a660c4c3bb2b165ab8486114c464c621bf37ecdba226629c20908c7f475c5b3a7628ce26d696436eab0b0148034dfcd",

  // initial beacon entry
  previousEntry:
    "0x15c30f4b6cf6dbbcbdcc10fe22f54c8170aea44e198139b776d512d8f027319a1b9e8bfaf1383978231ce98e42bafc8129f473fc993cf60ce327f7d223460663",

  // group signature over previousEntry
  groupSignature:
    "0x112d462728e89432b0fe40251eeb6608aed4560f3dc833a9877f5010ace9b1312006dbbe2f30c6e0e3e7ec47dc078b7b6b773379d44d64e44ec4e017bfa7375c",

  // group signature over previousEntry uint256
  groupSignatureUint256:
    "33136845259729814081977577432759716433762925710284529101628804651946005705295",

  // group signature over groupSignature
  nextGroupSignature:
    "0x144b0508140c2c63fac298ee0cdd98571971a2d958f7c97d2bab82f3e1e727542d08314d6f087aca6ec2173b9a1d928cb80ff45258984a8929977a58d8b2fc26",
}

export default blsData
