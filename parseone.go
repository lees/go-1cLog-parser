package main

import (
	"flag"
	"fmt"
	"strings"
	//"runtime"
)

func main() {

	//runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Need a file to process!")
		return
	}

	events := ConvertFile(flag.Args()[0])

	for event := range events {

		//fmt.Println(event[0])

		if event[8] != "E" {
			continue
		}

		str := "{" + strings.Join(event, ",") + "},\n"
		fmt.Print(str + "\n")

	}

}
