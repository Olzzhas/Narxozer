package graph

import (
	"github.com/olzzhas/narxozer/internal/data"
	"github.com/olzzhas/narxozer/internal/jsonlog"
)

type Resolver struct {
	Models data.Models
	Logger *jsonlog.Logger
}

func NewResolver(models data.Models, logger *jsonlog.Logger) *Resolver {
	return &Resolver{
		Models: models,
		Logger: logger,
	}
}
