package main

import (
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Создаем новый запрос для целевого сервера
	targetURL := "http://www.gnuplot.info/" + r.URL.Path // Замените на нужный URL

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из оригинального запроса
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Отправляем запрос к целевому серверу
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Ошибка отправки запроса", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем заголовки ответа
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Устанавливаем код состояния и тело ответа
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main1() {
	http.HandleFunc("/", handler)
	log.Println("Прокси-сервер запущен на порту 8003")
	err := http.ListenAndServe("185.102.139.161:8003", nil)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
