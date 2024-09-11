package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/assylzhan-a/company-task/config"
	"github.com/assylzhan-a/company-task/internal/db"
	"github.com/assylzhan-a/company-task/internal/db/repository"
	handler "github.com/assylzhan-a/company-task/internal/delivery/http"
	"github.com/assylzhan-a/company-task/internal/domain/entity"
	uc "github.com/assylzhan-a/company-task/internal/domain/usecase"
	"github.com/assylzhan-a/company-task/internal/worker"
	"github.com/assylzhan-a/company-task/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testRouter *chi.Mux
	testDB     *pgxpool.Pool
)

func TestMain(m *testing.M) {
	cfg := config.Load()
	cfg.DatabaseURL = "postgres://user:password@localhost:5432/company_db?sslmode=disable"

	log := logger.NewLogger(cfg.LogLevel)

	var err error
	testDB, err = db.NewPostgresConnection(cfg.DatabaseURL, log)
	if err != nil {
		log.Error("Failed to connect to test database", "error", err)
		os.Exit(1)
	}

	// Run migrations
	if err := db.RunMigrations("../../migrations/", cfg.DatabaseURL); err != nil {
		log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Set up repositories and use cases
	userRepo := repository.NewUserRepository(testDB)
	companyRepo := repository.NewCompanyRepository(testDB)

	userUseCase := uc.NewUserUseCase(userRepo)
	companyUseCase := uc.NewCompanyUseCase(companyRepo, log)

	// Set up router
	testRouter = chi.NewRouter()
	handler.NewUserHandler(testRouter, userUseCase)
	handler.NewCompanyHandler(testRouter, companyUseCase)

	// Run tests
	code := m.Run()

	// Clean up database after tests
	cleanupDatabase(*log)

	// Close DB connection
	testDB.Close()

	os.Exit(code)
}

func TestUserRegistrationAndLogin(t *testing.T) {
	registerPayload := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	registerBody, _ := json.Marshal(registerPayload)
	registerReq := httptest.NewRequest("POST", "/v1/users/register", bytes.NewBuffer(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	registerRec := httptest.NewRecorder()
	testRouter.ServeHTTP(registerRec, registerReq)

	assert.Equal(t, http.StatusCreated, registerRec.Code)

	loginPayload := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginReq := httptest.NewRequest("POST", "/v1/users/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	testRouter.ServeHTTP(loginRec, loginReq)

	assert.Equal(t, http.StatusOK, loginRec.Code)

	var loginResponse map[string]string
	json.Unmarshal(loginRec.Body.Bytes(), &loginResponse)
	assert.Contains(t, loginResponse, "token")
}

func TestCompanyOperations(t *testing.T) {
	token := getJWTToken(t)

	companyPayload := map[string]interface{}{
		"ID":                  uuid.New().String(),
		"name":                "TestCompany",
		"description":         "A test company",
		"amount_of_employees": 50,
		"registered":          true,
		"type":                "Corporations",
	}

	companyBody, _ := json.Marshal(companyPayload)
	createReq := httptest.NewRequest("POST", "/v1/companies", bytes.NewBuffer(companyBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createRec := httptest.NewRecorder()
	testRouter.ServeHTTP(createRec, createReq)

	assert.Equal(t, http.StatusCreated, createRec.Code)

	var createdCompany entity.Company
	json.Unmarshal(createRec.Body.Bytes(), &createdCompany)

	getReq := httptest.NewRequest("GET", fmt.Sprintf("/v1/companies/%s", createdCompany.ID), nil)
	getRec := httptest.NewRecorder()
	testRouter.ServeHTTP(getRec, getReq)

	assert.Equal(t, http.StatusOK, getRec.Code)

	var retrievedCompany entity.Company
	json.Unmarshal(getRec.Body.Bytes(), &retrievedCompany)
	assert.Equal(t, createdCompany.ID, retrievedCompany.ID)

	updatePayload := map[string]interface{}{
		"name": "Ucompany",
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest("PATCH", fmt.Sprintf("/v1/companies/%s", createdCompany.ID), bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateRec := httptest.NewRecorder()
	testRouter.ServeHTTP(updateRec, updateReq)

	assert.Equal(t, http.StatusOK, updateRec.Code)

	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/companies/%s", createdCompany.ID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteRec := httptest.NewRecorder()
	testRouter.ServeHTTP(deleteRec, deleteReq)

	assert.Equal(t, http.StatusNoContent, deleteRec.Code)
}

func getJWTToken(t *testing.T) string {
	loginPayload := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	loginBody, _ := json.Marshal(loginPayload)
	loginReq := httptest.NewRequest("POST", "/v1/users/login", bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	testRouter.ServeHTTP(loginRec, loginReq)

	var loginResponse map[string]string
	json.Unmarshal(loginRec.Body.Bytes(), &loginResponse)
	return loginResponse["token"]
}

func TestOutboxWorker(t *testing.T) {
	companyRepo := repository.NewCompanyRepository(testDB)
	testCompany := &entity.Company{
		ID:                uuid.New(),
		Name:              "OutboxCompany",
		Description:       new(string),
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              entity.CompanyType("Corporations"),
	}
	*testCompany.Description = "Test company for outbox"

	outboxEvent := &entity.OutboxEvent{
		ID:        uuid.New(),
		EventType: "company_created",
		Payload:   []byte(`{"id":"` + testCompany.ID.String() + `","name":"OutboxTestCompany"}`),
		CreatedAt: time.Now(),
	}

	err := companyRepo.CreateWithOutboxEvent(context.Background(), testCompany, outboxEvent)
	require.NoError(t, err)

	mockProducer := &mockKafkaProducer{}
	outboxWorker := worker.NewOutboxWorker(companyRepo, mockProducer, logger.NewLogger("debug"))
	err = outboxWorker.ProcessOutboxEvents(context.Background())
	require.NoError(t, err)

	events, err := companyRepo.GetOutboxEvents(context.Background(), 10)
	require.NoError(t, err)
	assert.Empty(t, events, "Outbox should be empty after processing")

	assert.True(t, mockProducer.MessageSent, "Message should have been sent to Kafka")
}

type mockKafkaProducer struct {
	MessageSent bool
}

func (m *mockKafkaProducer) Produce(ctx context.Context, topic string, key, value []byte) error {
	m.MessageSent = true
	return nil
}

func (m *mockKafkaProducer) Close() error {
	return nil
}

func cleanupDatabase(log logger.Logger) {
	tx, err := testDB.Begin(context.Background())
	if err != nil {
		log.Error("Failed to begin transaction for cleanup", "error", err)
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
		DROP TABLE users, companies, outbox_events, goose_db_version;
	`)
	if err != nil {
		log.Error("Failed to drop tables", "error", err)
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Error("Failed to commit transaction during cleanup", "error", err)
	}
}
