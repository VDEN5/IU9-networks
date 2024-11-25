package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Структура для хранения новостей
type NewsItem struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Link        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	PublishedAt time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// WebSocket обновление
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Получение данных из базы и отправка по WebSocket
func serveWs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Обновляем до WebSocket-соединения
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установке WebSocket-соединения:", err)
		return
	}
	defer conn.Close()

	for {
		// Получаем данные из таблицы
		var newsItems []NewsItem
		if err := db.Find(&newsItems).Error; err != nil {
			log.Println("Ошибка при получении данных:", err)
			return
		}

		// Конвертируем данные в JSON
		data, err := json.Marshal(newsItems)
		if err != nil {
			log.Println("Ошибка при кодировании данных:", err)
			return
		}

		// Отправляем данные по WebSocket
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("Ошибка при отправке сообщения:", err)
			return
		}

		// Обновление каждые 5 секунд
		time.Sleep(5 * time.Second)
	}
}

func (NewsItem) TableName() string {
	return "iu9vden" // Укажите здесь желаемое имя таблицы
}

func main() {
	// Подключаемся к базе данных
	dsn := "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Обработчик WebSocket-соединения
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(db, w, r)
	})

	// Статический сервер для HTML-дэшборда
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	r := gin.Default()

	// Обработчик GET-запроса для отображения формы
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Обработчик POST-запроса для получения данных из формы
	r.POST("/submit", func(c *gin.Context) {
		name := c.PostForm("name")
		c.HTML(http.StatusOK, "index.html", name)
	})

	// Установка шаблонов
	r.LoadHTMLGlob("templates/*")

	// Запуск сервера на порту 8080
	r.Run(":8080")
}
