package health

import (
	"context"
	"log/slog"
)

type Health struct {
	log *slog.Logger
}

func (h *Health) CheckHealth(ctx context.Context) {
	const op = "service.health.CheckHealth"

	log := h.log.With(slog.String("op", op))

	log.Info("check health")
}

func (h *Health) WatchHealth() {
	const op = "service.health.WatchHealth"

	log := h.log.With(slog.String("op", op))

	log.Info("watch health")
}

func New(log *slog.Logger) *Health {
	return &Health{
		log: log,
	}
}
