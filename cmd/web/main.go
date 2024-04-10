package main

import (
	"avito_test/cmd/web/banner_handler"
	bn "avito_test/pkg/db_avito_banner/banner"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"time"
)

func server() {
	r := gin.Default()
	r.GET("/user_banner", banner_handler.PostBanner)
	r.GET("/banner", banner_handler.GetBanner)
	r.POST("/banner", banner_handler.PostBanner)
	r.PATCH("/banner/:id", banner_handler.PatchBanner)
	r.DELETE("/banner/:id", banner_handler.PostBanner)

	r.Run("localhost:8080") //63342
}

func postBanner() {
	for i := 1; i < 6; i++ {
		test := bn.Banner{TagIds: []int64{int64(1 * i), int64(100 * i)},
			FeatureId: int64(i),
			Content:   "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
			IsActive:  true}
		// Кодируем структуру User в JSON
		bytesRepresentation, err := json.Marshal(test)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := http.Post("http://localhost:8080/banner", "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close() // закрываем тело ответа после работы с ним

		data, err := io.ReadAll(resp.Body) // читаем ответ
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%s", data) // печатаем ответ как строку
	}
}

func getBanner() {
	//resp, err := http.Get("http://localhost:8080/banner?tag_id=2&limit=4&offset=2")
	resp, err := http.Get("http://localhost:8080/banner?feature_id=2&limit=4&offset=0")

	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close() // закрываем тело ответа после работы с ним

	data, err := io.ReadAll(resp.Body) // читаем ответ
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(data)) // печатаем ответ как строку
}

func patchBanner() {
	//test := bn.Banner{
	//	TagIds: []int64{int64(1), int64(100)},
	//	//FeatureId: nil,
	//	Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
	//	IsActive: false}
	////Кодируем структуру User в JSON
	//bytesRepresentation, err := json.Marshal(test)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	bytesRepresentation := []byte("{\n\"tag_ids\": [1,100,1000],\n\"feature_id\": 1,\n\"content\": \"{\\\"title\\\": \\\"some_title111\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\",\n\"is_active\": false\n}")

	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("PATCH", "http://localhost:8080/banner/1", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
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

	// Читаем и конвертируем тело ответа в байты
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// Конвертируем тело ответа в строку и выводим
	fmt.Printf("API ответ в форме строки: %s\n", bodyBytes)

	// Вывод статуса ответа (если 200 - то успешный)
	fmt.Println("Статус ответа:", resp.Status)
}

func main() {
	go server()
	time.Sleep(time.Second * 2)
	//postBanner()
	//getBanner()
	patchBanner()
}
