package server

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	ipclient "github.com/M-Fisher/companies_api/app/clients/ipdata"
	"github.com/M-Fisher/companies_api/app/config"
	"github.com/M-Fisher/companies_api/app/internal/services/auth"
	"github.com/M-Fisher/companies_api/app/internal/services/companies"
	"github.com/M-Fisher/companies_api/app/internal/services/events"
	"github.com/M-Fisher/companies_api/app/internal/services/events/clients/kafka"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

type Server struct {
	Router           *mux.Router
	CompaniesService companies.CompaniesService
	EventsService    events.EventsService
	AuthService      auth.AuthService
	Storage          *postgres.DB
	Log              *zap.Logger
	Config           *config.Config
}

func NewServer(cfg *config.Config) *Server {
	srv := Server{
		Config: cfg,
	}
	srv.configureLogger()
	dbService, err := postgres.NewPostgres(&cfg.Postgres, srv.Log)
	if err != nil {
		log.Fatal("Failed to create storage", zap.Error(err))
	}

	kafkaClient, err := kafka.NewKafkaProducer(&cfg.Kafka, srv.Log)
	if err != nil {
		log.Fatal("Kafka failed to dial leader", zap.Error(err))
	}

	evService := events.NewService(&cfg.Kafka, kafkaClient, srv.Log)
	authService := auth.NewService(ipclient.NewClient(cfg.IPApiRequestTimeout))

	srv.Storage = dbService
	srv.AuthService = authService
	srv.CompaniesService = companies.NewService(dbService, authService, evService, srv.Log)
	srv.EventsService = evService

	return &srv
}

func (s *Server) Run() {
	srv := s.initApp()
	s.Log.Info("Starting server app", zap.String("port", s.Config.Port))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Log.Fatal("http server failed", zap.Error(err))
		}
	}()
}

func (s *Server) Stop() {
	s.Log.Info("Stopping server")
	s.Storage.Close()
	err := s.EventsService.Stop()
	if err != nil {
		s.Log.Error("Failed to stop events service", zap.Error(err))
	}
}

func NewTestServer(cfg *config.Config) *httptest.Server {
	srv := NewServer(cfg)
	return httptest.NewServer(srv.createHTTPHandler())
}
