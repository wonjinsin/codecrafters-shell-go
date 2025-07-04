package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type command string

const (
	exit  command = "exit"
	echo  command = "echo"
	type_ command = "type"
)

var commands = map[command]func(args string){
	exit:  exitC,
	echo:  echoC,
	type_: typeC,
}

var builtins = map[command]bool{
	exit:  true,
	echo:  true,
	type_: true,
}

var path = os.Getenv("PATH")

func exitC(args string) {
	code := 1
	if args == "0" {
		code = 0
	}
	os.Exit(code)
}

func echoC(args string) {
	fmt.Fprintln(os.Stdout, args)
}

func typeC(args string) {
	if builtins[command(args)] {
		fmt.Fprintf(os.Stdout, "type: %s is a shell builtin\n", args)
		return
	}

	for p := range strings.SplitSeq(path, ":") {
		fullPath := filepath.Join(p, args)
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			if fileInfo.Mode()&0111 != 0 {
				fmt.Fprintf(os.Stdout, "%s: not found\n", args)
				return
			}
			fmt.Fprintf(os.Stdout, "%s is %s\n", args, fullPath)
			return
		}
		if os.IsNotExist(err) {
			continue
		}
		fmt.Fprintf(os.Stdout, "%s: unknown error, %s\n", args, err)
	}
	fmt.Fprintf(os.Stdout, "%s: not found\n", args)
}

func handle(msg string) {
	parts := strings.Split(msg, " ")
	hdr := parts[0]
	args := strings.Join(parts[1:], " ")

	if _, ok := commands[command(hdr)]; !ok {
		fmt.Fprintln(os.Stderr, hdr+": command not found")
		return
	}

	commands[command(hdr)](args)
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleSignal(sigs)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			continue
		}

		cmd = cmd[:len(cmd)-1]
		cmd = strings.TrimSpace(cmd)
		handle(cmd)
	}
}

func handleSignal(sigs chan os.Signal) {
	for range sigs {
		os.Exit(0)
	}
}
