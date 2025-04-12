//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gadhittana01/cosmos-validation-tracking/app"
	querier "github.com/gadhittana01/cosmos-validation-tracking/db/repository"
	"github.com/gadhittana01/cosmos-validation-tracking/handler"
	"github.com/gadhittana01/cosmos-validation-tracking/scheduler"
	"github.com/gadhittana01/cosmos-validation-tracking/service"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

var validatorHandlerSet = wire.NewSet(
	querier.NewRepository,
	handler.NewValidatorHandler,
	service.NewValidatorSvc,
)

var loggerSet = wire.NewSet(
	utils.NewLogger,
)

var recoveryMiddlewareSet = wire.NewSet(
	utils.NewRecoveryMiddlewareSvc,
)

var httpClientSet = wire.NewSet(
	utils.NewDefaultHTTPClient,
)

var validatorSchedulerSet = wire.NewSet(
	scheduler.NewValidatorScheduler,
	handler.NewSchedulerHandler,
)

var cacheSet = wire.NewSet(
	wire.Bind(new(utils.RedisClient), new(*redis.Client)),
	utils.NewRedisClient,
	utils.NewCacheSvc,
)

func InitializeApp(
	route *chi.Mux,
	DB utils.PGXPool,
	config *utils.BaseConfig,
) (app.App, error) {
	wire.Build(
		validatorHandlerSet,
		app.NewApp,
		loggerSet,
		recoveryMiddlewareSet,
		httpClientSet,
		validatorSchedulerSet,
		cacheSet,
	)

	return nil, nil
}
