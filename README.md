# HomeTree

Welcome to hometree a library for building **ho**momorphic **me**rkle **tree**s for fast cryptographic 
signing and integrity checking of components in a finite set.

## Use Cases

Hometree has several use cases including but not limited to:

- Integrity checksum for composite data sets including data sets, databases and local and distributed directories
- Fingerprinting reproducible builds 
- Signing software bill of materials (SBOMs)

## Install

```sh
go get -u github.com/KevinFasusi/hometree
```

## Examples

Build a homomorphic Merkle tree from the tokens in a sentence:

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
	fmt.Printf("root hash digest: %s", root.Value)
}
```
```
OUTPUT:
    root hash digest: 5fce597b418510a5e06c87edc5839ca998489b6808dded8f81807b2b1048f2bcb...
```

Compare Merkle Trees using Diff:

```go

	var bs [][]byte
	var bas [][]byte
	var homeTree MerkleTree
	sentence := "Homomorphisms are structure-preserving maps between two algebraic structures"
	altSentence := "Homomorphisms are structure-preserving maps between two algebraic structures allowing a homomorphic " +
		"hash to be updated"
	tokens := strings.Split(sentence, " ")
	for _, token := range tokens {
		bs = append(bs, []byte(token))
	}
	tokens = strings.Split(altSentence, " ")
	for _, token := range tokens {
		bas = append(bas, []byte(token))
	}

	rootS, _ := homeTree.Build(bs)
	in := homeTree.Traverse(InOrder)
	rootAs, _ := homeTree.Build(bas)
	diff, subTree := Diff(rootS, rootAs)
	fmt.Printf("Diff: %v\n", diff)
```

```
Output:

Diff: [4f0e959caed399fa3b09eb8c7e5eac031f9b3c1b17598a5ede4d......]
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
and *muhash* (multiplicative hash) is discussed in this [paper]().

## Homomorphic Merkle Tree

The homomorphic merkle tree combines the merkle tree's ability to collapse the integrity across breadth and depth into a 
single hash value and the speed and efficiency of the randomize-and-combine feature of the homomorphic hash. Combining 
the two provides efficient storage, traversal and insertion of subtrees. Using a relational or graph database it is 
feasible to track the relationship of changes to root hash values due to updates to leaf nodes and therefore the 
relationship of all interior nodes across updates.

## Papers

The following papers informed the implementation in this library:

[1] Ralph C Merkle. A digital signature based on a conventional encryption function. In Conference on the theory and
application of cryptographic techniques, pages 369–378. Springer

[2] Michael Szydlo. Merkle tree traversal in log space and time. In Advances in Cryptology-EUROCRYPT 2004:
International Conference on the Theory and Applications of Cryptographic Techniques, Interlaken, Switzerland,
May 2-6, 2004. Proceedings 23, pages 541–554. Springer, 2004

[3] Mihir Bellare and Daniele Micciancio. A new paradigm for collision-free hashing: Incrementality at reduced cost.
volume 1233. Cryptology-Eurocrypt 97 Proceesings, Lecture Notes in Computer Science, 1997.

[4] Mihir Bellare, O Goldreich, and S Goldwasser. Incremental cryptography: The case of hashing and signing.
volume 839. Cryptology-Eurocrypt 94 Proceesings, Lecture Notes in Computer Science, 1994.
