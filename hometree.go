package hometree

import (
	"fmt"
	"lukechampine.com/lthash"
	"reflect"
	"sync"
)

type hashType int

const (
	LTHASH hashType = iota
	MUHASH
)

func (h hashType) Strings() string {
	return [...]string{
		"LTHASH",
		"MUHASH",
	}[h]
}

type TraversalType int

const (
	InOrder TraversalType = iota
	PostOrder
	PreOrder
)

func (t TraversalType) Strings() string {
	return [...]string{
		"InOrder",
		"PostOrder",
		"PreOrder",
	}[t]
}

type Node struct {
	Value           string
	HomomorphicHash lthash.Hash
	Left, Right     *Node
}

type MerkleTree struct {
	root *Node
	lock sync.RWMutex
}

type MerkleTreeError struct {
	Err error
}

func (e MerkleTreeError) Error() string {
	return "MerkleTree Error:" + e.Error()
}

// read constructs leaf nodes from byte slice
func (m *MerkleTree) read(b [][]byte) ([]*Node, MerkleTreeError) {
	var nodes []*Node

	if b == nil {
		err := fmt.Errorf("Error byte slices == %d. Length must be at least 1 ", len(b))
		return nil, MerkleTreeError{err}
	}

	b = balance(b)

	//builds the base leaf nodes
	for _, block := range b {
		node := new(Node)
		node.HomomorphicHash = lthash.New16()
		node.HomomorphicHash.Add(block)
		node.Value = fmt.Sprintf("%x\n", node.HomomorphicHash.Sum(nil))
		//fmt.Printf("Leaf Node:%v\n", node.Value)
		nodes = append(nodes, node)
	}

	return nodes, MerkleTreeError{Err: nil}
}

// Build a merkle tree from a slice of a slice of bytes
func (m *MerkleTree) Build(b [][]byte) (*Node, MerkleTreeError) {

	if b == nil {
		err := fmt.Errorf("Error byte slices == %d. Length must be at least 1 ", len(b))
		return nil, MerkleTreeError{err}
	}
	nodes, err := m.read(b)
	if err.Err != nil {
		return nil, MerkleTreeError{err}
	}
	m.tree(nodes)
	return m.root, MerkleTreeError{nil}
}

// tree constructs a merkle tree from leaf nodes using homomorphic hash Lthash16 (lattice hash) to create interior nodes
func (m *MerkleTree) tree(nodes []*Node) {
	// builds the merkle tree recursively from pairwise hashes
	if m.hanoi(nodes)[0] != nil && nodes != nil {
		rootNode := m.hanoi(nodes)[0]
		m.root = rootNode
	} else {
		m.root = nil
	}
}

// hanoi builds the merkle tree recursively from pairwise hashes unitl a root hash is generated
func (m *MerkleTree) hanoi(nodes []*Node) []*Node {
	var hanoiNodes []*Node
	for i := 0; i < len(nodes); i = i + 2 {
		node := new(Node)
		node.HomomorphicHash = lthash.New16()
		node.HomomorphicHash.Add(nodes[i].HomomorphicHash.Sum(nil))
		node.Left = nodes[i]
		node.HomomorphicHash.Add(nodes[i+1].HomomorphicHash.Sum(nil))
		node.Right = nodes[i+1]
		node.Value = fmt.Sprintf("%x", node.HomomorphicHash.Sum(nil))
		//fmt.Printf("Interior Node: %v\n", node.Value)
		hanoiNodes = append(hanoiNodes, node)
	}
	if len(hanoiNodes) != 1 && nodes != nil {
		hanoiNodes = m.hanoi(m.balanceNodes(hanoiNodes))
	}
	return hanoiNodes
}

// balance the len of any slice, for use as a node to hash, by ensuring the number of entries in the slice is even.
func balance[V any](b []V) []V {
	if len(b)%2 != 0 {
		//fmt.Printf("Record: %v\nLength of ngram is not balanced: length==%d\n", b, len(b))
		b = append(b, b[len(b)-1])
		//fmt.Printf("Record balanced: %v\nLength of ngram is balanced: length==%d\n", b, len(b))
	}
	return b
}

// balanceNodes ensures the merkle tree even number of nodes at height h_n
func (m *MerkleTree) balanceNodes(nodes []*Node) []*Node {
	if len(nodes)%2 != 0 {
		//fmt.Printf("Record: %v\nLength of node is not balanced: length==%d\n", nodes, len(nodes))
		var node, oNode []*Node
		oNode = append(oNode, nodes[len(nodes)-1])
		_ = copy(oNode, node)
		nodes = append(nodes, oNode[0])
		//fmt.Printf("Record balanced: %v\nLength of node is balanced: length==%d\n", nodes, len(nodes))
	}
	return nodes
}

// Traverse the merkle tree
func (m *MerkleTree) Traverse(traversalType TraversalType) []string {
	//fmt.Printf("%x", m.root.Value) var v []string
	var digests []string

	switch traversalType {
	case InOrder:
		inOrderTraverse(m.root, &digests)
	case PreOrder:
		preOrderTraversal(m.root, &digests)
	case PostOrder:
		postOrderTraversal(m.root, &digests)
	}
	return digests
}

// Insert adds a new leaf node to exiting tree propagates the changes through the merkle tree
func (m *MerkleTree) Insert() {

}

// Remove leaf node from existing tree and propagates changes through the merkle tree
func (m *MerkleTree) Remove() {}

// Diff compare two merkle trees and returns the difference as a slice and a subtree
func Diff(hmt1 *Node, hmt2 *Node) ([]string, *Node) {
	var hmt1Digests []string
	var hmt2Digests []string
	var diff []string
	var subT *Node

	inOrderTraverse(hmt1, &hmt1Digests)
	inOrderTraverse(hmt2, &hmt2Digests)

	eq := reflect.DeepEqual(hmt1Digests, hmt2Digests)

	if eq {
		return diff, subT
	}

	var remove []string

	for _, i := range hmt2Digests {
		for _, j := range hmt1Digests {
			if i == j {
				remove = append(remove, j)
			}
		}
	}

	for _, hmt1D := range hmt1Digests {
		for n, hmt2D := range hmt2Digests {
			for _, k := range remove {
				if hmt1D == k && hmt1D == hmt2D {
					hmt2Digests = append(hmt2Digests[:n], hmt2Digests[n+1:]...)
				}
			}
		}
	}

	diff = hmt2Digests
	subT = diffSubTree(hmt2Digests, hmt2)
	return diff, subT
}

// diffSubTree locates the lowest depth node with the same digest found in the diff and returns the subtree
func diffSubTree(diff []string, node *Node) *Node {
	if node == nil {
		return nil
	}
	diffSubTree(diff, node.Left)
	for _, i := range diff {
		if i == node.Value {
			return node
		}
	}
	diffSubTree(diff, node.Right)
	return nil
}

func inOrderTraverse(node *Node, digests *[]string) {

	if node == nil {
		return
	}
	inOrderTraverse(node.Left, digests)
	*digests = append(*digests, node.Value)
	inOrderTraverse(node.Right, digests)
}

func preOrderTraversal(node *Node, digests *[]string) {
	if node == nil {
		return
	}
	*digests = append(*digests, node.Value)
	preOrderTraversal(node.Left, digests)
	preOrderTraversal(node.Right, digests)
}

func postOrderTraversal(node *Node, digests *[]string) {
	if node == nil {
		return
	}
	preOrderTraversal(node.Left, digests)
	preOrderTraversal(node.Right, digests)
	*digests = append(*digests, node.Value)
}
