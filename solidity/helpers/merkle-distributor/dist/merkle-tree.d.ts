/// <reference types="node" />
export default class MerkleTree {
    private readonly elements;
    private readonly bufferElementPositionIndex;
    private readonly layers;
    constructor(elements: Buffer[]);
    getLayers(elements: Buffer[]): Buffer[][];
    getNextLayer(elements: Buffer[]): Buffer[];
    static combinedHash(first: Buffer, second: Buffer): Buffer;
    getRoot(): Buffer;
    getHexRoot(): string;
    getProof(el: Buffer): Buffer[];
    getHexProof(el: Buffer): string[];
    private static getPairElement;
    private static bufDedup;
    private static bufArrToHexArr;
    private static sortAndConcat;
}
