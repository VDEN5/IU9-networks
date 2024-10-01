package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
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
	r.Run("185.102.139.168:8080")
}
