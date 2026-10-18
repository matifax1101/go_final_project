package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/matifax1101/go_final_project/pkg/api"
	"github.com/matifax1101/go_final_project/pkg/db"
)

const defaultPort = "7540"
const defaultDBFile = "scheduler.db"

func main() {
	if err := run(); err != nil {
		// run() завершился, следовательно, его defer уже выполнены.
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}
	password := os.Getenv("TODO_PASSWORD")
	secret := os.Getenv("TODO_JWT_SECRET")

	signingKey, err := jwtSigningKey(password, secret)
	if err != nil {
		return fmt.Errorf("JWT configuration: %w", err)
	}
	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("database initialization: %w", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	api.Init(mux, password, signingKey)
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Server started at http://localhost:%s", port)
	return server.ListenAndServe()
}

func jwtSigningKey(password, secret string) ([]byte, error) {
	if password == "" {
		return nil, nil
	}

	var rawKey []byte
	if secret != "" {
		if len(secret) < 32 {
			return nil, fmt.Errorf("TODO_JWT_SECRET должен содержать не менее 32 байт")
		}
		rawKey = []byte(secret)
	} else {
		rawKey = make([]byte, 32)
		if _, err := rand.Read(rawKey); err != nil {
			return nil, fmt.Errorf("generate random JWT key: %w", err)
		}
	}

	// Привязываем подпись к паролю: его смена отзывает все старые JWT.
	mac := hmac.New(sha256.New, rawKey)
	if _, err := mac.Write([]byte("go_final_project/jwt/v1\x00" + password)); err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}
