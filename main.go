package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	var targetURL string
	flag.StringVar(&targetURL, "url", "", "URL для сканирования")
	flag.Parse()

	if targetURL == "" {
		fmt.Println("Использование: go run main.go -url <URL>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Попытка получить ссылки с: %s\n", targetURL)

	resp, err := http.Get(targetURL)
	if err != nil {
		fmt.Printf("Ошибка при выполнении HTTP-запроса к %s: %v\n", targetURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Получен неожиданный статус HTTP %s от %s\n", resp.Status, targetURL)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа от %s: %v\n", targetURL, err)
		return
	}

	htmlContent := string(bodyBytes)

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Printf("Ошибка при парсинге HTML с %s: %v\n", targetURL, err)
		return
	}

	baseURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("Ошибка при парсинге базового URL %s: %v\n", targetURL, err)
		return
	}

	links := findLinks(doc, baseURL)

	if len(links) == 0 {
		fmt.Println("Ссылок не найдено.")
		return
	}

	fmt.Println("\nНайденные ссылки:")
	for i, link := range links {
		fmt.Printf("%d. %s\n", i+1, link)
	}
}

func findLinks(node *html.Node, base *url.URL) []string {
	var links []string

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				resolvedURL, err := base.Parse(attr.Val)
				if err != nil {
					continue
				}
				links = append(links, resolvedURL.String())
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, findLinks(child, base)...)
	}

	return links
}
