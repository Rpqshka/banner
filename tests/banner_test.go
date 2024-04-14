package tests

import (
	"banner"
	"banner/pkg/handler"
	"banner/pkg/repository"
	"banner/pkg/service"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type BannerSuite struct {
	suite.Suite
	db         *sqlx.DB
	repos      *repository.Repository
	services   *service.Service
	handlers   *handler.Handler
	srv        *banner.Server
	userToken  string
	adminToken string
}

func (s *BannerSuite) SetupSuite() {
	logrus.Println("Server not started yet. Starting server...")
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "admin",
		DBName:   "testdb",
		SSLMode:  "disable",
	})
	if err != nil {
		s.T().Fatalf("failed to initialize test db: %s", err.Error())
	}
	s.db = db

	//migrations
	m, err := migrate.New(
		"file://../migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			"postgres", "admin", "localhost", "5432", "testdb"))
	if err != nil {
		logrus.Fatalf("Error creating migration instance: %v", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error applying migrations: %v", err)
	}
	logrus.Debug("Migrations applied successfully")

	s.repos = repository.NewRepository(s.db)
	s.services = service.NewService(s.repos)
	s.handlers = handler.NewHandler(s.services)

	s.srv = new(banner.Server)
	go func() {
		if err := s.srv.Run("8080", s.handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occurred while running http test server: %s", err.Error())
		}
	}()
	logrus.Println("Server started.")
	logrus.Printf("Banner App Integration Tests Started")
	s.initData()
}

func (s *BannerSuite) TearDownSuite() {
	logrus.Println("Shutting down test server...")

	logrus.Println("Closing test db connection...")
	if err := s.db.Close(); err != nil {
		logrus.Errorf("Error occurred on test db connection close: %s", err.Error())
	}
}

func (s *BannerSuite) initData() {
	s.register(userRegister)
	s.register(adminRegister)

	s.userToken = s.login(userLogin)
	s.adminToken = s.login(adminLogin)

	s.createBanner(bannerTest1)
	s.createBanner(bannerNotActive)
	s.createBanner(bannerTest2)
}

func (s *BannerSuite) register(requestBody map[string]interface{}) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req = req.WithContext(c)

	s.handlers.InitRoutes().ServeHTTP(recorder, req)

	if !assert.Equal(s.T(), recorder.Code, http.StatusOK) {
		s.T().FailNow()
	}
}

func (s *BannerSuite) login(requestBody map[string]interface{}) string {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return ""
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return ""
	}

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req = req.WithContext(c)

	s.handlers.InitRoutes().ServeHTTP(recorder, req)

	if !assert.Equal(s.T(), recorder.Code, http.StatusOK) {
		s.T().FailNow()
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	if err != nil {
		s.Fail("Failed to parse JSON body")
		return ""
	}

	token, ok := responseBody["token"].(string)
	if !ok {
		s.Fail("Token not found or invalid format")
		return ""
	}
	return token
}

func (s *BannerSuite) createBanner(requestBody map[string]interface{}) {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return
	}

	reqAdmin, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	//create request by admin
	reqAdmin.Header.Set("Authorization", "Bearer "+s.adminToken)

	recorderAdmin := httptest.NewRecorder()

	cAdmin, _ := gin.CreateTestContext(recorderAdmin)
	reqAdmin = reqAdmin.WithContext(cAdmin)

	s.handlers.InitRoutes().ServeHTTP(recorderAdmin, reqAdmin)

	if !assert.Equal(s.T(), recorderAdmin.Code, http.StatusCreated) {
		s.T().FailNow()
	}

	//create request by user
	reqUser, err := http.NewRequest("POST", "http://localhost:8080/banner", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	reqUser.Header.Set("Authorization", "Bearer "+s.userToken)

	recorderUser := httptest.NewRecorder()

	cUser, _ := gin.CreateTestContext(recorderUser)
	reqUser = reqUser.WithContext(cUser)

	s.handlers.InitRoutes().ServeHTTP(recorderUser, reqUser)

	if !assert.Equal(s.T(), recorderUser.Code, http.StatusForbidden) {
		s.T().FailNow()
	}

}

func (s *BannerSuite) TestGetInactiveBannerByUser() {
	jsonBody, err := json.Marshal(inactiveBannerSearch)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+s.userToken)

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	req = req.WithContext(c)

	s.handlers.InitRoutes().ServeHTTP(recorder, req)
	if !assert.Equal(s.T(), recorder.Code, http.StatusNotFound) {
		s.T().FailNow()
	}
}

func (s *BannerSuite) TestGetInactiveBannerByAdmin() {
	jsonBody, err := json.Marshal(inactiveBannerSearch)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+s.adminToken)

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	req = req.WithContext(c)

	s.handlers.InitRoutes().ServeHTTP(recorder, req)
	logrus.Println(recorder.Body)
	if !assert.Equal(s.T(), recorder.Code, http.StatusOK) {
		s.T().FailNow()
	}
}

func (s *BannerSuite) TestGetBannerByUser() {
	jsonBody, err := json.Marshal(activeBannerSearch)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+s.userToken)

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	req = req.WithContext(c)

	s.handlers.InitRoutes().ServeHTTP(recorder, req)
	if !assert.Equal(s.T(), recorder.Code, http.StatusOK) {
		s.T().FailNow()
	}
}

func (s *BannerSuite) TestGetBannerByAdmin() {
	jsonBody, err := json.Marshal(activeBannerSearch)
	if err != nil {
		s.Fail("Failed to marshal JSON body")
		return
	}

	req, err := http.NewRequest("GET", "http://localhost:8080/user_banner", bytes.NewBuffer(jsonBody))
	if err != nil {
		s.Fail("Failed to create HTTP request")
		return
	}

	req.Header.Set("Authorization", "Bearer "+s.adminToken)

	recorder := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(recorder)
	req = req.WithContext(c)

	s.handlers.InitRoutes().ServeHTTP(recorder, req)
	if !assert.Equal(s.T(), recorder.Code, http.StatusOK) {
		s.T().FailNow()
	}
}

func TestBannerSuite(t *testing.T) {
	suite.Run(t, new(BannerSuite))
}
