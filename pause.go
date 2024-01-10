package bnp

import (
	"bufio"
	"fmt"
	"os"
)

// Pause pauses the CLI, issuing a message.
func Pause(msg string) {
	Pausef("%s\n", msg)
}

// Pausef pauses the CLI, issuing a formatted string.
func Pausef(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// Pauseln pauses the CLI, issuing a list of arguments.
func Pauseln(a ...any) {
	Pausef("%s\n", fmt.Sprintln(a...))
}
