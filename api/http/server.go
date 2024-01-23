package http

import (
	"effective_mobile/api"
	"effective_mobile/config"
	"effective_mobile/internal"
	"effective_mobile/internal/person"
	"effective_mobile/internal/person/delivery/http"
	"effective_mobile/status"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	pprofMDW "github.com/gofiber/fiber/v2/middleware/pprof"
	recoverMDW "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type httpServer struct {
	fiber  *fiber.App
	cfg    *config.Config
	logger *zap.Logger
}

func NewServer(cfg *config.Config, logger *zap.Logger) api.Api {
	return &httpServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *httpServer) Init() error {
	h.fiber = fiber.New(fiber.Config{
		Immutable:    true,
		ProxyHeader:  "X-Forwarded-For",
		ErrorHandler: status.ErrorHandler,
		IdleTimeout:  time.Minute,
	})
	h.fiber.Use(recoverMDW.New(recoverMDW.Config{
		EnableStackTrace: true,
	}))
	h.fiber.Use(pprofMDW.New())

	return nil
}

func (h *httpServer) MapHandlers(app *internal.App) error {
	h.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	h.fiber.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(status.HTTPresponse{
			Status: status.SuccessBstatus,
		})
	})
	ph := http.NewHandler(app.UC["person"].(person.UC))
	http.MapRoutes(h.fiber, ph)
	return nil
}

func (h *httpServer) Serve() error {
	l := h.logger.Sugar()
	l.Infof("start http server on %s:%s", h.cfg.HTTPServerHost, h.cfg.HTTPServerPort)

	err := h.fiber.Listen(h.cfg.HTTPServerHost + ":" + h.cfg.HTTPServerPort)
	return err
}

func (h *httpServer) Shutdown() error {
	l := h.logger.Sugar()
	l.Infof("shutdown http server on %s:%s", h.cfg.HTTPServerHost, h.cfg.HTTPServerPort)

	return h.fiber.Shutdown()
}
