const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');
const { toBN } = require('../test/helpers/utils');
const Wallet = require('ethereumjs-wallet').default;
const { promisify } = require('util');
const randomBytesAsync = promisify(require('crypto').randomBytes);
const { ether } = require('@openzeppelin/test-helpers');
const qr = require('qr-image');
const fs = require('fs');

function keccak128 (input) {
    return keccak256(input).slice(0, 16);
}

const AMOUNT = ether('50');
const COUNT = 1024;
const PREFIX = 'https://app.1inch.io/#/1/qr?';

function makeDrop (wallets, amount) {
    const elements = wallets.map((w, i) => w + toBN(amount).toString(16, 64));
    const leaves = elements.map(keccak128).map(x => MerkleTree.bufferToHex(x));
    const tree = new MerkleTree(leaves, keccak128, { sortPairs: true });
    const root = tree.getHexRoot();
    const proofs = leaves.map(tree.getProof, tree);
    return { elements, leaves, root, proofs };
}

async function genPriv () {
    return (await randomBytesAsync(16)).toString('hex').padStart(64, '0');
}

async function genPrivs (n) {
    return Promise.all(Array.from({ length: n }, genPriv));
}

function uriEncode (b) {
    return encodeURIComponent(b.toString('base64').replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '!'));
}

function saveQr (i, url) {
    const code = qr.imageSync(url, { type: 'png' });
    fs.writeFileSync(`src/qr/${i}.png`, code);
}

// function verifyProof (wallet, amount, proof, root) {
//     const tree = new MerkleTree([], keccak128, { sortPairs: true });
//     const element = wallet + toBN(amount).toString(16, 64);
//     const node = MerkleTree.bufferToHex(keccak128(element));
//     console.log(tree.verify(proof, node, root));
// }

// function uriDecode (s, root) {
//     const b = Buffer.from(s.substring(PREFIX.length + 2).replace(/-/g, '+').replace(/_/g, '/').replace(/!/g, '='), 'base64');
//     // const vBuf = b.slice(0, 1);
//     const kBuf = b.slice(1, 17);
//     const aBuf = b.slice(17, 29);
//     let pBuf = b.slice(29);
//
//     const proof = [];
//     while (pBuf.length > 0) {
//         proof.push(pBuf.slice(0, 16));
//         pBuf = pBuf.slice(16);
//     }
//
//     const key = kBuf.toString('hex').padStart(64, '0');
//     const wallet = Wallet.fromPrivateKey(Buffer.from(key, 'hex')).getAddressString();
//     const amount = (new BN(aBuf.toString('hex'), 16)).toString();
//
//     verifyProof(wallet, amount, proof, root);
// }

function genUrl (priv, amount, proof) {
    const vBuf = Buffer.from([0]);
    const kBuf = Buffer.from(priv.substring(32), 'hex');
    const aBuf = Buffer.from(toBN(amount).toString(16, 24), 'hex');
    const pBuf = Buffer.concat(proof.map(p => p.data));

    const baseArgs = uriEncode(Buffer.concat([vBuf, kBuf, aBuf, pBuf]));
    return PREFIX + 'd=' + baseArgs;
}

async function main () {
    const privs = await genPrivs(COUNT);
    // console.log(privs);

    const accounts = privs.map(p => Wallet.fromPrivateKey(Buffer.from(p, 'hex')).getAddressString());
    // console.log(accounts);

    const drop = makeDrop(accounts, AMOUNT);

    // console.log(drop);
    // console.log(drop.proofs[0]);

    for (let i = 0; i < COUNT; i++) {
        const url = genUrl(privs[i], AMOUNT, drop.proofs[i]);
        console.log(url);
        saveQr(i, url);
    }

    // uriDecode(url, drop.root);
}

main();
