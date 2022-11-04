package electrs

// https://blockstream.info/testnet/api/tx/c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479
const Tx = `
{
    "txid": "c580e0e352570d90e303d912a506055ceeb0ee06f97dce6988c69941374f5479",
    "version": 1,
    "locktime": 0,
    "vin":
    [
        {
            "txid": "e788a344a86f7e369511fe37ebd1d74686dde694ee99d06db5db3d4a14719b1d",
            "vout": 1,
            "prevout":
            {
                "scriptpubkey": "76a914e257eccafbc07c381642ce6e7e55120fb077fbed88ac",
                "scriptpubkey_asm": "OP_DUP OP_HASH160 OP_PUSHBYTES_20 e257eccafbc07c381642ce6e7e55120fb077fbed OP_EQUALVERIFY OP_CHECKSIG",
                "scriptpubkey_type": "p2pkh",
                "scriptpubkey_address": "n29kLYP73xJkHzx26VthkCAEPeeCBRqD81",
                "value": 1382770
            },
            "scriptsig": "47304402206f8553c07bcdc0c3b906311888103d623ca9096ca0b28b7d04650a029a01fcf9022064cda02e39e65ace712029845cfcf58d1b59617d753c3fd3556f3551b609bbb00121039d61d62dcd048d3f8550d22eb90b4af908db60231d117aeede04e7bc11907bfa",
            "scriptsig_asm": "OP_PUSHBYTES_71 304402206f8553c07bcdc0c3b906311888103d623ca9096ca0b28b7d04650a029a01fcf9022064cda02e39e65ace712029845cfcf58d1b59617d753c3fd3556f3551b609bbb001 OP_PUSHBYTES_33 039d61d62dcd048d3f8550d22eb90b4af908db60231d117aeede04e7bc11907bfa",
            "is_coinbase": false,
            "sequence": 4294967295
        }
    ],
    "vout":
    [
        {
            "scriptpubkey": "a9143ec459d0f3c29286ae5df5fcc421e2786024277e87",
            "scriptpubkey_asm": "OP_HASH160 OP_PUSHBYTES_20 3ec459d0f3c29286ae5df5fcc421e2786024277e OP_EQUAL",
            "scriptpubkey_type": "p2sh",
            "scriptpubkey_address": "2Mxy76sc1qAxiJ1fXMXDXqHvVcPLh6Lf12C",
            "value": 20000
        },
        {
            "scriptpubkey": "0014e257eccafbc07c381642ce6e7e55120fb077fbed",
            "scriptpubkey_asm": "OP_0 OP_PUSHBYTES_20 e257eccafbc07c381642ce6e7e55120fb077fbed",
            "scriptpubkey_type": "v0_p2wpkh",
            "scriptpubkey_address": "tb1quft7ejhmcp7rs9jzeeh8u4gjp7c807ld7vk4ss",
            "value": 1360550
        }
    ],
    "size": 220,
    "weight": 880,
    "fee": 2220,
    "status":
    {
        "confirmed": true,
        "block_height": 2135049,
        "block_hash": "0000000064d5a6c6180baee14ac75161e9f1c626b9eeda948acdb48595e49f8a",
        "block_time": 1641682320
    }
}
`

// https://blockstream.info/testnet/api/block/000000000000002af10911b8db32ed34dc6ea6515f84af5f7b82973c9a839e6d
const Block = `
{
    "id": "000000000000002af10911b8db32ed34dc6ea6515f84af5f7b82973c9a839e6d",
    "height": 2135502,
    "version": 536870916,
    "timestamp": 1641914003,
    "tx_count": 82,
    "size": 22494,
    "weight": 62091,
    "merkle_root": "1251774996b446f85462d5433f7a3e384ac1569072e617ab31e86da31c247de2",
    "previousblockhash": "000000000066450030efdf72f233ed2495547a32295deea1e2f3a16b1e50a3a5",
    "mediantime": 1641908180,
    "nonce": 778087099,
    "bits": 436256810,
    "difficulty": 22350181
}
`
