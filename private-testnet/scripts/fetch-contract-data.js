const {Storage} = require('@google-cloud/storage');

async function fetchContractData(bucketName, sourceFile, destinationFile) {

  const storage = new Storage();

   await storage
    .bucket(bucketName)
    .file(sourceFile)
    .download({destination: destinationFile});

    console.log(`gs://${bucketName}/${sourceFile} downloaded to ${destinationFile}.`);
}

fetchContractData("keep-test-contract-data", "keep-core/KeepToken.json", "../KeepToken.json")
fetchContractData("keep-test-contract-data", "keep-core/TokenStaking.json", "../TokenStaking.json")
