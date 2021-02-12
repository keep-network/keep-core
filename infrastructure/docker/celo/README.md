# Celo client Docker image for local development

Building the image:
```
$ docker build -t celo-dev .
```

Starting Celo blockchain locally exposing ports 8545 and 8546:
```
$ docker run -it -v $(pwd)/data:/mnt/data -p 8545:8545 -p 8546:8546 celo-dev
```

Starting Celo blockchain locally with RPC and WS ports randomly assigned
and checking which ports were assigned:
```
$ docker run -P -it -v $(pwd)/data:/mnt/data --name celo-dev celo-dev
```

```
$ docker port celo-dev
```

`$(pwd)/data` directory is where the chain data will be stored.

In case the developer wants to start from a fresh chain, it is enough to remove
chain data: `$ rm -rf data/celo`, or even remove the entire data directory: 
`$ rm -rf data` and start the container again. Keystore will be automatically 
populated and chain genesis executed.
