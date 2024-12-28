package main

import (
	"fmt"
	"html/template"
	"log"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/bytbox/go-pop3"
)

const (
	pop3Server   = "mail.nic.ru:110"
	pop3Username = "2@dactyl.su"
	pop3Password = "12345678990DactylSUDTS"
)

type Email struct {
	ID      int
	From    string
	Subject string
}

var templates = template.Must(template.New("").ParseFiles("templates/index.html"))

func decodeMimeHeader(encoded string) (string, error) {
	decoder := new(mime.WordDecoder)
	decoded, err := decoder.DecodeHeader(encoded)
	if err != nil {
		return encoded, err
	}
	return decoded, nil
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/delete-all", handleDeleteAll)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectToPOP3() (*pop3.Client, error) {
	conn, err := net.Dial("tcp", pop3Server)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to POP3 server: %w", err)
	}

	client, err := pop3.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create POP3 client: %w", err)
	}

	if err := client.Auth(pop3Username, pop3Password); err != nil {
		client.Quit()
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return client, nil
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	client, err := connectToPOP3()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Quit()

	count, _, err := client.Stat()
	if err != nil {
		http.Error(w, "Failed to fetch email count", http.StatusInternalServerError)
		return
	}

	emails := []Email{}
	for i := 1; i <= count; i++ {
		email, err := fetchEmail(client, i)
		if err != nil {
			log.Printf("Failed to fetch email %d: %v", i, err)
			continue
		}
		emails = append(emails, email)
	}

	err = templates.ExecuteTemplate(w, "index.html", emails)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func fetchEmail(client *pop3.Client, id int) (Email, error) {
	resp, err := client.Retr(id)
	if err != nil {
		return Email{}, err
	}

	msg, err := mail.ReadMessage(strings.NewReader(resp))
	if err != nil {
		return Email{}, err
	}

	from, err := decodeMimeHeader(msg.Header.Get("From"))
	if err != nil {
		from = msg.Header.Get("From")
	}

	subject, err := decodeMimeHeader(msg.Header.Get("Subject"))
	if err != nil {
		subject = msg.Header.Get("Subject")
	}

	return Email{
		ID:      id,
		From:    from,
		Subject: subject,
	}, nil
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	client, err := connectToPOP3()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Quit()

	if err := client.Dele(id); err != nil {
		http.Error(w, "Failed to delete email", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleDeleteAll(w http.ResponseWriter, r *http.Request) {
	client, err := connectToPOP3()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Quit()

	count, _, err := client.Stat()
	if err != nil {
		http.Error(w, "Failed to fetch email count", http.StatusInternalServerError)
		return
	}

	for i := 1; i <= count; i++ {
		if err := client.Dele(i); err != nil {
			log.Printf("Failed to delete email %d: %v", i, err)
			continue
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
