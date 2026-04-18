package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func parseConfig() (Config, error) {
	recursive := flag.Bool("r", false,
		"recursively downloads the images in a URL received as a parameter.")

	depth := flag.Uint("l", 5,
		"indicates the maximum depth level of the recursive download.")

	path := flag.String("p", "./data/",
		"indicates the path where the downloaded files will be saved.")

	flag.Parse()

	if flag.NArg() == 0 {
		return Config{}, fmt.Errorf("missing URL")
	}

	url := flag.Args()

	rawURL := url[0]
	if !strings.HasPrefix(rawURL, "http://") &&
		!strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	if len(*path) == 0 {
		return Config{}, fmt.Errorf("path cannot be empty")
	} else if (*path)[len(*path)-1] != '/' {
		*path += "/"
	}

	return Config{
		isRecursive: *recursive,
		depth:       *depth,
		path:        *path,
		url:         rawURL,
	}, nil
}

func main() {
	config, err := parseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Error:"+Reset, err.Error())
		fmt.Fprintln(os.Stderr,
			Bold+White+"Usage:"+Reset+"./spider [-rlp] URL")
		os.Exit(2)
	}

	spider := NewSpider(config)
	if err := spider.Run(); err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Error:"+Reset, err.Error())
	}
}
