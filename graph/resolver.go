package graph

import (
	"github.com/olzzhas/narxozer/internal/data"
	"github.com/olzzhas/narxozer/internal/jsonlog"
	"time"
)

type Resolver struct {
	Models     data.Models
	Logger     *jsonlog.Logger
	JWTManager JWTManager
}

func NewResolver(models data.Models, logger *jsonlog.Logger) *Resolver {
	return &Resolver{
		Models:     models,
		Logger:     logger,
		JWTManager: *NewJWTManager("JWT_MANAGER_SECRET", 24*time.Hour),
	}
}
