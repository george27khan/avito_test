package main

import (
	bn "avito_test/pkg/postgres_db/banner"
	db "avito_test/pkg/postgres_db/connection"
	"avito_test/pkg/postgres_db/user"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
)

func load_env() {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
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
	fmt.Println("Статус ответа:", resp.Status)
	return resp.Header.Get("token")
}

// createUsers создание тестовых пользователей
func createUsers() error {
	var u user.User
	ctx := context.Background()
	conn, tx, err := db.ConnectPoolTrx(ctx)
	defer conn.Release()
	defer func() {
		if err != nil {
			tx.Rollback(ctx)

		} else {
			tx.Commit(ctx)
		}
	}()
	if err != nil {
		return fmt.Errorf("createUsers error - %v", err.Error())
	}
	for i := 1; i <= 3; i++ {
		u = user.User{UserName: fmt.Sprintf("admin_%v", i),
			Password: fmt.Sprintf("admin_%v", i),
			IsAdmin:  true}
		if err = u.Insert(ctx, conn); err != nil {
			return fmt.Errorf("createUsers error - %v", err.Error())
		}
	}
	for i := 1; i <= 5; i++ {
		u = user.User{UserName: fmt.Sprintf("user_%v", i),
			Password: fmt.Sprintf("user_%v", i),
			IsAdmin:  false}
		if err = u.Insert(ctx, conn); err != nil {
			return fmt.Errorf("createUsers error - %v", err.Error())
		}
	}
	return nil
}

// createBanners создание тестовых баннеров
func createBanners() error {
	var isActive bool
	token := Authentication("admin_1", "admin_1")

	for i := 1; i < 1000; i++ {
		if i%10 == 0 {
			isActive = false
		} else {
			isActive = true
		}
		test := bn.Banner{TagIds: []int64{int64(1 + i), int64(2 + i)},
			FeatureId: int64(i),
			Content:   fmt.Sprintf("{\"title\": \"title%v\", \"text\": \"some_text\", \"url\": \"some_url\"}", i),
			IsActive:  isActive}
		// Кодируем структуру User в JSON
		bytesRepresentation, err := json.Marshal(test)
		if err != nil {
			log.Fatalln(err)
		}
		req, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return nil
		}
		// Устанавливаем заголовки запроса
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("token", token)

		// Отправляем запрос
		client := &http.Client{}    // создаем http клиент
		resp, err := client.Do(req) // передаем выше подготовленный запрос на отправку
		if err != nil {
			log.Println("Ошибка при выполнении запроса: ", err)
			return nil
		}
		defer resp.Body.Close() // не забываем закрыть тело

		// Читаем и конвертируем тело ответа в байты
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil
		}

		// Конвертируем тело ответа в строку и выводим
		fmt.Printf("API ответ в форме строки: %s\n", bodyBytes)

		// Вывод статуса ответа (если 200 - то успешный)
		fmt.Println("Статус ответа:", resp.Status)
	}
	return nil
}

func main() {
	load_env()
	err := createUsers()
	if err != nil {
		fmt.Println("err ", err.Error())
	}
	err = createBanners()
	if err != nil {
		fmt.Println("err ", err.Error())
	}
}
