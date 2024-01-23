package internal

import (
	"effective_mobile/config"
	"effective_mobile/internal/person/repository"
	"effective_mobile/internal/person/usecase"
	"effective_mobile/pkg/postgres"
	"fmt"

	"go.uber.org/zap"
)

type App struct {
	UC     map[string]interface{}
	Repo   map[string]interface{}
	DB     postgres.Postgres
	cfg    *config.Config
	logger *zap.Logger
}

func NewApp(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		UC:     make(map[string]interface{}),
		Repo:   make(map[string]interface{}),
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Init() error {
	var err error
	a.DB, err = postgres.Connect(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.cfg.PostgresHost,
		a.cfg.PostgresPort,
		a.cfg.PostgresUser,
		a.cfg.PostgresPassword,
		a.cfg.PostgresDBName,
		a.cfg.PostgresSSLMode))
	if err != nil {
		return err
	}

	a.Repo["person"] = repository.NewPostgres(a.DB)
	a.UC["person"] = usecase.New(a.DB, a.cfg, a.logger)
	return nil
}
