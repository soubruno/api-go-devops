package user

import (
	"errors"
	"testing"
)

type MockUserRepository struct {
	users []User
}

func (m *MockUserRepository) GetAll() ([]User, error) {
	return m.users, nil
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Create(u User) (*User, error) {
	u.ID = len(m.users) + 1
	m.users = append(m.users, u)
	return &u, nil
}

func (m *MockUserRepository) Update(id int, u User) (*User, error) {
	for i, user := range m.users {
		if user.ID == id {
			u.ID = id
			m.users[i] = u
			return &u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Delete(id int) error {
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

func TestUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)

	// Teste Create
	user, err := service.CreateUser(User{Name: "Teste", Email: "teste@example.com"})
	if err != nil {
		t.Fatalf("esperava sucesso ao criar usuario, erro: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("esperava ID 1, obteve %d", user.ID)
	}

	// Teste GetByID
	found, err := service.GetUserByID(1)
	if err != nil || found.Name != "Teste" {
		t.Fatalf("falha ao buscar usuario por ID")
	}

	// Teste GetAll
	all, err := service.GetAllUsers()
	if err != nil || len(all) != 1 {
		t.Fatalf("esperava 1 usuario retornado")
	}
}