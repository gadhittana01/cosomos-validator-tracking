// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gadhittana01/cosmos-validation-tracking/app"
	"github.com/gadhittana01/cosmos-validation-tracking/db/repository"
	"github.com/gadhittana01/cosmos-validation-tracking/handler"
	"github.com/gadhittana01/cosmos-validation-tracking/scheduler"
	"github.com/gadhittana01/cosmos-validation-tracking/service"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// Injectors from injector.go:

func InitializeApp(route *chi.Mux, DB utils.PGXPool, config *utils.BaseConfig) (app.App, error) {
	repository := querier.NewRepository(DB)
	loggerSvc := utils.NewLogger(config)
	client := utils.NewRedisClient(config)
	cacheSvc := utils.NewCacheSvc(config, client, loggerSvc)
	validatorSvc := service.NewValidatorSvc(repository, loggerSvc, cacheSvc)
	validatorHandler := handler.NewValidatorHandler(validatorSvc, loggerSvc)
	httpClient := utils.NewDefaultHTTPClient()
	validatorScheduler := scheduler.NewValidatorScheduler(repository, config, loggerSvc, httpClient, cacheSvc)
	schedulerHandler := handler.NewSchedulerHandler(validatorScheduler, loggerSvc)
	recoveryMiddlewareSvc := utils.NewRecoveryMiddlewareSvc(loggerSvc)
	appApp := app.NewApp(route, config, validatorHandler, schedulerHandler, loggerSvc, recoveryMiddlewareSvc)
	return appApp, nil
}

// injector.go:

var validatorHandlerSet = wire.NewSet(querier.NewRepository, handler.NewValidatorHandler, service.NewValidatorSvc)

var loggerSet = wire.NewSet(utils.NewLogger)

var recoveryMiddlewareSet = wire.NewSet(utils.NewRecoveryMiddlewareSvc)

var httpClientSet = wire.NewSet(utils.NewDefaultHTTPClient)

var validatorSchedulerSet = wire.NewSet(scheduler.NewValidatorScheduler, handler.NewSchedulerHandler)

var cacheSet = wire.NewSet(wire.Bind(new(utils.RedisClient), new(*redis.Client)), utils.NewRedisClient, utils.NewCacheSvc)
