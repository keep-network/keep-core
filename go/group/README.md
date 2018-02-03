# Overview

Here's the overall (swim lane) process for requesting, generating and publishing a random number:

<img src="https://raw.githubusercontent.com/l3x/images/master/random-num-request.png">

## Steps Summarized:

1. Groups are formed and ready to accept random number requests
2. Client requests a new random number
   - KMA verifies client has purchased enough Keep to process the request
3. Group processes request and generates the random number
   - Each WRK (worker node) performs processing and tries to be the one to submit the number
   - submitted numbers are verified by KMA
   - Each WRK verifies other WRK submissions and my accuse another WRK of wrong doing
4. Group is dissolved

# Components

Here's a list of the components from the diagram above and what each one does.

**CLI** - Client that want us to generate a random number
* request random number

**KMA** - Keep Master that acts as gate keeper between the chain and random#-generating-groups 
* init new group creation
**  set # of nodes in group
**  select nodes that can join
* verify stake (has client purchased enough Keep to participate?)
* verify group ready to process random # requests
** publish group key to chain
   
**WRK** - Nodes that want to join group (worker nodes)     
* request to join group
* generate pub/priv key for group communication
* broadcast accusation (optional)

**GRP** - Group of worker nodes
* submit random number

**ETH** - Ethereum blockchain
* run contract to verify stake
* accepts published random number

# Lex's Playground

This playground will implement **group creation** (as seen in the diagram above)

It takes the [p2ptest](https://github.com/keep-network/go-experiments/tree/master/p2ptest) code and reorganizes it such that each node:

* has it's own IP address
* has three subcomponents
	* initRuntime
	* runWorkers
	* exposeAPI

Each node's API consists of the following:
 
| **API**                | **Subscribers** |
| -----------------------| :-------------: |
| GET peerId             | KMA, WRK, GRP |
| GET peerList           | KMA, WRK, GRP |
| GET groupId            | KMA, WRK, GRP |
| GET stats              | KMA |
| POST requstToJoinGroup | KMA |
| POST signaturePart     | GRP, WRK |
| POST processBeacon     | GRP, WRK |
| POST accusation        | KMA |


POST processBeacon does the following:
* verifies signaturePart
* constructs full BLS signature if we have enough signatureParts 
	* returns the random number

Future work can incorporate the APIs into a pubsub subsystem for the p2p network.

## Playground Status

### Where we are

The code has been going through a lot of initial flux including various dockerizing solutions.

Dep has not played well with my favorite dev tool (jetbrains). Consternation included error messages like the following:

```
unexpected directory layout:
    import path: github.com/libp2p/go-libp2p-crypto
    root: /Users/lex/dev/go/src
    dir: /Users/lex/dev/go/src/github.com/keep-network/go-experiments/p2ptest/vendor/github.com/libp2p/go-libp2p-crypto
    expand root: /Users/lex/dev/go/src
    expand dir: /Users/lex/dev/go/src/github.com/keep-network/go-experiments/p2ptest/vendor/github.com/libp2p/go-libp2p-crypto
    separator: /
```

Errors like that were never part of a Glide dependency management solution. (I am personally unhappy that dep is Go's official choice.)

Got past those dep-related errors/issues and have dep working now.

Will upload working code in the next day or two.

Since this code mostly pertains to **group** creation I'll put it in a group folder.

Will weave in ideas from the code in the beacon folder.

### Where we can go

* The swim lane diagram may not be correct but it helps me visualize the system's components and can be a good point of discussion on Monday.
	* How do clients initially request a random number?
	* Are we using the concept of a Leader node rather than KMA (Keep Master node) to set # of members in group, dissolve group, etc?
	* Is there a need for assigning a "group number"?

* Add pubsub infrastructure

