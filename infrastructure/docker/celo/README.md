# Celo client Docker image for local development

Building the image:
```
$ docker build -t celo-dev .
```

Starting Celo blockchain locally with the image built:
```
$ docker run -it -v $(pwd)/data:/data -p 8545:8545 -p 8546:8546 celo-dev
```

`$(pwd)/data` directory is where the chain data will be stored. Since it is 
a Docker volume, chain data survive container restarts.

In case the developer wants to start from a fresh chain, it is enough to remove
chain data: `$ rm -rf data/celo`, or even remove the entire data directory: 
`$ rm -rf data` and start the container again. Keystore will be automatically 
populated and chain genesis executed.