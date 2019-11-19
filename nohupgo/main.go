package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

/// startWithNoHUP run new process (fork and exec) in Backround
/// with ignoring SIGHUP and detached from TTY
func startWithNoHUP(args []string, stdout, stderr *os.File) (*os.Process, error) {
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()

	if stdout == nil {
		stdout = devnull
	}
	if stderr == nil {
		stderr = devnull
	}

	var err error
	args[0], err = exec.LookPath(args[0])
	if err != nil {
		return nil, err
	}

	// ignoreing SIGHUP
	signal.Ignore(syscall.SIGHUP)

	attr := &os.ProcAttr{
		Files: []*os.File{
			devnull, // stdin
			stdout,
			stderr,
		},
		Sys: &syscall.SysProcAttr{
			Setsid:     true,
			Foreground: false,
		},
	}
	p, err := os.StartProcess(args[0], args, attr)
	if err != nil {
		return nil, err
	}

	return p, err
}

const (
	output  = "nohupgo.txt"
	helpStr = `
nohupgo: nohup implementation in Go

Usage:
	nohupgo sh -c 'sleep 3; echo hello && touch hello'
`
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println(helpStr)
		return
	}

	out, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()
	p, err := startWithNoHUP(os.Args[1:], out, out)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] starting on pid: %d\n", p.Pid)
}
