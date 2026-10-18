package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	password := os.Getenv("TODO_PASSWORD")
	if Token == "" && password != "" {
		token, err := loginForTests(password)
		if err != nil {
			fmt.Fprintln(os.Stderr, "получение тестового JWT:", err)
			os.Exit(1)
		}
		Token = token
	}
	os.Exit(m.Run())
}

func loginForTests(password string) (string, error) {
	body, err := json.Marshal(map[string]string{"password": password})
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Post(getURL("api/signin"), "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("signin вернул HTTP %d", response.StatusCode)
	}
	var result struct {
		Token string `json:"token"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Token == "" {
		return "", fmt.Errorf("сервер не выдал JWT: %s", result.Error)
	}
	return result.Token, nil
}
