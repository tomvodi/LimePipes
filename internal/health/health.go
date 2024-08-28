package health

import (
	"context"
	"fmt"
	healthlib "github.com/alexliesenfeld/health"
	"github.com/alexliesenfeld/health/middleware"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes/internal/config"
	"github.com/tomvodi/limepipes/internal/interfaces"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type healthCheck struct {
	cfg         config.HealthConfig
	gormDb      *gorm.DB
	httpHandler http.Handler
}

func (h *healthCheck) GetCheckHandler() (http.Handler, error) {
	if h.httpHandler == nil {
		return nil, fmt.Errorf("health check handler not initialized")
	}

	return h.httpHandler, nil
}

func (h *healthCheck) init() error {
	db, err := h.gormDb.DB()
	if err != nil {
		return err
	}

	checker := healthlib.NewChecker(
		// Set the time-to-live for our cache to 1 second (default).
		healthlib.WithCacheDuration(time.Duration(h.cfg.CacheDurationSeconds)*time.Second),
		// Configure a global timeout that will be applied to all checks.
		healthlib.WithTimeout(time.Duration(h.cfg.GlobalTimeoutSeconds)*time.Second),

		healthlib.WithPeriodicCheck(
			time.Duration(h.cfg.RefreshPeriodSeconds)*time.Second,
			time.Duration(h.cfg.InitialDelaySeconds)*time.Second,
			healthlib.Check{
				Name:  "database",
				Check: db.PingContext,
			},
		),
		healthlib.WithStatusListener(
			func(_ context.Context, state healthlib.CheckerState) {
				log.Info().Msgf(
					"health status changed to %s", state.Status,
				)
			}),
	)

	h.httpHandler = healthlib.NewHandler(checker,
		healthlib.WithMiddleware(
			middleware.BasicLogger(), // This middleware will log incoming requests.
		),

		healthlib.WithStatusCodeUp(200),

		healthlib.WithStatusCodeDown(503),
	)

	return nil
}

func NewHealthCheck(
	cfg config.HealthConfig,
	gormDb *gorm.DB,
) (interfaces.HealthChecker, error) {
	checker := &healthCheck{
		cfg:    cfg,
		gormDb: gormDb,
	}

	err := checker.init()
	if err != nil {
		return nil, err
	}

	return checker, nil
}
