#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <winsock2.h>
#include <ws2tcpip.h>

#pragma comment(lib, "ws2_32.lib")

#define PORT 8080
#define BUFFER_SIZE 1024

int main() {
    WSADATA wsaData;
    SOCKET listenSocket, clientSocket;
    struct sockaddr_in serverAddr, clientAddr;
    int addrLen = sizeof(clientAddr);
    char buffer[BUFFER_SIZE];

    // Инициализация Winsock
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0) {
        printf("Ошибка инициализации Winsock: %d\n", WSAGetLastError());
        return 1;
    }

    // Создание сокета
    listenSocket = socket(AF_INET, SOCK_STREAM, 0);
    if (listenSocket == INVALID_SOCKET) {
        printf("Не удалось создать сокет: %d\n", WSAGetLastError());
        WSACleanup();
        return 1;
    }

    // Настройка структуры адреса сервера
    serverAddr.sin_family = AF_INET;
    serverAddr.sin_addr.s_addr = INADDR_ANY; // Принимаем соединения с любого IP
    serverAddr.sin_port = htons(PORT); // Порт в сетевом порядке байтов

    // Привязка сокета к адресу и порту
    if (bind(listenSocket, (struct sockaddr*)&serverAddr, sizeof(serverAddr)) == SOCKET_ERROR) {
        printf("Ошибка привязки сокета: %d\n", WSAGetLastError());
        closesocket(listenSocket);
        WSACleanup();
        return 1;
    }

    // Начинаем прослушивание входящих соединений
    if (listen(listenSocket, SOMAXCONN) == SOCKET_ERROR) {
        printf("Ошибка прослушивания: %d\n", WSAGetLastError());
        closesocket(listenSocket);
        WSACleanup();
        return 1;
    }

    printf("Ожидание входящих соединений на порту %d...\n", PORT);

    // Основной цикл серверной части
    while (1) {
        clientSocket = accept(listenSocket, (struct sockaddr*)&clientAddr, &addrLen);
        if (clientSocket == INVALID_SOCKET) {
            printf("Ошибка при принятии соединения: %d\n", WSAGetLastError());
            continue;
        }

        // Получаем данные от клиента
        int bytesReceived = recv(clientSocket, buffer, BUFFER_SIZE, 0);
        if (bytesReceived > 0) {
            buffer[bytesReceived] = '\0'; // Завершаем строку
            printf("Получено сообщение от клиента: %s\n", buffer);
            
            // Отправляем ответ клиенту
            send(clientSocket, "Сообщение получено", 18, 0);
        } else if (bytesReceived == 0) {
            printf("Соединение закрыто клиентом.\n");
        } else {
            printf("Ошибка при получении данных: %d\n", WSAGetLastError());
        }

        closesocket(clientSocket);
    }

    // Закрываем сокет прослушивания и очищаем Winsock
    closesocket(listenSocket);
    WSACleanup();
    return 0;
}
