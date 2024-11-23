package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Corray333/bumblebee/internal/config"
	"github.com/Corray333/bumblebee/internal/domains/product"
	"github.com/Corray333/bumblebee/internal/storage"
	"github.com/Corray333/bumblebee/pkg/server/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type controller interface {
	Build()
	Run()
}

type App struct {
	server      *http.Server
	controllers []controller
}

func (app *App) AddController(c controller) {
	app.controllers = append(app.controllers, c)
}
func newRouter() *chi.Mux {
	router := chi.NewMux()
	router.Use(logger.New(slog.Default()))

	// TODO: get allowed origins, headers and methods from cfg
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Set-Cookie", "Refresh", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(c.Handler)

	router.Get("/api/swagger/*", httpSwagger.WrapHandler)

	return router
}

func New() *App {
	config.MustInit(os.Args[1])

	app := &App{}

	router := newRouter()

	// TODO: add timeouts
	server := &http.Server{
		Addr:    "0.0.0.0:" + viper.GetString("port"),
		Handler: router,
	}

	app.server = server

	store, err := storage.New()
	if err != nil {
		panic(err)
	}

	// fileManager := files.New()

	productController := product.NewDomainController(router, store)
	app.AddController(productController)

	return app
}

func (app *App) Init() *App {
	for _, c := range app.controllers {
		c.Build()
	}
	return app
}

func (app *App) Run() {
	slog.Info("Server started at " + app.server.Addr)
	for _, c := range app.controllers {
		go c.Run()
	}
	slog.Error(app.server.ListenAndServe().Error())
}
