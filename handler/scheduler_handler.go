package handler

import (
	"net/http"

	"github.com/gadhittana01/cosmos-validation-tracking/scheduler"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/go-chi/chi"
)

type SchedulerHandler interface {
	SetupSchedulerRoutes(route *chi.Mux)
}

type schedulerHandlerImpl struct {
	validatorScheduler scheduler.ValidatorScheduler
	logger             utils.LoggerSvc
}

func NewSchedulerHandler(validatorScheduler scheduler.ValidatorScheduler, logger utils.LoggerSvc) SchedulerHandler {
	return &schedulerHandlerImpl{
		validatorScheduler: validatorScheduler,
		logger:             logger,
	}
}

// SchedulerForHourlyCollectValidatorData godoc
// @Id schedulerForHourlyCollectValidatorData
// @Summary      Scheduler For Hourly Collect Validator Data
// @Description  Scheduler For Hourly Collect Validator Data
// @Tags         validator
// @Accept 		 json
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /api/v1/scheduler/validator/hourly [get]
func (h *schedulerHandlerImpl) SchedulerForHourlyCollectValidatorData(w http.ResponseWriter, r *http.Request) {
	h.validatorScheduler.SchedulerForHourlyCollectValidatorData(r.Context())

	utils.GenerateSuccessResp[any](w, nil, 200)
}

// SchedulerForDailyCollectValidatorData godoc
// @Id schedulerForDailyCollectValidatorData
// @Summary      Scheduler For Daily Collect Validator Data
// @Description  Scheduler For Daily Collect Validator Data
// @Tags         validator
// @Accept 		 json
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /api/v1/scheduler/validator/daily [get]
func (h *schedulerHandlerImpl) SchedulerForDailyCollectValidatorData(w http.ResponseWriter, r *http.Request) {
	h.validatorScheduler.SchedulerForDailyCollectValidatorData(r.Context())

	utils.GenerateSuccessResp[any](w, nil, 200)
}

func (h *schedulerHandlerImpl) SetupSchedulerRoutes(route *chi.Mux) {
	setupSchedulerV1Routes(route, h)
}

func setupSchedulerV1Routes(route *chi.Mux, h *schedulerHandlerImpl) {
	route.Post("/api/v1/scheduler/validator/hourly", h.SchedulerForHourlyCollectValidatorData)
	route.Post("/api/v1/scheduler/validator/daily", h.SchedulerForDailyCollectValidatorData)
}
