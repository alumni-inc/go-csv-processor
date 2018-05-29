package util

import (
	"log"
	"os"
	"testing"

	"io/ioutil"
	"path/filepath"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const dir = "testdata"

func TestDetermineEncoding(t *testing.T) {

	m := map[string]encoding.Encoding{
		"data.csv":           unicode.UTF8,
	}

	for path, v := range m {
		f, err := os.Open(filepath.Join(dir, path))
		if err != nil {
			log.Fatal(err)
		}
		if enc, err := determineEncoding(f); err != nil {
			log.Fatal(err)
		} else if enc != v {
			log.Fatal("The file: ", path, " should be the ", v, " encoding but not matched.")
		}
	}
}

func TestGetRecords(t *testing.T) {

	// Check all csv files with each different encoding under the testdata directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		//fmt.Println("file -----", f.Name())
		file, err := os.Open(filepath.Join(dir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		records, err := GetRecords(file)
		if err != nil {
			log.Fatal(err)
		}
		// print
		if len(records) != 1000 {
			log.Fatal("records length should be 1000")
		}
	}
}
