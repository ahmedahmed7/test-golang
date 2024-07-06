package internal

import (
	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
)

type App struct {
	logger *charmLog.Logger
}

func NewApp(logger *charmLog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) RegisterRoutes(r *mux.Router) {
	//TODO: Implement routes
}
