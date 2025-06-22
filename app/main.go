package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type command string

const (
	exit command = "exit"
	echo command = "echo"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
		}

		cmd = cmd[:len(cmd)-1]
		cmd = strings.TrimSpace(cmd)

		cmdHeader := strings.Split(cmd, " ")[0]
		cmdArgs := strings.Join(strings.Split(cmd, " ")[1:], " ")

		if cmdHeader == string(exit) {
			os.Exit(0)
		}

		if cmdHeader == string(echo) {
			fmt.Fprintln(os.Stdout, cmdArgs)
			continue
		}

		fmt.Fprintln(os.Stderr, cmd+": command not found")
	}
}
