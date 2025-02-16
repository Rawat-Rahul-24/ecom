package user

import (
	"bytes"
	"ecom/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is empty", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "user",
			Password: "123",
			Email: "asdcom",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rresponseWriter := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rresponseWriter, req)

		if rresponseWriter.Code != http.StatusBadRequest {
			t.Errorf("expected error code %d, got %d", http.StatusBadRequest, rresponseWriter.Code)
		}
	})

	t.Run("should create user succesfully", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "user",
			LastName: "user",
			Password: "123",
			Email: "asd@gmail.com",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rresponseWriter := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rresponseWriter, req)

		if rresponseWriter.Code != http.StatusCreated {
			t.Errorf("expected error code %d, got %d", http.StatusCreated, rresponseWriter.Code)
		}
	})
}

type mockUserStore struct {}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m * mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}