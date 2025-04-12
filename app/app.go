package app

import (
	"fmt"
	"net/http"

	"github.com/gadhittana01/cosmos-validation-tracking/handler"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	openApiMiddleware "github.com/go-openapi/runtime/middleware"
)

type App interface {
	Start()
}

type AppImpl struct {
	route                 *chi.Mux
	config                *utils.BaseConfig
	validatorHandler      handler.ValidatorHandler
	validatorScheduler    handler.SchedulerHandler
	logger                utils.LoggerSvc
	recoveryMiddlewareSvc utils.RecoveryMiddlewareSvc
}

func NewApp(route *chi.Mux,
	config *utils.BaseConfig,
	validatorHandler handler.ValidatorHandler,
	validatorScheduler handler.SchedulerHandler,
	logger utils.LoggerSvc,
	recoveryMiddlewareSvc utils.RecoveryMiddlewareSvc,
) App {
	return &AppImpl{
		route:                 route,
		config:                config,
		validatorHandler:      validatorHandler,
		validatorScheduler:    validatorScheduler,
		logger:                logger,
		recoveryMiddlewareSvc: recoveryMiddlewareSvc,
	}
}

func (s *AppImpl) Start() {
	s.route.Use(s.recoveryMiddlewareSvc.Recovery)
	s.route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	opts := openApiMiddleware.SwaggerUIOpts{
		SpecURL: s.config.SwaggerURL,
		Path:    "/v1/validator/docs",
	}
	sh := openApiMiddleware.SwaggerUI(opts, nil)
	s.route.Handle("/v1/validator/docs/*", sh)
	s.route.Handle("/swagger.yaml", http.FileServer(http.Dir("./docs")))

	// Health check endpoint
	s.route.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	s.validatorHandler.SetupValidatorRoutes(s.route)
	s.validatorScheduler.SetupSchedulerRoutes(s.route)

	s.route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.GenerateErrorResp[any](w, nil, 404)
	})

	s.logger.Info(fmt.Sprintf("server started on port %d", s.config.ServerPort))
	port := fmt.Sprintf(":%d", s.config.ServerPort)
	err := http.ListenAndServe(port, s.route)
	if err != nil {
		panic(err)
	}
}
