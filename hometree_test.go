package hometree

import (
	"reflect"
	"strings"
	"testing"
)

func TestHomomorphicMerkleTree_Read(t *testing.T) {
	var homeTree MerkleTree
	example1 := []byte("washing up liquid")
	example2 := []byte("soap")
	example3 := []byte("batteries")
	var examples [][]byte
	examples = append(examples, example1)
	examples = append(examples, example2)
	examples = append(examples, example3)
	_, err := homeTree.read(examples)
	if err.Err != nil {
		t.Errorf("Error homeTree.read(): %s", err.Error())
	}
}

func TestHomomorphicMerkleTree_ReadNil(t *testing.T) {
	var homeTree MerkleTree

	examples := [][]byte(nil)
	_, mktErr := homeTree.read(examples)
	if mktErr.Err == nil {
		t.Errorf("Reading from empty slice should cause an exception. %s", mktErr.Err)
	}
}

func TestBalance(t *testing.T) {
	itemsStr := []string{
		"one",
		"two",
		"three",
	}

	itemsByte := make([]byte, 3)

	itemsByte[0] = byte(1)
	itemsByte[1] = byte(2)
	itemsByte[2] = byte(2)

	balancedStr := balance(itemsStr)
	balancedBytes := balance(itemsByte)

	if len(balancedStr)%2 != 0 || len(balancedBytes)%2 != 0 {
		t.Errorf("Length after balancing expected == 4, got blancedStr == %d, balancedBytes == %d",
			len(balancedStr), len(balancedBytes))
	}
}

func TestMerkleTree_Build(t *testing.T) {
	var b [][]byte
	var homeTree MerkleTree
	sentence := "Homomorphisms are structure-preserving maps between two algebraic structures"
	tokens := strings.Split(sentence, " ")
	for _, token := range tokens {
		b = append(b, []byte(token))
	}
	var n *Node
	root, _ := homeTree.Build(b)
	if reflect.TypeOf(root) != reflect.TypeOf(n) {
		t.Errorf("Build not returning a Node struct. Expected *hometree.Node, returned %v", reflect.TypeOf(n))
	}
}

func TestMerkleTree_Traversal(t *testing.T) {
	var b [][]byte
	var homeTree MerkleTree
	sentence := "Homomorphisms are structure-preserving maps between two algebraic structures"
	tokens := strings.Split(sentence, " ")
	for _, token := range tokens {
		b = append(b, []byte(token))
	}

	_, _ = homeTree.Build(b)

	in := homeTree.Traverse(InOrder)
	pre := homeTree.Traverse(PreOrder)
	post := homeTree.Traverse(PostOrder)

	if reflect.DeepEqual(in, pre) || reflect.DeepEqual(in, post) || reflect.DeepEqual(pre, post) {
		t.Errorf("Sorts should not return a slice with equal order ")
	}

	sim := 0
	for _, v := range in {
		for _, k := range pre {
			if v == k {
				sim++
			}
		}
	}

	if sim != len(pre) && sim != len(in) {
		t.Errorf("Elements of the InOrder and PreOrder traversal are not equal")
	}

	sim = 0
	for _, v := range in {
		for _, k := range post {
			if v == k {
				sim++
			}
		}
	}

	if sim != len(in) && sim != len(post) {
		t.Errorf("Elements of the InOrder and PostOrder traversal are not equal")
	}

	sim = 0
	for _, v := range pre {
		for _, k := range post {
			if v == k {
				sim++
			}
		}
	}
	if sim != len(pre) && sim != len(post) {
		t.Errorf("Elements of the InOrder and PreOrder traversal are not equal")
	}

}

func TestDiff(t *testing.T) {
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
	_, subTree := Diff(rootS, rootAs)

	_, subDiffTree := Diff(subTree, rootAs)

	n, _ := Diff(subDiffTree, rootS)

	if !reflect.DeepEqual(in, n) {
		t.Errorf("Subtree diff should equal original Node")
	}
}
