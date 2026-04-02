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
	White   = "\033[37m"
)

type Config struct {
	isRecursive bool
	depth       uint
	path        string
	url         string
}

type Spider struct {
	cfg        Config
	visited    map[string]bool
	downloaded map[string]bool
	muImage    sync.Mutex
	muCrawl    sync.Mutex
	wg         sync.WaitGroup
}

var (
	regexpImg    = regexp.MustCompile(`(?i)<img[^>]+\bsrc="([^"]+)"`)
	regexpAnchor = regexp.MustCompile(`(?i)<a[^>]+href="([^">]+)"`)
)

func NewSpider(cfg Config) *Spider {
	return &Spider{
		cfg:        cfg,
		downloaded: make(map[string]bool),
		visited:    make(map[string]bool),
	}
}

func (s *Spider) Run() error {
	if err := os.MkdirAll(s.cfg.path, 0755); err != nil {
		return err
	}

	s.wg.Add(1)
	go func() {
		err := s.CrawlURL(s.cfg.depth, s.cfg.url)
		if err != nil {
			fmt.Println("Error starting crawl:", err)
		}
	}()
	s.wg.Wait()

	return nil
}

func (s *Spider) CrawlURL(recursionDepth uint, rawURL string) error {
	defer s.wg.Done()

	s.muCrawl.Lock()
	if s.visited[rawURL] || recursionDepth == 0 {
		s.muCrawl.Unlock()
		return nil
	}
	s.visited[rawURL] = true
	s.muCrawl.Unlock()

	fmt.Println(Bold + Yellow + rawURL + Default)

	resp, err := http.Get(rawURL)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		 return fmt.Errorf("GET %s: status %s", rawURL, resp.Status)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	imgMatches := regexpImg.FindAllStringSubmatch(string(body), -1)
	anchorMatches := regexpAnchor.FindAllStringSubmatch(string(body), -1)

	baseURL, _ := url.Parse(rawURL)
	if s.cfg.isRecursive {
		for _, match := range anchorMatches {
			href := match[1]
			absURL, err := baseURL.Parse(href)
			if err != nil {
				continue
			}
			if absURL.Host == baseURL.Host {
				s.wg.Add(1)
				go s.CrawlURL(recursionDepth-1, absURL.String())
			}
		}
	}

	for _, match := range imgMatches {
		imgPath := match[1]
		imgAbs, err := baseURL.Parse(imgPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if hasImageExtension(imgAbs.Path) {
			s.wg.Add(1)
			go s.DownloadImage(imgAbs.String())
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

func (s *Spider) DownloadImage(imgURL string) error {
	defer s.wg.Done()

	s.muImage.Lock()
	if s.downloaded[imgURL] {
		s.muImage.Unlock()
		return nil
	}
	s.muImage.Unlock()

	resp, err := http.Get(imgURL)
	if err != nil {
		fmt.Println("Failed to download:", imgURL, err)
		return err
	}
	defer resp.Body.Close()

	parts := strings.Split(imgURL, "/")
	fileName := parts[len(parts)-1]
	filePath := s.cfg.path + fileName

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", filePath, err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Failed to save image:", imgURL, err)
		return err
	}
	fmt.Println(Blue+"Downloaded:", imgURL+Default)

	s.muImage.Lock()
	s.downloaded[imgURL] = true
	s.muImage.Unlock()

	return nil
}
