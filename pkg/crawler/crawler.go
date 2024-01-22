package crawler

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	Ext  string
	Size int64
	//HashType hashType
}

type Crawler interface {
	Crawl(root string, out io.Writer) error
}

type DirectoryCrawler struct {
	Dir         string
	Regex       []*regexp.Regexp
	Conf        Config
	FileDigests [][]byte
}

func (d *DirectoryCrawler) Crawl(out io.Writer) error {
	return filepath.Walk(d.Dir, d.signatureWalk)
}

func (d *DirectoryCrawler) signatureWalk(path string, info os.FileInfo, err error) error {

	if !info.IsDir() {
		signature := d.FileSignature(path)
		fmt.Printf("FILE Name: %s\nSIG: %X\n", path, signature)
		d.FileDigests = append(d.FileDigests, signature)
		// read each file as bytes and create leaf node
	}

	return nil
}

func (d *DirectoryCrawler) FileSignature(path string) []byte {
	var sum []byte
	file, err := os.Open(path)

	if err != nil {
		fmt.Printf("ERR\n")
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal()
		}
	}(file)

	buf := make([]byte, 8192)
	fileSignature := sha256.New()

	for b := 0; err == nil; {
		b, err = file.Read(buf)
		if err == nil {
			_, err = fileSignature.Write(buf[:b])
		}
	}
	sum = fileSignature.Sum(nil)
	return sum
}

func NewCrawler(root string, conf *Config) *DirectoryCrawler {
	return &DirectoryCrawler{
		Dir:  root,
		Conf: *conf}
}
