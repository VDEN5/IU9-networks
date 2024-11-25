package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	ftpHost = "students.yss.su"
	ftpUser = "ftpiu8"
	ftpPass = "3Ru7yOTA"
)

func connectToFTP() (*ftp.ServerConn, error) {
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

func uploadFile(conn *ftp.ServerConn, localPath, remotePath string) error {
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

func downloadFile(conn *ftp.ServerConn, remotePath, localPath string) error {
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

func createDirectory(conn *ftp.ServerConn, dirPath string) error {
	err := conn.MakeDir(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}
	fmt.Println("Директория создана:", dirPath)
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

func listDirectory(conn *ftp.ServerConn, dirPath string) error {
	entries, err := conn.List(dirPath)
	if err != nil {
		return fmt.Errorf("ошибка получения содержимого директории: %w", err)
	}
	fmt.Println("Содержимое директории", dirPath)
	for _, entry := range entries {
		fmt.Println(" -", entry.Name)
	}
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

func main() {
	conn, err := connectToFTP()
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer conn.Logout()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("FTP-клиент запущен. Введите команду для выполнения:")
	fmt.Println("Доступные команды: upload, download, mkdir, rmdir, rmdir_recursive, ls, cd, rm")
	fmt.Println("Для выхода введите 'exit'.")

	for {
		fmt.Print("> ")
		scanner.Scan()
		commandLine := scanner.Text()
		fmt.Println(commandLine)
		parts := strings.Fields(commandLine)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		if command == "exit" {
			fmt.Println("Завершаем работу.")
			break
		}

		switch command {
		case "upload":
			if len(args) < 2 {
				fmt.Println("Использование: upload <localPath> <remotePath>")
				continue
			}
			err := uploadFile(conn, args[0], args[1])
			if err != nil {
				log.Println(err)
			}

		case "download":
			if len(args) < 2 {
				fmt.Println("Использование: download <remotePath> <localPath>")
				continue
			}
			err := downloadFile(conn, args[0], args[1])
			if err != nil {
				log.Println(err)
			}

		case "mkdir":
			if len(args) < 1 {
				fmt.Println("Использование: mkdir <dirPath>")
				continue
			}
			err := createDirectory(conn, args[0])
			if err != nil {
				log.Println(err)
			}

		case "rmdir":
			if len(args) < 1 {
				fmt.Println("Использование: rmdir <dirPath>")
				continue
			}
			err := deleteEmptyDirectory(conn, args[0])
			if err != nil {
				log.Println(err)
			}

		case "rmdir_recursive":
			if len(args) < 1 {
				fmt.Println("Использование: rmdir_recursive <dirPath>")
				continue
			}
			err := removeDirectoryRecursive(conn, args[0])
			if err != nil {
				log.Println(err)
			}

		case "ls":
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}
			err := listDirectory(conn, dir)
			if err != nil {
				log.Println(err)
			}

		case "cd":
			if len(args) < 1 {
				fmt.Println("Использование: cd <dirPath>")
				continue
			}
			err := changeDirectory(conn, args[0])
			if err != nil {
				log.Println(err)
			}

		case "rm":
			if len(args) < 1 {
				fmt.Println("Использование: delete <filePath>")
				continue
			}
			err := deleteFile(conn, args[0])
			if err != nil {
				log.Println(err)
			}

		default:
			fmt.Println("Неизвестная команда:", command)
		}
	}
}
