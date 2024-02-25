# HomeTree

Welcome to hometree a library for building **ho**momorphic **me**rkle **tree**s for fast cryptographic 
signing and integrity checking of components in a finite set.

## Use Cases

Hometree has several use cases including but not limited to:

- Fingerprinting reproducible builds 
- Signing SBOMs
- Integrity checksum for composite data sets including data sets, databases and local and distributed directories

## Install

```sh
get -u github.com/KevinFasusi/hometree
```

## Examples

Build a homomorphic merkle tree from the tokens in a sentence:

```go
    package main

import (
	"fmt"
	"github.com/KevinFasusi/hometree"
	"strings"   
)

func main() {
	var b [][]byte
	var homeTree hometree.MerkleTree
	sentence := "Homomorphisms are structure-preserving maps between two algebraic structures"
	tokens := strings.Split(sentence, " ")
	for _, token := range tokens {
		b = append(b, []byte(token))
	}

	root := homeTree.Build(b)
	fmt.Printf("Root Checksum: %s", root.value)
}

```
```
OUTPUT:
    root hash: 5fce597b418510a5e06c87edc5839ca998489b6808dded8f81807b2b1048f2bcb...
```

## Merkle Tree

A Merkle tree is a binary tree with a height **H**. The nodes forming the base of the tree at **H**=0 are the leaf nodes. 
Traditionally, the leaf nodes are chosen arbitrarily or by using a pseudo random number. The hash of the pseudo
random number is referred to as the leaf-preimage. All nodes derived from the leaf nodes are named interior nodes.
For a tree of height **H**, there are 2^**H** leaves and 2^**H**-1 interior nodes [1][2]. 

The height **h** of an interior node is relative to its distance in the path to a leaf node below. 
The Merkle tree uses a hash function and values assigned to each leaf node, to hash all leaf and interior nodes until 
a single root hash is generated.

The root hash created from the leaf and intermediate interior nodes can be used as an integrity check over the set.

## Homomorphic Hash

The property of homomorphism is present in a map between two algebraic structures that preserve their operations. 
In cryptography homomorphism is often associated with fully homomorphic encryption. However, Bellare, Goldreich and 
Goldwasser [4] detail the properties and construction of a hashing mechanism with partial updates or "incrementality". 
The hashing scheme allows the for the update of a digest (through some operation dependent on the construction
of the homomorphic hash) when the original hashed message has been altered without re-hashing the entire message. 
The method uses already approved cryptographic hash functions such as BLAKE2b etc.  

The homomorphic hash used to build this merkle tree is lthash by [Luke Champine](https://github.com/lukechampine/lthash). 
An implementation of [MuHash](https://github.com/KevinFasusi/muHash) (currently a work in progress) will be added to 
this library as another option. For more about homomorphic hashing and the difference between *ltHash* (lattice hash) 
and *muhash* (multiplicative hash) is discussed [here]().


## Papers

The following papers informed the implementation in this library:

[1] Ralph C Merkle. A digital signature based on a conventional encryption function. In Conference on the theory and
application of cryptographic techniques, pages 369–378. Springe

[2] Michael Szydlo. Merkle tree traversal in log space and time. In Advances in Cryptology-EUROCRYPT 2004:
International Conference on the Theory and Applications of Cryptographic Techniques, Interlaken, Switzerland,
May 2-6, 2004. Proceedings 23, pages 541–554. Springer, 2004

[3] Mihir Bellare and Daniele Micciancio. A new paradigm for collision-free hashing: Incrementality at reduced cost.
volume 1233. Cryptology-Eurocrypt 97 Proceesings, Lecture Notes in Computer Science, 1997.

[4] Mihir Bellare, O Goldreich, and S Goldwasser. Incremental cryptography: The case of hashing and signinng.
volume 839. Cryptology-Eurocrypt 94 Proceesings, Lecture Notes in Computer Science, 1994
r, 1987.
