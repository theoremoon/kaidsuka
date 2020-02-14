package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func run() error {
	flag.Usage = func() {
		fmt.Printf("Usage:\n  %s [OPTIONS] <git repo URL>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	deploykey := flag.String("deploykey", "", "GitHub Deploykey.")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return nil
	}

	fs := memfs.New()
	storer := memory.NewStorage()

	if deploykey != nil && *deploykey != "" {
		if !strings.HasPrefix(args[0], "git@") {
			return fmt.Errorf("repository url format: git@<hostingsite.com>:<username>/<repo>.git")
		}
		auth, err := ssh.NewPublicKeysFromFile("git", *deploykey, "")
		if err != nil {
			return err
		}
		_, err = git.Clone(storer, fs, &git.CloneOptions{
			URL:  args[0],
			Auth: auth,
		})
		if err != nil {
			return err
		}
	} else {
		_, err := git.Clone(storer, fs, &git.CloneOptions{
			URL: args[0],
		})
		if err != nil {
			return err
		}
	}

	infos, err := fs.ReadDir("/")
	if err != nil {
		return err
	}
	for _, info := range infos {
		fmt.Println(info.Name())
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
