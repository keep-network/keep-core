- Eth account with contracts: `0xecfbf98c33478afb97a3a08704f1fa94299986c8`

- Unlock

> personal.unlockAccount("0xecfbf98c33478afb97a3a08704f1fa94299986c8", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
true

- Truffle config:

```
require('babel-register');
require('babel-polyfill');

module.exports = {
  networks: {
    keep_dev: {
      host: "eth-miner-node.default.svc.cluster.local",
      port: 8545,
      network_id: "*",
      gas: 4712388,
      from: "0xecfbf98c33478afb97a3a08704f1fa94299986c8"
    }
  }
};

```

- keep-client accounts:

Pre-generated in genesis file.

```
0x7faf88360448efcfd611933d8893d6fde30cbb1a - boostrap-peer-0
0x9b62bdbc11a9f4d131f32dafb28f707750d2f248 - peer-0
0xa51a6685d7ce0bc25582f3573938d2b1a0daee44 - peer-1
0xc47d085e4e555584d8db83ed7abec7bdc5b9d24d - peer-2
0x9892a772c34b89eb58abc00f2c67f4752bcca735 - peer-3
```

```
personal.unlockAccount("0x7faf88360448efcfd611933d8893d6fde30cbb1a", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
personal.unlockAccount("0x9b62bdbc11a9f4d131f32dafb28f707750d2f248", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
personal.unlockAccount("0xa51a6685d7ce0bc25582f3573938d2b1a0daee44", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
personal.unlockAccount("0xc47d085e4e555584d8db83ed7abec7bdc5b9d24d", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
personal.unlockAccount("0x9892a772c34b89eb58abc00f2c67f4752bcca735", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
```

Result:
```
> personal.unlockAccount("0x7faf88360448efcfd611933d8893d6fde30cbb1a", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
true
> personal.unlockAccount("0x9b62bdbc11a9f4d131f32dafb28f707750d2f248", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
true
> personal.unlockAccount("0xa51a6685d7ce0bc25582f3573938d2b1a0daee44", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
true
> personal.unlockAccount("0xc47d085e4e555584d8db83ed7abec7bdc5b9d24d", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
true
> personal.unlockAccount("0x9892a772c34b89eb58abc00f2c67f4752bcca735", "doughnut_armenian_parallel_firework_backbite_employer_singlet", 150000);
true
```

```

- stake keep-client accounts

Using `demo.js` script

```
sthompson22@99-0-0-6:~/projects/keep-core/contracts/solidity/scripts(sthompson22/playing/local-keep-client⚡) » truffle exec ./demo.js --network keep_dev
Using network 'keep_dev'.
successfully staked KEEP tokens for account 0xecfbf98c33478afb97a3a08704f1fa94299986c8
successfully staked KEEP tokens for account 0x7faf88360448efcfd611933d8893d6fde30cbb1a
successfully staked KEEP tokens for account 0x9b62bdbc11a9f4d131f32dafb28f707750d2f248
successfully staked KEEP tokens for account 0xa51a6685d7ce0bc25582f3573938d2b1a0daee44
successfully staked KEEP tokens for account 0xc47d085e4e555584d8db83ed7abec7bdc5b9d24d
successfully staked KEEP tokens for account 0x9892a772c34b89eb58abc00f2c67f4752bcca735
```

- fetch eth account keyfiles

```
kubectl cp  default/eth-miner-node-0:/mnt/eth-data/keystore ./
```