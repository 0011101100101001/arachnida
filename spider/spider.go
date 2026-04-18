package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Italic  = "\033[3m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	White   = "\033[37m"
)

type Config struct {
	isRecursive bool
	depth       uint
	path        string
	url         string
}

type Spider struct {
	config     Config
	visited    map[string]bool
	downloaded map[string]bool
	mutexImage sync.Mutex
	mutexCrawl sync.Mutex
	waitGroup  sync.WaitGroup
}

var (
	regexpImage  = regexp.MustCompile(`(?i)<img[^>]+\bsrc="([^"]+)"`)
	regexpAnchor = regexp.MustCompile(`(?i)<a[^>]+href="([^">]+)"`)
)

func NewSpider(config Config) *Spider {
	return &Spider{
		config:     config,
		downloaded: make(map[string]bool),
		visited:    make(map[string]bool),
	}
}

func (spider *Spider) Run() error {
	spider.PrintConfig()

	if err := os.MkdirAll(spider.config.path, 0755); err != nil {
		return err
	}

	spider.waitGroup.Add(1)
	go func() {
		err := spider.CrawlURL(spider.config.depth, spider.config.url)
		if err != nil {
			fmt.Fprintln(os.Stderr, Bold+Red+"Error:   "+Reset, err)
		}
	}()
	spider.waitGroup.Wait()

	return nil
}

func (spider *Spider) PrintConfig() {
	fmt.Println(Bold + Italic + Magenta + "Spider" + Reset)
	fmt.Println(
		Bold+Blue+"  URL:"+Bold+White,
		strings.TrimPrefix(spider.config.url, "https://"),
	)
	if strings.HasPrefix(spider.config.path, "/") ||
		strings.HasPrefix(spider.config.path, "./") {
		fmt.Println(Bold+Blue+"  Path:"+Bold+White, spider.config.path+Reset)
	} else {
		fmt.Println(
			Bold+Blue+"  Path:"+Bold+White, "./"+spider.config.path+Reset,
		)
	}
	if spider.config.isRecursive {
		fmt.Println(
			Bold + Blue + "  Recursive: " + Bold + Green + "✔" + Reset)
		fmt.Println(
			Bold+Blue+"  Depth:"+Bold+White, spider.config.depth, "\n"+Reset)
	} else {
		fmt.Println(
			Bold + Blue + "  Recursive: " + Bold + Red + "✖\n" + Reset)
	}
}

func (spider *Spider) CrawlURL(recursionDepth uint, rawURL string) error {
	defer spider.waitGroup.Done()

	if recursionDepth == 0 {
		return nil
	}

	spider.mutexCrawl.Lock()
	if spider.visited[rawURL] {
		spider.mutexCrawl.Unlock()
		return nil
	}
	spider.visited[rawURL] = true
	spider.mutexCrawl.Unlock()

	fmt.Println(Magenta+"Visite:  "+Bold+White, rawURL+Reset)

	resp, err := http.Get(rawURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("GET %s: status %s", rawURL, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	imageMatches := regexpImage.FindAllStringSubmatch(string(body), -1)
	anchorMatches := regexpAnchor.FindAllStringSubmatch(string(body), -1)

	baseURL, _ := url.Parse(rawURL)
	if spider.config.isRecursive {
		for _, match := range anchorMatches {
			href := match[1]
			absoluteURL, err := baseURL.Parse(href)
			if err != nil {
				continue
			}
			if absoluteURL.Host == baseURL.Host {
				spider.waitGroup.Add(1)
				go spider.CrawlURL(recursionDepth-1, absoluteURL.String())
			}
		}
	}

	for _, match := range imageMatches {
		imagePath := match[1]
		imageAbsolutePath, err := baseURL.Parse(imagePath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if hasImageExtension(imageAbsolutePath.Path) {
			spider.waitGroup.Add(1)
			go spider.DownloadImage(imageAbsolutePath.String())
		}
	}
	return nil
}

func hasImageExtension(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return true
	}
	return false
}

func (spider *Spider) DownloadImage(imageURL string) error {
	defer spider.waitGroup.Done()

	spider.mutexImage.Lock()
	if spider.downloaded[imageURL] {
		spider.mutexImage.Unlock()
		return nil
	}
	spider.mutexImage.Unlock()

	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Failed:  "+Reset, imageURL)
		return err
	}
	defer resp.Body.Close()

	parts := strings.Split(imageURL, "/")
	fileName := parts[len(parts)-1]
	filePath := spider.config.path + fileName

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Failed:  "+Reset, filePath)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Failed:  "+Reset, imageURL)
		return err
	}
	fmt.Println(Bold+Green+"Download:"+Bold+White, imageURL+Reset)

	spider.mutexImage.Lock()
	spider.downloaded[imageURL] = true
	spider.mutexImage.Unlock()

	return nil
}
