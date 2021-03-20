package bnp

import (
	"bufio"
	"fmt"
	"os"
)

// Pause pauses the CLI
func Pause(msg string) {
	fmt.Println(msg)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
