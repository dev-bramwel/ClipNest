package app

type App struct {
	Config *config.Config
	Router http.Handler
	DB *sql.DB
}

func New(cfg *config.Config) (*App, error) {
	// initialize database
	// initialize storage
	// initialize processor
	// initialize services
	// initialize handlers
	// initialize router

	return &App{}, nil
}

func (a *App) Run() error {
	// start HTTP server
	return nil
}
