package user

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	repo := &PostgresUserRepository{db: db}
	repo.initTable()
	return repo
}

func (r *PostgresUserRepository) initTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL
	);`
	_, err := r.db.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao inicializar tabela users: %v", err)
	}
}

func (r *PostgresUserRepository) GetAll() ([]User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresUserRepository) GetByID(id int) (*User, error) {
	var u User
	err := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Create(u User) (*User, error) {
	err := r.db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) Update(id int, u User) (*User, error) {
	result, err := r.db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	u.ID = id
	return &u, nil
}

func (r *PostgresUserRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}