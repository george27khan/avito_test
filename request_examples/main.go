package main

import (
	bn "avito_test/pkg/postgres_db/banner"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func endRequest(req *http.Request) {
	token := Authentication("admin_1", "admin_1")
	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("token", token)
	fmt.Println(req.Method, req.URL.String())
	fmt.Println("Headers:", req.Header)
	fmt.Println("Body:", req.Body)

	// Отправляем запрос
	client := &http.Client{}    // создаем http клиент
	resp, err := client.Do(req) // передаем выше подготовленный запрос на отправку
	if err != nil {
		log.Println("Ошибка при выполнении запроса: ", err)
		return
	}
	defer resp.Body.Close()

	// Читаем и конвертируем тело ответа в байты
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// Конвертируем тело ответа в строку и выводим
	fmt.Printf("Тело ответа: %s\n", bodyBytes)
	// Вывод статуса ответа (если 200 - то успешный)
	fmt.Println("Статус ответа:", resp.Status)
}

func postBanner() {
	test := bn.Banner{TagIds: []int64{10000, 10001},
		FeatureId: 10000,
		Content:   "{\"title\": \"some_title10000\", \"text\": \"some_text\", \"url\": \"some_url\"}",
		IsActive:  true}
	bytesRepresentation, err := json.Marshal(test)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	endRequest(req)
}

func getBanner() {
	req, err := http.NewRequest("GET", "http://localhost:8080/banner?feature_id=10000&limit=4&offset=0", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	endRequest(req)
}

func getUserBanner() {
	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=10000&feature_id=10000&use_last_revision=False", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	endRequest(req)
}

func patchBanner() {
	bytesRepresentation := []byte("{\n\"tag_ids\": [1,2,10003],\n\"feature_id\": 1,\n\"content\": \"{\\\"title\\\": \\\"some_title111111\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\",\n\"is_active\": false\n}")
	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("PATCH", "http://localhost:8080/banner/1", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	endRequest(req)
}

func delBanner() {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/banner/3", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	endRequest(req)
}

func Authentication(username string, password string) (token string) {
	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("GET", "http://localhost:8080/authentication", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	req.SetBasicAuth(username, password)
	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	// Отправляем запрос
	client := &http.Client{}    // создаем http клиент
	resp, err := client.Do(req) // передаем выше подготовленный запрос на отправку
	if err != nil {
		log.Println("Ошибка при выполнении запроса: ", err)
		return
	}
	defer resp.Body.Close() // не забываем закрыть тело
	// Вывод статуса ответа (если 200 - то успешный)
	return resp.Header.Get("token")
}

func main() {
	fmt.Println("postBanner")
	postBanner()

	fmt.Println("\ngetUserBanner")
	getUserBanner()

	fmt.Println("\ngetBanner")
	getBanner()

	fmt.Println("\npatchBanner")
	patchBanner()

	fmt.Println("\ndelBanner")
	delBanner()

}
