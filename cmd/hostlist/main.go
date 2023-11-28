package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/puttsk/hostlist"
)

func main() {
	var expand bool
	flag.BoolVar(&expand, "expand", false, "Expand hostlist expression")
	flag.BoolVar(&expand, "e", false, "See. -expand")

	var compress bool
	flag.BoolVar(&compress, "compress", false, "Compress list of hostnames to hostlist expression")
	flag.BoolVar(&compress, "c", false, "See. -compress")

	flag.Parse()

	if expand && compress {
		fmt.Println("Must choose either expand (-e) or compress (-c) mode.")
		os.Exit(1)
	}
	// Choose either expand or compress mode
	expand = !compress

	if expand {
		hosts, err := hostlist.Expand(flag.Arg(0))
		if err != nil {
			fmt.Print("Error: " + err.Error())
		}
		fmt.Printf("%s\n", strings.Join(hosts, " "))
	} else if compress {
		hosts := flag.Args()
		expr, _ := hostlist.Compress(hosts)

		fmt.Println(expr)
	}
}
