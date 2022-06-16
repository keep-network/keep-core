import { ethers } from "ethers";

const password = process.env.KEEP_ETHEREUM_PASSWORD || "password"

async function unlockAccounts() {
    const provider = new ethers.providers.JsonRpcProvider(process.env.NETWORK)
    const accounts = await provider.listAccounts()

    console.log(`Total accounts: ${accounts.length}`)
    console.log(`---------------------------------`)

    for (let i = 0; i < accounts.length; i++) {
      const account = accounts[i]
      
      try {
        console.log(`\nUnlocking account: ${account}`)
        const signerAccount = provider.getSigner(account)
        await signerAccount.unlock(password)
        console.log(`Account unlocked!`)
      } catch (error) {
        console.log(`\nAccount: ${account} not unlocked!`)
        console.error(error)
      }
      console.log(`\n---------------------------------`)
    }

}

unlockAccounts()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });