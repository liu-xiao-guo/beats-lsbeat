package main

// This file is mandatory as otherwise the lsbeat.test binary is not generated correctly.

import (
	"flag"
	"testing"

	"github.com/liu-xiao-guo/lsbeat/cmd"
)

var systemTest *bool

func init() {
	testing.Init()
	systemTest = flag.Bool("systemTest", false, "Set to true when running system tests")

	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup("systemTest"))
	cmd.RootCmd.PersistentFlags().AddGoFlag(flag.CommandLine.Lookup("test.coverprofile"))
}

// Test started when the test binary is started. Only calls main.
func TestSystem(t *testing.T) {

	if *systemTest {
		main()
	}
}
