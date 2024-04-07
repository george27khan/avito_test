package main

import (
	"avito_test/cmd/web/banner_handler"
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
	r.GET("/banner", banner_handler.PostBanner)
	r.POST("/banner", banner_handler.PostBanner)
	r.PATCH("/banner/:id", banner_handler.PostBanner)
	r.DELETE("/banner/:id", banner_handler.PostBanner)

	r.Run("localhost:8080") //63342
}

func postBanner() {
	type PostBody struct {
		TagIds    []int  `json:"tag_ids"`
		FeatureId int    `json:"feature_id"`
		Content   string `json:"content"`
		IsActive  bool   `json:"is_active"`
	}

	test := PostBody{TagIds: []int{0},
		FeatureId: 0,
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

func getBanner() {
	resp, err := http.Get("http://localhost:8080/banner?feature_id=0&tag_id=0")
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

func main() {
	go server()
	time.Sleep(time.Second * 5)
	postBanner()
	getBanner()
}
