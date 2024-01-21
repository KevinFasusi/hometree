package hometree

import (
	"fmt"
	"lukechampine.com/lthash"
	"sync"
)

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
	value           string
	HomomorphicHash lthash.Hash
	left, right     *Node
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

// Insert a node into a homomorphic merkle tree

// Read constructs leaf nodes from byte slice
func (m *MerkleTree) Read(b []byte) ([]*Node, MerkleTreeError) {
	var nodes []*Node

	if b == nil {
		err := fmt.Errorf("Error byte slices == %d. Length must be at least 1 ", len(b))
		return nil, MerkleTreeError{err}
	}

	b = Balance(b)

	//builds the base leaf nodes
	for _, v := range b {
		node := new(Node)
		node.HomomorphicHash = lthash.New16()
		node.HomomorphicHash.Add([]byte{v})
		node.value = fmt.Sprintf("%X\n", node.HomomorphicHash.Sum(nil))
		nodes = append(nodes, node)
	}

	return nodes, MerkleTreeError{Err: nil}
}

// Build constructs a merkle tree from leaf nodes using homomorphic hash Lthash16 (lattice hash)
func (m *MerkleTree) Build(nodes []*Node) {
	// builds the merkle tree recursively from pairwise hashes
	if m.Hanoi(nodes)[0] != nil && nodes != nil {
		rootNode := m.Hanoi(nodes)[0]
		m.root = rootNode
	} else {
		m.root = nil
	}
}

// Hanoi builds the merkle tree recursively from pairwise hashes
func (m *MerkleTree) Hanoi(nodes []*Node) []*Node {
	var hanoiNodes []*Node
	for i := 0; i < len(nodes); i = i + 2 {
		node := new(Node)
		node.HomomorphicHash = lthash.New16()
		node.HomomorphicHash.Add(nodes[i].HomomorphicHash.Sum(nil))
		node.left = nodes[i]
		node.HomomorphicHash.Add(nodes[i+1].HomomorphicHash.Sum(nil))
		node.right = nodes[i+1]
		node.value = fmt.Sprintf("%X", node.HomomorphicHash.Sum(nil))
		hanoiNodes = append(hanoiNodes, node)
	}
	if len(hanoiNodes) != 1 && nodes != nil {
		hanoiNodes = m.Hanoi(BalanceNodes(hanoiNodes))
	}
	return hanoiNodes
}

// Balance the len of any slice, for use as a leaf node to hash, by ensuring the number of entries in the slice is even.
func Balance[V any](b []V) []V {
	// Balance merkle tree by ensuring the tree has an even number of leaf nodes
	if len(b)%2 != 0 {
		//fmt.Printf("Record: %v\nLength of ngram is not balanced: length==%d\n", b, len(b))
		b = append(b, b[len(b)-1])
		//fmt.Printf("Record balanced: %v\nLength of ngram is balanced: length==%d\n", b, len(b))
	}
	return b
}

// BalanceNodes ensures the number of nodes to hash.
func BalanceNodes(nodes []*Node) []*Node {
	// Balance merkle tree by ensuring the tree has an even number of leaf nodes
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

func (m *MerkleTree) Traverse(traversalType TraversalType) {
	fmt.Printf("%x", m.root.value)
	switch traversalType {
	case InOrder:
		inOrderTraverse(m.root)
	case PreOrder:
		preOrderTraversal(m.root)
	case PostOrder:
		postOrderTraversal(m.root)
	}
}

func inOrderTraverse(node *Node) {
	if node == nil {
		return
	}
	inOrderTraverse(node.left)
	_, err := fmt.Println(node.value)
	if err != nil {
		return
	}
	inOrderTraverse(node.right)
}

func preOrderTraversal(node *Node) {
	if node == nil {
		return
	}
	_, err := fmt.Println(node.value)
	if err != nil {
		return
	}
	preOrderTraversal(node.left)
	preOrderTraversal(node.right)
}

func postOrderTraversal(node *Node) {
	if node == nil {
		return
	}
	preOrderTraversal(node.left)
	preOrderTraversal(node.right)
	_, err := fmt.Println(node.value)
	if err != nil {
		return
	}
}

// Edit finds change and edits existing hmtree
func Edit(hmt *Node, record []string) *Node {
	return hmt
}

// Diff finds the difference between two merkle trees and computes the difference
func Diff(hmt1 *Node, hmt2 *Node) uint16 {
	var diff uint16
	return diff
}
