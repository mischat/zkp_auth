# ZKP Auth Prototype using the Chaun Pederson Algorithm 
by mischa@mmt.me.uk

I put most of my time and effort into trying to understand how the alogorithm works, and to make sure I understood what numbers I was plugging into the system. 

I have never built or looked at ZKPs in that much detail before. I started out with a toy example and then I looked at how I could make the prototype work with a big numbers, and I feel like I got there in the end.

It should be noted that I had to make changes to the protobuf definitions to be able to make the system work with big numbers. You should be able to see this in the git history

## Notes
 - The client created is not an interactive CLI, perhaps one for the future
 - Currently, there is no persistence on either the client or the server
  - Client:
    - I need to ascertain whether or not we can use a incrementing contiguous nonce, whether or not that impacts the security of the system. I will write up how I would approach the further build out of the Client below. 
  - Server: 
    - I have merely implemented a key value store for the server. I would look at using some managed store if I was pushing this out into production. If the choice was to use AWS i would probably be looking DynamoDB for this feature. 
  
## Running the Code

In order to run this prototype ZKP Auth system, you need golang.

There is a server and a client, you need to pass in the public variables to both the client and the server for system to run.

All of the code will respond to `--help` if you would like to get a feel for what parameters a given script / executable might support.

### Setup 

There is a script in `scripts` folder, which you can use to create your own public variables. 

In the below example, the script will find a biggest prime number upto the variable `-p` and it will then generate a Schnorr group, only with some public variables that can be used to configure a client and server. 

```
cd scripts/gennumbers/

go run main.go -p 200000
```

### Running the client and the server
This is how to instantiate with some big numbers:

Need to run please start the `server` and run the `client` from their respective directories: 

```
cd server/ && go run main.go -p "115792089237316195423570985008687907852837564279074904382605163141518161494337" -q "341948486974166000522343609283189" -g "74446558554923317135296388588396736831887322850186029432124219757485062736903" -h "79726485623116979445189935890227226532411986477410367519098002861237945910855"

and 

cd client && go run main.go -p "115792089237316195423570985008687907852837564279074904382605163141518161494337" -q "341948486974166000522343609283189" -g "74446558554923317135296388588396736831887322850186029432124219757485062736903" -h "79726485623116979445189935890227226532411986477410367519098002861237945910855" -u alice0@example.com 
```
## Testing 

There are a handful of unit tests, most of the testing here is to ensure that the numbers are calculated correctly and that the public variables needed to power the ZK auth are indeed sound. 

The key for me was understanding how the public variables were generated so that I could ensure that the numbers flying around actually worked as needed. 

## Future Development 

I would like to touch on the client side and the server side considerations. 

### Client Side Design Decisions

I need to ascertain whether or not the security of the system is compromised if the clients were to uss a incrementing / contiguous nonce, like used in most blockchains for the value K selected by the client during the authentication dance. 

As it stands, the implementation attempts to randomly generate a big number each time, if the number is big enough, this should be OK for our prototype. 

#### Using a "blockchain" styled incrementing / contiguous nonce 

If we were using incrementing nonces, we could make it so that the client needs to store that information themselves, but we could also expose an endpoint on the server's side. The server could store the last used nonce for all of the userIds. 

This seems like something that wouldn't be too onerous on either side.  

#### Using random bigInts

If we couldn't use incrementing nonces, and we need to rely on using random numbers, we should think about maintaining a map of the numbers used.   

Given that we are assuming that the server is acting in good faith and not maliciously, the server could merely maintain a hashMap of all of the challenges `c` that it has ever presented to a given user, some sort of associative array. And that, along with the client randomly picking a big number for `k` should suffice. 

The idea of making the client have to maintain lots of state doesn't seem like a good design decision. 

### Server Side Design Decisions  

As it stands the server uses in memory state, and as a result the state is lost when the service goes down. 

The state should move to something akin to AWS's dynamoDB, that way it would be possible to have a number of different, performant docker containers reading and writing from the same dynamoDB instance. This would allow AWS to take the load. 

### Operational Considerations and next steps

It should be reletively easy to get all of this up and working in docker containers, golang binaries on top of alpine linux images have a nice small footprint and can be deployed almost anywhere. 
