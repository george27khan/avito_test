package banner_test

import (
	"encoding/json"
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
	testCases := map[string]map[string]interface{}{
		"success_active_banner": {"username": "user_1",
			"password":          "user_1",
			"tag_id":            "2",
			"feature_id":        "1",
			"use_last_revision": "true",
			"expected_body":     "\"{\\\"title\\\": \\\"title1\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusOK,
		},
		"success_active_banner_cash": {"username": "user_1",
			"password":          "user_1",
			"tag_id":            "2",
			"feature_id":        "1",
			"use_last_revision": "false",
			"expected_body":     "\"{\\\"title\\\": \\\"title1\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusOK,
		},
		"empty_use_last_revision": {"username": "admin_1",
			"password":          "admin_1",
			"tag_id":            "11",
			"feature_id":        "10",
			"use_last_revision": "",
			"expected_body":     "\"{\\\"title\\\": \\\"title10\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusOK,
		},
		"success_noactive_banner_admin": {"username": "admin_1",
			"password":          "admin_1",
			"tag_id":            "11",
			"feature_id":        "10",
			"use_last_revision": "true",
			"expected_body":     "\"{\\\"title\\\": \\\"title10\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusOK,
		},
		"success_noactive_banner": {"username": "user_1",
			"password":          "user_1",
			"tag_id":            "11",
			"feature_id":        "10",
			"use_last_revision": "true",
			"expected_body":     "\"{\\\"title\\\": \\\"title10\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusForbidden,
		},
		"bad_tag_id": {"username": "admin_1",
			"password":          "admin_1",
			"tag_id":            "",
			"feature_id":        "10",
			"use_last_revision": "true",
			"expected_body":     "\"{\\\"title\\\": \\\"title10\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusBadRequest,
		},
		"bad_feature_id": {"username": "admin_1",
			"password":          "admin_1",
			"tag_id":            "11",
			"feature_id":        "",
			"use_last_revision": "true",
			"expected_body":     "\"{\\\"title\\\": \\\"title10\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"",
			"expected_status":   http.StatusBadRequest,
		},
	}
	client := &http.Client{} // создаем http клиент
	for test, params := range testCases {
		var bannerResp, bannerExp interface{}

		req, err := http.NewRequest("GET", "http://localhost:8080/user_banner", nil)
		q := req.URL.Query()

		q.Add("tag_id", params["tag_id"].(string))
		q.Add("feature_id", params["feature_id"].(string))
		q.Add("use_last_revision", params["use_last_revision"].(string))
		req.URL.RawQuery = q.Encode()
		os.Stdout.Write([]byte(req.URL.String()))

		token := Authentication(params["username"].(string), params["password"].(string))
		if err != nil {
			log.Println("Ошибка при создании запроса:", err)
			return
		}
		// Устанавливаем заголовки запроса
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("token", token)
		// Отправляем запрос
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

		if test != "success_noactive_banner" && test != "bad_tag_id" && test != "bad_feature_id" {
			json.Unmarshal(bodyBytes, &bannerResp)
			json.Unmarshal([]byte(params["expected_body"].(string)), &bannerExp)
			assert.Equal(t, bannerResp, bannerExp)
		}
		assert.Equal(t, params["expected_status"].(int), resp.StatusCode)
	}
}
