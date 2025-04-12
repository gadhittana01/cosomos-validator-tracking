package handler

import (
	"net/http"

	"github.com/gadhittana01/cosmos-validation-tracking/constant"
	"github.com/gadhittana01/cosmos-validation-tracking/dto"
	"github.com/gadhittana01/cosmos-validation-tracking/service"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/go-chi/chi"
)

type ValidatorHandler interface {
	SetupValidatorRoutes(route *chi.Mux)
}

type ValidatorHandlerImpl struct {
	validatorService service.ValidatorSvc
	logger           utils.LoggerSvc
}

func NewValidatorHandler(validatorService service.ValidatorSvc, logger utils.LoggerSvc) ValidatorHandler {
	return &ValidatorHandlerImpl{
		validatorService: validatorService,
		logger:           logger,
	}
}

// GetHourlyDelegationSnapshot godoc
// @Id getHourlyDelegationSnapshot
// @Summary      Get Hourly Delegation Snapshot
// @Description  Get Hourly Delegation Snapshot
// @Tags         validator
// @Accept 		 json
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200{data=dto.PaginationResp[dto.GetHourlySnapshotResponse]}
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /api/v1/validators/{validatorAddress}/delegations/hourly [get]
func (h *ValidatorHandlerImpl) GetHourlyDelegationSnapshot(w http.ResponseWriter, r *http.Request) {
	validatorAddress := utils.ValidateURLParamString(r, "validatorAddress")
	page := utils.ValidateQueryParamInt(r, "page", constant.DefaultPage)
	limit := utils.ValidateQueryParamInt(r, "limit", constant.DefaultLimit)

	resp := h.validatorService.GetHourlySnapshot(r.Context(), dto.GetHourlySnapshotRequest{
		ValidatorAddress: validatorAddress,
		Page:             int32(page),
		Limit:            int32(limit),
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

// GetDailyDelegationSnapshot godoc
// @Id getDailyDelegationSnapshot
// @Summary      Get Daily Delegation Snapshot
// @Description  Get Daily Delegation Snapshot
// @Tags         validator
// @Accept 		 json
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200{data=dto.PaginationResp[dto.GetDailySnapshotResponse]}
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /api/v1/validators/{validatorAddress}/delegations/daily [get]
func (h *ValidatorHandlerImpl) GetDailyDelegationSnapshot(w http.ResponseWriter, r *http.Request) {
	validatorAddress := utils.ValidateURLParamString(r, "validatorAddress")
	page := utils.ValidateQueryParamInt(r, "page", constant.DefaultPage)
	limit := utils.ValidateQueryParamInt(r, "limit", constant.DefaultLimit)

	resp := h.validatorService.GetDailySnapshot(r.Context(), dto.GetDailySnapshotRequest{
		ValidatorAddress: validatorAddress,
		Page:             int32(page),
		Limit:            int32(limit),
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

// GetDelegatorHistory godoc
// @Id getDelegatorHistory
// @Summary      Get Delegator History
// @Description  Get Delegator History
// @Tags         validator
// @Accept 		 json
// @Produce      json
// @Success      200  {object}  dto.SuccessResp200{data=dto.PaginationResp[dto.GetDelegatorHistoryResponse]}
// @Failure      400  {object}  dto.FailedResp400
// @Failure      401  {object}  dto.FailedResp401
// @Failure      404  {object}  dto.FailedResp404
// @Failure      500  {object}  dto.FailedResp500
// @Router       /api/v1/validators/{validatorAddress}/delegator/{delegatorAddress}/history [get]
func (h *ValidatorHandlerImpl) GetDelegatorHistory(w http.ResponseWriter, r *http.Request) {
	validatorAddress := utils.ValidateURLParamString(r, "validatorAddress")
	delegatorAddress := utils.ValidateURLParamString(r, "delegatorAddress")
	sortBy := utils.ValidateURLParamString(r, "sortBy", "date")
	page := utils.ValidateQueryParamInt(r, "page", constant.DefaultPage)
	limit := utils.ValidateQueryParamInt(r, "limit", constant.DefaultLimit)

	resp := h.validatorService.GetDelegatorHistory(r.Context(), dto.GetDelegatorHistoryRequest{
		ValidatorAddress: validatorAddress,
		DelegatorAddress: delegatorAddress,
		SortBy:           sortBy,
		Page:             int32(page),
		Limit:            int32(limit),
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *ValidatorHandlerImpl) SetupValidatorRoutes(route *chi.Mux) {
	setupValidatorV1Routes(route, h)
}

func setupValidatorV1Routes(route *chi.Mux, h *ValidatorHandlerImpl) {
	route.Get("/api/v1/validators/{validatorAddress}/delegations/hourly", h.GetHourlyDelegationSnapshot)
	route.Get("/api/v1/validators/{validatorAddress}/delegations/daily", h.GetDailyDelegationSnapshot)
	route.Get("/api/v1/validators/{validatorAddress}/delegator/{delegatorAddress}/history", h.GetDelegatorHistory)
}
