package main

import "fmt"
import "flag"
import "multithread_du"

func main() {
	flag.Parse()
	var args = flag.Args()
	var filename = "."

	if len(args) > 0 {
		filename = args[0]
	}

	first_filename := filename

	total_size := multithread_du.TotalFileSize(first_filename)

	fmt.Printf("%d\t%s\n", total_size, first_filename)
}
