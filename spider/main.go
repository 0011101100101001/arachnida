package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	Default = "\033[0m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	White   = "\033[37m"
)

func parseConfig() (Config, error) {
	recursive := flag.Bool("r", false,
		"-r: recursively downloads the images in a URL received as a parameter.")

	depth := flag.Uint("l", 5,
		"-r -l [N]: indicates the maximum depth level of the recursive download.")

	path := flag.String("p", "./data/",
		"-p [PATH]: indicates the path where the downloaded files will be saved.")

	flag.Parse()
	if flag.NArg() == 0 {
		return Config{}, fmt.Errorf("missing URL")
	}

	url := flag.Args()

	fmt.Println(Bold+Blue+"Url:"+Bold+White, url[0])
	fmt.Println(Bold+Blue+"Recursive:"+Bold+White, *recursive)
	fmt.Println(Bold+Blue+"Depth:"+Bold+White, *depth)
	fmt.Println(Bold+Blue+"Path:"+Bold+White, *path+Default)

	return Config{
		isRecursive: *recursive,
		depth:       *depth,
		path:        *path,
		url:         url[0],
	}, nil
}

func main() {
	log.SetFlags(0)

	config, err := parseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage: ./spider [-rlp] URL")
		os.Exit(2)
	}

	spider := NewSpider(config)
	if err := spider.Run(); err != nil {
		log.Fatalf("spider: %v", err)
	}
}
