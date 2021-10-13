const { MerkleTree } = require('merkletreejs');
const keccak256 = require('keccak256');
const { toBN } = require('../test/helpers/utils');
const Wallet = require('ethereumjs-wallet').default;
const { promisify } = require('util');
const randomBytesAsync = promisify(require('crypto').randomBytes)
const { BN, ether } = require('@openzeppelin/test-helpers');
var qr = require('qr-image');
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

function verifyProof (wallet, amount, proof, root) {
    const tree = new MerkleTree([], keccak128, { sortPairs: true });
    const element = wallet + toBN(amount).toString(16, 64);
    const node = MerkleTree.bufferToHex(keccak128(element));
    console.log(tree.verify(proof, node, root));
}

async function gen_priv() {
    return (await randomBytesAsync(16)).toString('hex').padStart(64, '0');
}

async function gen_privs(n) {
    return Promise.all(Array.from({length: n}, gen_priv));
}

function uri_encode(b) {
    return encodeURIComponent(b.toString('base64').replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '!'));
}

function save_qr(i, url) {
    const code = qr.imageSync(url, { type: 'png' });
    fs.writeFileSync(`src/qr/${i}.png`, code);
}

function uri_decode(s, root) {
    const b = Buffer.from(s.substring(PREFIX.length + 2).replace(/-/g, '+').replace(/_/g, '/').replace(/!/g, '='), 'base64');
    // const v_buf = b.slice(0, 1);
    const k_buf = b.slice(1, 17);
    const a_buf = b.slice(17, 29);
    var p_buf = b.slice(29);

    const proof = [];
    while (p_buf.length > 0) {
        proof.push(p_buf.slice(0, 16));
        p_buf = p_buf.slice(16);
    }

    const key = k_buf.toString('hex').padStart(64, '0');
    const wallet = Wallet.fromPrivateKey(Buffer.from(key, 'hex')).getAddressString()
    const amount = (new BN(a_buf.toString('hex'), 16)).toString();

    verifyProof(wallet, amount, proof, root);
}

function gen_url(priv, amount, proof) {
    const v_buf = Buffer.from([0]);
    const k_buf = Buffer.from(priv.substring(32), 'hex');
    const a_buf = Buffer.from(toBN(amount).toString(16, 24), 'hex');
    const p_buf = Buffer.concat(proof.map(p => p.data));

    const base_args = uri_encode(Buffer.concat([v_buf, k_buf, a_buf, p_buf]));
    return PREFIX + 'd=' + base_args;
}

async function main () {
    const privs = await gen_privs(COUNT);
    // console.log(privs);

    const accounts = privs.map(p => Wallet.fromPrivateKey(Buffer.from(p, 'hex')).getAddressString());
    // console.log(accounts);

    const drop = makeDrop(accounts, AMOUNT);

    // console.log(drop);
    // console.log(drop.proofs[0]);

    for (let i = 0; i < COUNT; i++) {
        const url = gen_url(privs[i], AMOUNT, drop.proofs[i]);
        console.log(url);
        save_qr(i, url);
    }

    // uri_decode(url, drop.root);
}

main();
