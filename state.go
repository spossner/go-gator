package main

import (
	"github.com/spossner/gator/internal/config"
	"github.com/spossner/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
