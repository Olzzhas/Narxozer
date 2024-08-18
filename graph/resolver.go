package graph

import (
	"github.com/olzzhas/narxozer/auth"
	"github.com/olzzhas/narxozer/internal/data"
	"github.com/olzzhas/narxozer/internal/jsonlog"
	"time"
)

type Resolver struct {
	Models     data.Models
	Logger     *jsonlog.Logger
	JWTManager auth.JWTManager
}

func NewResolver(models data.Models, logger *jsonlog.Logger) *Resolver {
	return &Resolver{
		Models:     models,
		Logger:     logger,
		JWTManager: *auth.NewJWTManager("your-secret-key", 24*time.Hour),
	}
}
