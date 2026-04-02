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
	mutex      sync.Mutex
}

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

	err := s.CrawlURL(s.cfg.depth, s.cfg.url)
	if err != nil {
		return err
	}

	return nil
}

func (s *Spider) CrawlURL(depth uint, rawURL string) error {

	if depth == 0 {
		return nil
	}

	resp, err := http.Get(rawURL)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reImg := regexp.MustCompile(`(?i)<img[^>]+src="([^">]+)"`)
	// reAnchor := regexp.MustCompile(`(?i)<a[^>]+href="([^">]+)"`)
	imgMatches := reImg.FindAllStringSubmatch(string(body), -1)
	// anchorMatches := reAnchor.FindAllStringSubmatch(string(body), -1)

	base, _ := url.Parse(rawURL)
	for _, matche := range imgMatches {
		imgPath := matche[1]
		imgAbs, err := base.Parse(imgPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(imgAbs.Path)
		if hasImageExtension(imgAbs.Path) {
			s.DownloadImage(imgAbs.String())
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
	if s.downloaded[imgURL] {
		return nil
	}

	resp, err := http.Get(imgURL)
	if err != nil {
		fmt.Println("Failed to download:", imgURL, err)
	}
	defer resp.Body.Close()

	parts := strings.Split(imgURL, "/")
	fileName := parts[len(parts)-1]
	filePath := s.cfg.path + "/" + fileName

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", filePath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Failed to save image:", filePath, err)
	} else {
		fmt.Println("Downloaded:", filePath)
	}
	s.downloaded[imgURL] = true
	return nil
}
