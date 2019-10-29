package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {
	err := exec.Command("tar", "cf", "ex.tar", "ex").Run()
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile("ex.tar")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("---BEFORE---")
	buf := bytes.NewBuffer(data)
	tr := tar.NewReader(buf)
	var thdr *tar.Header
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if hdr.FileInfo().IsDir() {
			continue
		}

		content := make([]byte, hdr.Size)
		tr.Read(content)
		fmt.Printf("\t%s: %s\n", hdr.Name, string(content))
		if hdr.Name == "ex/target" {
			thdr = hdr
		}
	}

	fmt.Println("---MODIFYING---")
	newContent := []byte("this is modified data")
	buf = bytes.NewBuffer(data)
	tw := tar.NewWriter(buf)
	thdr.Size = int64(len(newContent))
	err = tw.WriteHeader(thdr)
	if err != nil {
		log.Fatal(err)
	}
	wrote, err := tw.Write(newContent)
	if err != nil {
		log.Fatal(err)
	}
	tw.Flush()
	tw.Close()
	fmt.Printf("\t%d BYTES WROTE\n", wrote)

	fmt.Println("---AFTER---")
	tr = tar.NewReader(buf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if hdr.FileInfo().IsDir() {
			continue
		}

		content := make([]byte, hdr.Size)
		tr.Read(content)
		fmt.Printf("\t%s: %s\n", hdr.Name, string(content))
	}
}
