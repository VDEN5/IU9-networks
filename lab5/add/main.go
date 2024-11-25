package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
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

// Функция для транслитерации русского текста
func transliterate(text string) string {
	translitMap := map[rune]string{
		'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "Yo", 'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y",
		'К': "K", 'Л': "L", 'М': "M", 'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U", 'Ф': "F",
		'Х': "Kh", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Shch", 'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i", 'й': "y",
		'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f",
		'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	}

	var translitText strings.Builder
	for _, ch := range text {
		if val, ok := translitMap[ch]; ok {
			translitText.WriteString(val)
		} else {
			translitText.WriteRune(ch)
		}
	}
	return translitText.String()
}

// Метод TableName задает собственное имя таблицы
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

	// Создаем таблицу (если еще не создана)
	err = db.AutoMigrate(&NewsItem{})
	if err != nil {
		log.Fatalf("Не удалось создать таблицу: %v", err)
	}

	// Подключаемся к RSS и парсим данные
	rssURL := "https://news.rambler.ru/rss/Guadeloupe/" // Замените на нужный URL RSS-канала
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		log.Fatalf("Не удалось получить RSS данные: %v", err)
	}

	for _, item := range feed.Items {
		// Заменяем русский текст на транслит
		transliteratedDescription := transliterate(item.Description)

		// Создаем запись
		news := NewsItem{
			Title:       transliterate(item.Title), // Транслитерация заголовка
			Link:        item.Link,
			Description: transliteratedDescription,
		}

		// Парсим дату публикации
		published, err := time.Parse(time.RFC1123Z, item.Published)
		if err == nil {
			news.PublishedAt = published
		} else {
			news.PublishedAt = time.Now()
		}

		// Сохраняем запись в базу данных
		err = db.Create(&news).Error
		if err != nil {
			log.Printf("Не удалось сохранить новость: %v", err)
		}
	}

	fmt.Println("Завершено обновление новостей.")
}
