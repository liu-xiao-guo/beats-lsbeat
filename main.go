package main

import (
	"os"

	"github.com/liu-xiao-guo/lsbeat/cmd"

	_ "github.com/liu-xiao-guo/lsbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
