package test

import (
	"bytes"
	"encoding/json"

	"go-api/db"
	"go-api/handler"
	"go-api/model"
	"go-api/repository"
	"go-api/router"
	"go-api/service"

	"net/http"
	"net/http/httptest"
	"testing"
)

var authToken string

func authenticate() {
	payload := model.User{
		Name:     "testuser",
		Password: "senha123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	if resp.Code != http.StatusOK {
		panic("🔥 Falha no login do teste")
	}

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)
	authToken = res["token"]
}

var testHandler *handler.UserHandler

func init() {
	db.Connect()
	repo := repository.NewUserRepository(db.DB)
	svc := service.NewUserService(repo)
	testHandler = handler.NewUserHandler(svc)
	authenticate()
}

func authRequest(method, url string, body []byte) *http.Request {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	return req
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r := router.SetupRoutes(testHandler)
	r.ServeHTTP(rr, req)
	return rr
}

func TestRegister(t *testing.T) {
	payload := model.User{
		Name:     "testuser",
		Password: "senha123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := executeRequest(req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Esperava 201, recebeu %d", resp.Code)
	}
}

func TestLogin(t *testing.T) {
	payload := model.User{
		Name:     "testuser",
		Password: "senha123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := executeRequest(req)

	if resp.Code != http.StatusOK {
		t.Errorf("Login falhou: código %d", resp.Code)
	}
}

func TestGetUsers(t *testing.T) {
	req := authRequest("GET", "/users", nil)
	resp := executeRequest(req)

	if resp.Code != http.StatusOK {
		t.Errorf("Erro ao buscar usuários: código %d", resp.Code)
	}

	var users []model.User
	if err := json.Unmarshal(resp.Body.Bytes(), &users); err != nil {
		t.Errorf("Erro ao parsear usuários: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	payload := model.User{
		Name: "Nome Atualizado",
	}
	body, _ := json.Marshal(payload)

	req := authRequest("PUT", "/users/6", body)
	resp := executeRequest(req)

	if resp.Code != http.StatusOK {
		t.Errorf("Erro ao atualizar: %d", resp.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	req := authRequest("DELETE", "/users/7", nil)
	resp := executeRequest(req)

	if resp.Code != http.StatusOK {
		t.Errorf("Erro ao deletar: %d", resp.Code)
	}
}
