package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	exitMessage = "exit 0"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
		}

		command = command[:len(command)-1]
		command = strings.TrimSpace(command)

		if command == exitMessage {
			os.Exit(0)
		}

		fmt.Println(command + ": command not found")
	}
}
