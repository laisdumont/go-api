package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"golang.org/x/crypto/bcrypt"

	"go-api/handler"
	"go-api/model"
	"go-api/repository"
	"go-api/router"
	"go-api/service"
)

var (
	testHandler *handler.UserHandler
	authToken   string
	teardown    func()
)

func setupTestDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.0"),
		mysql.WithDatabase("testdb"),
		mysql.WithUsername("testuser"),
		mysql.WithPassword("senha123"),
	)
	if err != nil {
		t.Fatalf("Erro ao iniciar container do MySQL: %v", err)
	}

	endpoint, err := container.ConnectionString(ctx, "mysql")
	if err != nil {
		t.Fatalf("Erro ao obter string de conexão: %v", err)
	}

	db, err := sql.Open("mysql", endpoint)
	if err != nil {
		t.Fatalf("Erro ao conectar no MySQL: %v", err)
	}

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		t.Fatalf("Banco demorou demais pra responder: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("Erro ao criar tabela: %v", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Erro ao gerar hash da senha: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users (name, password) VALUES (?, ?)`, "testuser", string(hashed))
	if err != nil {
		t.Fatalf("Erro ao inserir usuário de teste: %v", err)
	}

	teardown = func() {
		db.Close()
		if err := container.Terminate(ctx); err != nil {
			log.Printf("Erro ao encerrar container: %v", err)
		}
	}

	return db
}

func TestMain(m *testing.M) {
	db := setupTestDB(&testing.T{})

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	testHandler = handler.NewUserHandler(svc)

	token, err := loginForTests()
	if err != nil {
		panic("🔥 Falha no login do teste: " + err.Error())
	}
	authToken = token

	code := m.Run()
	teardown()
	os.Exit(code)
}

func loginForTests() (string, error) {
	payload := model.User{
		Name:     "testuser",
		Password: "senha123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := executeRequest(req)
	if resp.Code != http.StatusOK {
		return "", fmt.Errorf("esperava 200, recebeu %d", resp.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res["token"], nil
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

	req := authRequest("PUT", "/users/1", body)
	resp := executeRequest(req)

	if resp.Code != http.StatusOK {
		t.Errorf("Erro ao atualizar: %d", resp.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	req := authRequest("DELETE", "/users/1", nil)
	resp := executeRequest(req)

	if resp.Code != http.StatusOK {
		t.Errorf("Erro ao deletar: %d", resp.Code)
	}
}
