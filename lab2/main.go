package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Структура для хранения данных
type PageData struct {
	Title string
	Links []struct {
		Text string
		Href string
	}
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Println("Сервер запущен на http://185.102.139.168:8080") // Замените на ваш IP-адрес
	log.Fatal(http.ListenAndServe("185.102.139.168:8080", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	url := "https://news.rambler.ru/"

	// Выполняем GET-запрос
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Проверяем статус-код ответа
	if res.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("Статус-код ответа: %d", res.StatusCode), res.StatusCode)
		return
	}

	// Загружаем содержимое страницы в goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем страницу данных
	var pageData PageData
	pageData.Title = "bfkdbje"
	newsBlocks := doc.Find(".XSvLK2D0")

	// Iterate over each news block
	newsBlocks.Each(func(i int, s *goquery.Selection) {
		// Get the link
		link, exists := s.Find("a").Attr("href")
		if !exists {
			fmt.Println("No link found")
			return
		}

		// Get the title
		title := s.Find(".lNJ9PP5h").Text()
		pageData.Links = append(pageData.Links, struct {
			Text string
			Href string
		}{Text: title, Href: link})
		// Print the news
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Link: %s\n\n", link)
	})
	newsBlock := doc.Find(".CUzJFxJK")
	// Iterate over each news item
	newsBlock.Each(func(i int, s *goquery.Selection) {
		// Get the link
		link, exists := s.Find("a").Attr("href")
		if !exists {
			fmt.Println("No link found")
			return
		}

		// Get the title
		title := s.Find(".fPnVl30V").Text()
		pageData.Links = append(pageData.Links, struct {
			Text string
			Href string
		}{Text: title, Href: link})
		// Print the news
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Link: %s\n\n", link)
	})
	newsBlock.Each(func(i int, s *goquery.Selection) {
		// Get the link
		link, exists := s.Find("a").Attr("href")
		if !exists {
			fmt.Println("No link found")
			return
		}

		// Get the title
		title := s.Find(".fPnVl30V").Text()
		pageData.Links = append(pageData.Links, struct {
			Text string
			Href string
		}{Text: title, Href: link})
		// Print the news
		fmt.Printf("Title: %s\n", title)
		fmt.Printf("Link: %s\n\n", link)
	})
	// Извлекаем ссылки
	/*doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		text := s.Text()
		pageData.Links = append(pageData.Links, struct {
			Text string
			Href string
		}{Text: text, Href: link})
	})*/

	// Генерируем HTML
	t, err := template.New("webpage").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>{{.Title}}</title>
		</head>
		<body>
			<h1>{{.Title}}</h1>
			<h2>Ссылки:</h2>
			<ul>
				{{range .Links}}
					<li><a href="{{.Href}}">{{.Text}}</a></li>
				{{end}}
			</ul>
		</body>
		</html>
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Выполняем шаблон
	if err := t.Execute(w, pageData); err != nil {
		log.Fatal(err)
	}
}
