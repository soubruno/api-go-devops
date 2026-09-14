package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"api-go/user"

	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbPort == "" {
		dbPort = "5432"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var db *sql.DB
	var err error

	// Loop de retry para aguardar o container do PostgreSQL estar pronto para conexões
	for i := 0; i < 15; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Conexao com PostgreSQL realizada com sucesso!")
				break
			}
		}
		log.Printf("Aguardando banco de dados responder... tentativa %d/15\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Nao foi possivel conectar ao PostgreSQL: %v", err)
	}
	defer db.Close()

	repo := user.NewUserRepository(db)
	service := user.NewUserService(repo)
	controller := user.NewUserController(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListUsers(w, r)
		case http.MethodPost:
			controller.CreateUser(w, r)
		default:
			http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUser(w, r)
		case http.MethodPut:
			controller.UpdateUser(w, r)
		case http.MethodDelete:
			controller.DeleteUser(w, r)
		default:
			http.Error(w, "Metodo nao permitido", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}