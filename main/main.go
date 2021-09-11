package main

import (
	"flag"

	"github.com/h0x0er/headlysis"
)

func main() {

	options := new(headlysis.Options)

	flag.StringVar(&options.Url, "url", "", "URL for analysis (use this option only if testing single url)")
	flag.StringVar(&options.UrlFile, "url-file", "", "url file containing urls (use this options if testing multiple urls)")
	flag.StringVar(&options.OutputFile, "output-file", "", "file name for storing result")
	flag.IntVar(&options.Threads, "threads", 10, "number of threads to be used")
	flag.BoolVar(&options.Verbose, "verbose", false, "show errors")

	flag.Parse()

	headlysis.Headlysis(options)

}
