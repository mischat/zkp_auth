# ZKP Auth Prototype using the Chaun Pederson Algorithm 

I put most of my time and effort into trying to understand how the alogorithm works, as I have never built or looked at ZKPs in that much detail. I started out with a toy example and then i looked at how i could make the prototype work with a big numbers.

It should be noted that I had to make changes to the protobuf definitions to be able to make the system work with big numbers.

## Notes
 - The client created is not an interactive CLI, perhaps one for the future
 - Currently, there is no persistence on either the client or the server
  - Client:
    - I need to ascertain whether or not we can use a incrementing contiguous nonce, whether or not that impacts the security of the system. I will write up how I would approach the further build out of the Client below. 
  - Server: 
    - I have merely implemented a key value store for the server. I would look at using some managed store if I was pushing this out into production. If the choice was to use AWS i would probably be looking DynamoDB for this feature. 
  

## Setup 

In order to run this prototype ZKP Auth system, the the 


## Running the code 
This is how to instantiate with some big numbers:

Need to run these from both the `server` and the `client` directories: 

```
go run main.go -p "115792089237316195423570985008687907852837564279074904382605163141518161494337" -q "341948486974166000522343609283189" -g "74446558554923317135296388588396736831887322850186029432124219757485062736903" -h "79726485623116979445189935890227226532411986477410367519098002861237945910855"

```

## Testing 
