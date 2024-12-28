package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jlaffaye/ftp"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var messages []string
var mu sync.Mutex
var a1, a2, a3 string

func connectToFTP(ftpHost, ftpUser, ftpPass string) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(ftpHost+":21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	err = conn.Login(ftpUser, ftpPass)
	if err != nil {
		return nil, err
	}
	fmt.Println("Подключение успешно!")
	return conn, nil
}
func createDirectory(conn *ftp.ServerConn, dirPath string) error {
	err := conn.MakeDir(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}
	fmt.Println("Директория создана:", dirPath)
	return nil
}
func uploadFile1(conn *ftp.ServerConn, localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	fmt.Println("Файл успешно загружен:", remotePath)
	return nil
}
func uploadFile2(conn *ftp.ServerConn, remotePath string, file *os.File) error {
	err := conn.Stor(remotePath, file)
	if err != nil {
		return fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	fmt.Println("Файл успешно загружен:", remotePath)
	return nil
}
func downloadFile1(conn *ftp.ServerConn, remotePath, localPath string) error {
	resp, err := conn.Retr(remotePath)
	if err != nil {
		return fmt.Errorf("ошибка скачивания файла: %w", err)
	}
	defer resp.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = localFile.ReadFrom(resp)
	if err != nil {
		return err
	}
	fmt.Println("Файл успешно скачан:", localPath)
	return nil
}
func deleteFile(conn *ftp.ServerConn, filePath string) error {
	err := conn.Delete(filePath)
	if err != nil {
		return fmt.Errorf("ошибка удаления файла: %w", err)
	}
	fmt.Println("Файл удален:", filePath)
	return nil
}
func changeDirectory(conn *ftp.ServerConn, dirPath string) error {
	err := conn.ChangeDir(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка перехода в директорию: %w", err)
	}
	fmt.Println("Перешли в директорию:", dirPath)
	return nil
}

func deleteEmptyDirectory(conn *ftp.ServerConn, dirPath string) error {
	err := conn.RemoveDir(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка удаления пустой директории: %w", err)
	}
	fmt.Println("Пустая директория удалена:", dirPath)
	return nil
}

func removeDirectoryRecursive(c *ftp.ServerConn, dirName string) error {
	entries, err := c.List(dirName)
	if err != nil {
		return fmt.Errorf("ошибка получения содержимого директории %s: %w", dirName, err)
	}

	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		path := dirName + "/" + entry.Name
		if entry.Type == ftp.EntryTypeFolder {
			if err := removeDirectoryRecursive(c, path); err != nil {
				return err
			}
		} else {
			if err := c.Delete(path); err != nil {
				return fmt.Errorf("ошибка удаления файла %s: %w", path, err)
			}
		}
	}

	err = c.RemoveDir(dirName)
	if err != nil {
		return fmt.Errorf("ошибка удаления директории %s: %w", dirName, err)
	}

	fmt.Println("Директория рекурсивно удалена:", dirName)
	return nil
}
func listDirectory(conn *ftp.ServerConn, dirPath string) (error, string) {
	entries, err := conn.List(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка получения содержимого директории: %w", err), ""
	}
	res := ""
	res += "Содержимое директории" + dirPath
	fmt.Println("Содержимое директории", dirPath)
	for _, entry := range entries {
		qw := entry.Name
		if strings.Count(qw, ".") == 0 {
			qw = "dir:   " + qw
		} else {
			qw = "file:  " + qw
		}
		fmt.Println(" -", qw)
		res += " -" + qw
	}
	return nil, res
}

type FormData struct {
	Name    string
	Email   string
	Message string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Process form data
		r.ParseForm()
		data := FormData{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("name1"),
			Message: r.FormValue("name2"),
		}
		//fmt.Println(data)
		a1, a2, a3 = data.Name, data.Email, data.Message
		fmt.Println(a1, a2, a3)
		// Redirect to /home
		http.Redirect(w, r, "/blabla", http.StatusSeeOther) // Use StatusSeeOther for POST-redirect-GET pattern
		return                                              // Important: exit after redirect

	} else {
		// Display form
		tmpl, err := template.ParseFiles("templates/form.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
func main() {

	//scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("FTP-клиент запущен. Введите команду для выполнения:")
	fmt.Println("Доступные команды: upload, download, mkdir, rmdir, rmdir_recursive, ls, cd, rm")
	fmt.Println("Для выхода введите 'exit'.")
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/blabla", serveHome)
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		fmt.Printf("Ошибка при создании директории: %v\n", err)
		return
	}

	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)

	fmt.Println("Сервер запущен на http://localhost:8086")
	if err := http.ListenAndServe("185.104.251.226:8086", nil); err != nil {
		panic("Ошибка запуска сервера: " + err.Error())
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println(a1, a2, a3)
	conn1, err := connectToFTP(a1, a2, a3)
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn1.Logout()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()
	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("file")
		ew := r.FormValue("text")
		fmt.Println(ew)
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Создаем уникальное имя файла, чтобы избежать перезаписи
		filePath := fmt.Sprintf("uploads/%s", header.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
			return
		}
		uploadFile2(conn1, ew, out)

		lastUploadedFilePath = filePath // Обновляем путь к последнему загруженному файлу
		fmt.Fprintln(w, "Файл успешно загружен!")
		return
	}
	for {
		var msg string
		_, ms, err := conn.ReadMessage()
		msg = string(ms)
		parts := strings.Fields(msg)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]
		fmt.Println(msg + "cjknkj")
		if command == "ls" {
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			err, msg1 := listDirectory(conn1, dir)
			messages = append(messages, msg1)
			if err != nil {
				log.Println(err)
			}
		}
		if command == "mkdir" {
			if len(args) < 1 {
				fmt.Println("Использование: mkdir <dirPath>")
				continue
			}
			err := createDirectory(conn1, args[0])
			if err != nil {
				log.Println(err)
			}
		}
		if command == "rmdir" {
			if len(args) < 1 {
				fmt.Println("Использование: rmdir <dirPath>")
				continue
			}
			err := deleteEmptyDirectory(conn1, args[0])
			if err != nil {
				log.Println(err)
			}
		}
		if command == "rmdir_recursive" {
			if len(args) < 1 {
				fmt.Println("Использование: rmdir_recursive <dirPath>")
				continue
			}
			err := removeDirectoryRecursive(conn1, args[0])
			if err != nil {
				log.Println(err)
			}
		}
		if command == "cd" {
			if len(args) < 1 {
				fmt.Println("Использование: cd <dirPath>")
				continue
			}
			err := changeDirectory(conn1, args[0])
			if err != nil {
				log.Println(err)
			}
		}
		if command == "rm" {
			if len(args) < 1 {
				fmt.Println("Использование: delete <filePath>")
				continue
			}
			err := deleteFile(conn1, args[0])
			if err != nil {
				log.Println(err)
			}
		}
		if err != nil {
			fmt.Println("Ошибка чтения сообщения:", err)
			break
		}

		mu.Lock()
		messages = append(messages, msg)
		mu.Unlock()

		err = conn.WriteJSON(messages)
		if err != nil {
			fmt.Println("Ошибка отправки сообщения:", err)
			break
		}
	}
}

var lastUploadedFilePath string

func uploadFile(w http.ResponseWriter, r *http.Request) {
	conn1, err := connectToFTP(a1, a2, a3)
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn1.Logout()
	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("file")
		ew := r.FormValue("text")
		fmt.Println(ew)
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Создаем уникальное имя файла, чтобы избежать перезаписи
		filePath := fmt.Sprintf("uploads/%s", header.Filename)
		out, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
			return
		}
		uploadFile2(conn1, ew, out)

		lastUploadedFilePath = filePath // Обновляем путь к последнему загруженному файлу
		fmt.Fprintln(w, "Файл успешно загружен!")
		return
	}
	http.ServeFile(w, r, "upload.html")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	if lastUploadedFilePath == "" {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	file, err := os.Open(lastUploadedFilePath)
	if err != nil {
		http.Error(w, "Ошибка при открытии файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Устанавливаем заголовки для скачивания
	sd, _ := file.Stat()
	http.ServeContent(w, r, "", sd.ModTime(), file)

}
