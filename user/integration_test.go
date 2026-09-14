package user

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestIntegrationPostgres(t *testing.T) {
	connStr := os.Getenv("TEST_DB_URL")
	if connStr == "" {
		t.Skip("Pulando teste de integracao: TEST_DB_URL nao configurada")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("falha ao abrir conexao: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("falha ao pingar banco de integracao: %v", err)
	}

	repo := NewUserRepository(db)
	u, err := repo.Create(User{Name: "Integracao", Email: "integracao@test.com"})
	if err != nil {
		t.Fatalf("falha ao inserir usuario: %v", err)
	}

	found, err := repo.GetByID(u.ID)
	if err != nil || found.Email != "integracao@test.com" {
		t.Fatalf("falha ao recuperar usuario no banco real")
	}

	_ = repo.Delete(u.ID)
}