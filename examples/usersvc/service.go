package usersvc

import (
	"context"
	"errors"
	"sync"
	"time"
)

//go:generate protogo httptest ./service.go Service --mode=server --spec=./httpserver.httptest.yaml --out=./httpserver_test.go
//go:generate protogo httptest ./service.go Service --mode=client --spec=./httpclient.httptest.yaml --out=./httpclient_test.go

type Service interface {
	GetUser(ctx context.Context, name string) (user *User, err error)
	ListUsers(ctx context.Context) (users []*User, err error)
	CreateUser(ctx context.Context, user *User) (err error)
	UpdateUser(ctx context.Context, name string, user *User) (err error)
	DeleteUser(ctx context.Context, name string) (err error)
}

type User struct {
	Name  string    `json:"name,omitempty" httptest:"name,omitempty"`
	Sex   string    `json:"sex,omitempty" httptest:"sex,omitempty"`
	Birth time.Time `json:"birth,omitempty" httptest:"birth,omitempty"`
}

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

type InmemService struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewInmemService() *InmemService {
	return &InmemService{
		users: map[string]*User{},
	}
}

func (s *InmemService) GetUser(ctx context.Context, name string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[name]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (s *InmemService) ListUsers(ctx context.Context) ([]*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*User
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *InmemService) CreateUser(ctx context.Context, user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[user.Name]; ok {
		return ErrAlreadyExists
	}
	s.users[user.Name] = user
	return nil
}

func (s *InmemService) UpdateUser(ctx context.Context, name string, user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.Name = name
	s.users[name] = user
	return nil
}

func (s *InmemService) DeleteUser(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[name]; !ok {
		return ErrNotFound
	}
	delete(s.users, name)
	return nil
}
