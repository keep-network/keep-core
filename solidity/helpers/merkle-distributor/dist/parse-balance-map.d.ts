interface MerkleDistributorInfo {
    merkleRoot: string;
    tokenTotal: string;
    claims: {
        [account: string]: {
            index: number;
            amount: string;
            proof: string[];
            flags?: {
                [flag: string]: boolean;
            };
        };
    };
}
declare type OldFormat = {
    [account: string]: number | string;
};
declare type NewFormat = {
    address: string;
    earnings: string;
    reasons: string;
};
export declare function parseBalanceMap(balances: OldFormat | NewFormat[]): MerkleDistributorInfo;
export {};
