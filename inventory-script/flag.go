package main

import (
	"flag"
)

func flagoption() (bool, string) {
	listPtr := flag.Bool("list", false, "Output all hosts info, works as inventory script")
	hostPtr := flag.String("host", "", "Output specific host info, works as inventory script")
	flag.Parse()
	return *listPtr, *hostPtr
}
