package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type queue struct {
	items []string
}

func newQueue() *queue {
	return &queue{}
}

func (q *queue) push(item string) {
	q.items = append(q.items, item)
}

func (q *queue) pop() string {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *queue) isEmpty() bool {
	return len(q.items) == 0
}

type syncSet struct {
	items map[string]struct{}
}

func newSyncMap() *syncSet {
	return &syncSet{
		items: make(map[string]struct{}),
	}
}

func (m *syncSet) insert(key string) {
	m.items[key] = struct{}{}
}

func (m *syncSet) contains(key string) bool {
	_, ok := m.items[key]
	return ok
}

type flags struct {
	mirror    bool
	userAgent string
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"

func main() {
	f := &flags{}

	flag.BoolVar(&f.mirror, "m", false, "скачать сайт на локальную машину")
	flag.StringVar(&f.userAgent, "U", userAgent, "заголовок User-Agent")

	flag.Parse()

	if len(flag.Args()) < 1 {
		_, _ = fmt.Fprintln(os.Stderr, "Необходимо указать URL: wget [flags] https://example.com")
		os.Exit(1)
	}

	args := flag.Args()
	argURL := args[0]
	baseURL, err := url.Parse(argURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Ошибка при разборе URL: %v", err)
		return
	}

	body, err := downloadURL(argURL, f)
	if err != nil {
		log.Printf("Ошибка при скачивании URL: %v", err)
		return
	}

	urls := newSyncMap()
	q := newQueue()

	saveLinkAndAddNewLinks := func(link string) error {
		absURL, err := resolveURL(baseURL, link)
		if err != nil {
			return fmt.Errorf("resolveURL: %w", err)
		}
		body, err = downloadURL(absURL, f)
		if err != nil {
			return fmt.Errorf("downloadURL: %w", err)
		}

		err = extractLinks(body, q, urls, baseURL.Host)
		if err != nil {
			return fmt.Errorf("extractLinks: %w", err)
		}
		return nil
	}

	if f.mirror {
		err = extractLinks(body, q, urls, baseURL.Host)
		if err != nil {
			log.Printf("extractLinks: %v", err)
			return
		}

		for !q.isEmpty() {
			link := q.pop()
			err := saveLinkAndAddNewLinks(link)
			if err != nil {
				log.Printf("saveLinkAndAddNewLinks: %v", err)
				return
			}
		}
	}
}

func downloadURL(urlStr string, f *flags) ([]byte, error) {
	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("oшибка при разборе URL: %w", err)
	}

	client := &http.Client{}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    baseURL,
		Header: make(http.Header),
	}
	req.Header.Set("User-Agent", f.userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oшибка при выполнении GET-запроса: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("oшибка при чтении ответа: %w", err)
	}

	// Определяем путь для сохранения файла
	filePath, err := getFilePath(baseURL, urlStr)
	if err != nil {
		return nil, fmt.Errorf("oшибка при определении пути: %w", err)
	}
	filePath = filepath.Join(filePath, "index.html")
	// Создаем директорию для сохранения файла
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("oшибка при создании директории: %w", err)
	}

	// Сохраняем содержимое файла
	err = os.WriteFile(filePath, body, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("oшибка при сохранении файла: %w", err)
	}

	log.Printf("Файл успешно скачан и сохранен: %v", filePath)

	return body, nil
}

func getFilePath(baseURL *url.URL, urlStr string) (string, error) {
	relURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе URL: %w", err)
	}

	// Определяем относительный путь файла
	fileName := baseURL.ResolveReference(relURL).Path

	// Добавляем базовую директорию и имя файла
	filePath := filepath.Join("website", fileName)

	return filePath, nil
}

func extractLinks(body []byte, q *queue, urls *syncSet, host string) error {

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("ошибка при разборе HTML: %w", err)
	}

	// Добавляем ссылки в очередь
	var addLinkInQueue func(*html.Node)
	addLinkInQueue = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					ok := urls.contains(attr.Val)
					if !ok && (strings.HasPrefix(attr.Val, "/") || strings.Contains(attr.Val, host)) {
						urls.insert(attr.Val)
						q.push(attr.Val)
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			addLinkInQueue(c)
		}
	}

	addLinkInQueue(doc)
	return nil
}

func resolveURL(baseURL *url.URL, urlStr string) (string, error) {
	relURL, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("ошибка при разборе URL: %w", err)
	}

	absURL := baseURL.ResolveReference(relURL).String()
	return absURL, nil
}
