package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type command string

const (
	exit command = "exit"
	echo command = "echo"
)

var commands = map[command]func(args string){
	exit: exitC,
	echo: echoC,
}

func exitC(args string) {
	code := 1
	if args != "0" {
		code = 0
	}
	os.Exit(code)
}

func echoC(args string) {
	fmt.Fprintln(os.Stdout, args)
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
	sigs := make(chan os.Signal)
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
	for _ = range sigs {
		os.Exit(0)
	}
}
