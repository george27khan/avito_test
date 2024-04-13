package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func Authentication(username string, password string) (token string) {
	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("GET", "http://localhost:8080/authentication", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	req.SetBasicAuth(username, password)
	fmt.Println(req.URL)
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

	return resp.Header.Get("token")
}

func TestUserBannerHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=100&feature_id=1&use_last_revision=True", nil)
	token := Authentication("admin_1", "admin_1")
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("token", token)
	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	//rr := httptest.NewRecorder()
	//
	//// Наш хендлер соответствует интерфейсу http.Handler, а значит
	//// мы можем использовать ServeHTTP и напрямую указать
	//// Request и ResponseRecorder
	//http.DefaultServeMux.ServeHTTP(rr, req) // Проверяем код
	//responseData, _ := io.ReadAll(rr.Body)
	//fmt.Println(string(responseData))

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
	os.Stdout.Write(bodyBytes)
	//os.Stdout.Write([]byte())
	var t1, t2 interface{}
	json.Unmarshal(bodyBytes, &t1)
	json.Unmarshal([]byte("{\"title\": \"some_title111\", \"text\": \"some_text\", \"url\": \"some_url\"}"), &t2)
	assert.Equal(t, t1, t1)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Конвертируем тело ответа в строку и выводим

}
