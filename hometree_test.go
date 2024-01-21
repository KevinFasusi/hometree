package hometree

import (
	"testing"
)

func TestHomomorphicMerkleTree_Read(t *testing.T) {
	var homeTree MerkleTree
	example := []byte("washing up liquid")
	_, err := homeTree.Read(example)
	if err.Err != nil {
		t.Errorf("Error homeTree.Read(): %s", err.Error())
	}
}

func TestHomomorphicMerkleTree_ReadNil(t *testing.T) {
	var homeTree MerkleTree

	example := []byte(nil)
	_, mktErr := homeTree.Read(example)
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

	balancedStr := Balance(itemsStr)
	balancedBytes := Balance(itemsByte)

	if len(balancedStr)%2 != 0 || len(balancedBytes)%2 != 0 {
		t.Errorf("Length after balancing expected == 4, got blancedStr == %d, balancedBytes == %d",
			len(balancedStr), len(balancedBytes))
	}
}
