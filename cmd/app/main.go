package main

import (
	"avito_test/internal/app/handlers"
	bn "avito_test/pkg/postgres_db/banner"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

func server() {
	r := gin.Default()

	//отключение middleware для аутентификации
	r.GET("/authentication", handlers.GetToken)

	r.Use(handlers.AuthMiddleware())
	r.GET("/user_banner", handlers.GetUserBanner)
	r.GET("/banner_version", handlers.GetBannerVersion)
	r.GET("/banner", handlers.GetBanner)
	r.POST("/banner", handlers.PostBanner)
	r.PATCH("/banner/:id", handlers.PatchBanner)
	r.DELETE("/banner/:id", handlers.DeleteBanner)

	r.Run("0.0.0.0:8080") //63342
}

func postBanner() {
	token := Authentication("admin_1", "admin_1")

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
		req, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return
		}
		// Устанавливаем заголовки запроса
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("token", token)

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
}

func getBanner() {
	token := Authentication("admin_1", "admin_1")
	req, err := http.NewRequest("GET", "http://localhost:8080/banner?tag_id=1&limit=4&offset=0", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("token", token)

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

func getBannerVersion() {
	//resp, err := http.Get("http://localhost:8080/banner?tag_id=2&limit=4&offset=2")
	resp, err := http.Get("http://localhost:8080/banner/version?feature_id=2&limit=4&offset=0")

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

func getUserBanner() {
	token := Authentication("admin_1", "admin_1")

	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner?tag_id=2&feature_id=1&use_last_revision=False", nil)
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("token", token)

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

func patchBanner() {
	token := Authentication("admin_1", "admin_1")
	bytesRepresentation := []byte("{\n\"tag_ids\": [1,100,1000],\n\"feature_id\": 1,\n\"content\": \"{\\\"title\\\": \\\"some_title111\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\",\n\"is_active\": false\n}")

	// Создаем новый HTTP-запрос с методом POST
	req, err := http.NewRequest("PATCH", "http://localhost:8080/banner/1", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Println("Ошибка при создании запроса:", err)
		return
	}
	// Устанавливаем заголовки запроса
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("token", token)

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

func delBanner() {
	token := Authentication("admin_1", "admin_1")
	// Создаем новый HTTP-запрос с методом POST
	for i := 1; i < 6; i++ {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/banner/%v", i), nil)

		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return
		}
		// Устанавливаем заголовки запроса
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("token", token)

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

func main() {
	load_env()
	server()
	//time.Sleep(time.Second * 2)
	//postBanner()//++
	//getBanner()//+
	//patchBanner()
	//delBanner() //+
	//getUserBanner()
}
