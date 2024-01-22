package main

import (
	"flag"
	"hometree/pkg/crawler"
	"hometree/pkg/hometree"
	"os"
	"path/filepath"
)

func main() {
	root := flag.String("root", ".", "Root directory")
	flag.Parse()
	path, err := filepath.Abs(*root)
	if err != nil {
		panic("Err")
	}

	config := crawler.Config{
		Ext:  "",
		Size: 0,
	}
	buildHomeTree(path, config)
}

func buildHomeTree(path string, config crawler.Config) {
	crawl := crawler.NewCrawler(path, &config)
	err := crawl.Crawl(os.Stdout)
	if err != nil {
		return
	}
	var homeTree hometree.MerkleTree
	_ = homeTree.BuildHomeTree(crawl.FileDigests)
	homeTree.Traverse(hometree.InOrder)
}
