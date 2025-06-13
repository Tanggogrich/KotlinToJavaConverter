package main

import "flag"

func main() {
	var folder string
	flag.StringVar(&folder, "convert", "", "convert kotlin files into java files from concrete directory")
	flag.Parse()

	if folder == "" {
		Convert(folder)
		return
	}
}
