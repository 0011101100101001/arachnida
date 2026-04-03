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
	Default = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	White   = "\033[37m"
	Header  = "░▒█▀▀▀█░▒█▀▀█░▀█▀░▒█▀▀▄░▒█▀▀▀░▒█▀▀▄\n" +
		"░░▀▀▀▄▄░▒█▄▄█░▒█░░▒█░▒█░▒█▀▀▀░▒█▄▄▀\n" +
		"░▒█▄▄▄█░▒█░░░░▄█▄░▒█▄▄█░▒█▄▄▄░▒█░▒█"
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
			fmt.Println("Error starting crawl:", err)
		}
	}()
	spider.waitGroup.Wait()

	return nil
}

func (spider *Spider) PrintConfig() {
	fmt.Println(Bold + Magenta + Header + "\n" + Default)
	fmt.Println(Bold+Blue+"URL:"+Bold+White, spider.config.url)
	fmt.Println(Bold+Blue+"Recursive:"+Bold+White, spider.config.isRecursive)
	fmt.Println(Bold+Blue+"Depth:"+Bold+White, spider.config.depth)
	fmt.Println(Bold+Blue+"Path:"+Bold+White, spider.config.path+"\n"+Default)
}

func (spider *Spider) CrawlURL(recursionDepth uint, rawURL string) error {
	defer spider.waitGroup.Done()

	spider.mutexCrawl.Lock()
	if spider.visited[rawURL] || recursionDepth == 0 {
		spider.mutexCrawl.Unlock()
		return nil
	}
	spider.visited[rawURL] = true
	spider.mutexCrawl.Unlock()

	fmt.Println(Bold+Magenta+"Visite:  "+Bold+White, rawURL+Default)

	resp, err := http.Get(rawURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, Red+"Error:"+Default, err)
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
		fmt.Fprintln(os.Stderr, Bold+Red+"Failed:  "+Default, imageURL)
		return err
	}
	defer resp.Body.Close()

	parts := strings.Split(imageURL, "/")
	fileName := parts[len(parts)-1]
	filePath := spider.config.path + fileName

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Failed:  "+Default, filePath)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, Bold+Red+"Failed:  "+Default, imageURL)
		return err
	}
	fmt.Println(Bold+Green+"Download:"+Bold+White, imageURL+Default)

	spider.mutexImage.Lock()
	spider.downloaded[imageURL] = true
	spider.mutexImage.Unlock()

	return nil
}
