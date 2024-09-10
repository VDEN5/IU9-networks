package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SlyMarbo/rss"
)

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	feed, _ := rss.Fetch("http://www.kommersant.ru/RSS/main.xml")
	fmt.Println(feed.Author)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
 <!DOCTYPE html>
 <html>
 <head>
  <title>Пример</title>
 </head>
 <body>
  <h1>Добро пожаловать на наш сайт!</h1>
  <p>Чтобы перейти к странице, нажмите на ссылку ниже:</p>
  <a href="%s">Перейти по ссылке</a>
  <br>
  <p>Или можете щелкнуть <a href="https://www.example.com">%s</a></p>
  <p>Или можете щелкнуть <a href="https://www.example.com">%s</a></p>
  <p>Или можете щелкнуть <a href="https://www.example.com">%s</a></p>
 </body>
 </html>
 `
	html = fmt.Sprintf(html, feed.Link, feed.Items[0].Title, feed.Items[1].Title, feed.Items[2].Title)
	fmt.Fprintf(w, html)
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)
	err := http.ListenAndServe("185.102.139.168:9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

type Feed struct {
	Nickname    string // This is not set by the package, but could be helpful.
	Title       string
	Description string
	Link        string // Link to the creator's website.
	UpdateURL   string // URL of the feed itself.
	Image       *Image // Feed icon.
	Items       []*Item
	ItemMap     map[string]struct{} // Used in checking whether an item has been seen before.
	Refresh     time.Time           // Earliest time this feed should next be checked.
	Unread      uint32              // Number of unread items. Used by aggregators.
}

type Item struct {
	Title     string
	Summary   string
	Content   string
	Link      string
	Date      time.Time
	DateValid bool
	ID        string
	Read      bool
}

type Image struct {
	Title  string
	URL    string
	Height uint32
	Width  uint32
}
