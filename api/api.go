package api

import "effective_mobile/internal"

type Api interface {
	Init() error
	MapHandlers(app *internal.App) error
	Serve() error
	Shutdown() error
}
