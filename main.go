/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/marcy-ot/ddfmt/cmd"
)

func main() {
	cmd.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}
