#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <winsock2.h>
#include <ws2tcpip.h>

#pragma comment(lib, "ws2_32.lib")

#define SERVER_IP "127.0.0.1" // IP-адрес сервера (локальный в данном случае)
#define PORT 8080              // Порт сервера
#define BUFFER_SIZE 1024

int main() {
    WSADATA wsaData;
    SOCKET sock;
    struct sockaddr_in serverAddr;
    char buffer[BUFFER_SIZE];

    // Инициализация Winsock
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        printf("Ошибка инициализации Winsock: %d\n", WSAGetLastError());
        return 1;
    }

    // Создание сокета
    sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock == INVALID_SOCKET) {
        printf("Не удалось создать сокет: %d\n", WSAGetLastError());
        WSACleanup();
        return 1;
    }

    // Настройка структуры адреса сервера
    serverAddr.sin_family = AF_INET;
    serverAddr.sin_port = htons(PORT); // Порт в сетевом порядке байтов
    inet_pton(AF_INET, SERVER_IP, &serverAddr.sin_addr); // Преобразование IP-адреса

    // Подключение к серверу
    if (connect(sock, (struct sockaddr*)&serverAddr, sizeof(serverAddr)) == SOCKET_ERROR) {
        printf("Ошибка подключения к серверу: %d\n", WSAGetLastError());
        closesocket(sock);
        WSACleanup();
        return 1;
    }

    // Отправка сообщения на сервер
    const char *message = "hlhlkjbjhlvhjhvu";
    send(sock, message, strlen(message), 0);

    // Получение ответа от сервера
    int bytesReceived = recv(sock, buffer, BUFFER_SIZE, 0);
    if (bytesReceived > 0) {
        buffer[bytesReceived] = '\0'; // Завершаем строку
        printf("Ответ от сервера: %s\n", buffer);
    } else if (bytesReceived == 0) {
        printf("Сервер закрыл соединение.\n");
    } else {
        printf("Ошибка при получении данных: %d\n", WSAGetLastError());
    }

    // Закрытие сокета и очистка Winsock
    closesocket(sock);
    WSACleanup();
    return 0;
}
