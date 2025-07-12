package main

import (
	"fmt"
	"os"

	"github.com/patrickdappollonio/vsc-replacer/cmd"
)

func main() {
	if err := cmd.MainCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}
}
