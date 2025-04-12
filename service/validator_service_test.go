package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gadhittana01/cosmos-validation-tracking/constant"
	querier "github.com/gadhittana01/cosmos-validation-tracking/db/repository"
	mockrepo "github.com/gadhittana01/cosmos-validation-tracking/db/repository/mock"
	"github.com/gadhittana01/cosmos-validation-tracking/dto"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	mockutl "github.com/gadhittana01/cosmos-validation-tracking/utils/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var errInvalidReq = errors.New("invalid request")

func initValidatorSvc(
	t *testing.T,
	ctrl *gomock.Controller,
	config *utils.BaseConfig,
) (ValidatorSvc, *mockrepo.MockRepository, *mockutl.MockLoggerSvc, utils.CacheSvc) {
	mockRepo := mockrepo.NewMockRepository(ctrl)
	mockLogger := mockutl.NewMockLoggerSvc(ctrl)
	cacheSvc := utils.InitCacheSvc(t, config, mockLogger)

	return NewValidatorSvc(mockRepo, mockLogger, cacheSvc), mockRepo, mockLogger, cacheSvc
}

func TestGetHourlySnapshot(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	validatorSvcMock, mockRepo, mockLogger, cacheSvc := initValidatorSvc(t, ctrl, config)
	request := dto.GetHourlySnapshotRequest{
		ValidatorAddress: "cosmosvaloper1...",
		Limit:            10,
		Page:             1,
	}
	timestamp := time.Now()
	response := dto.GetHourlySnapshotResponse{
		Address:   "cosmos1...",
		Amount:    1000,
		Change:    1000,
		Date:      timestamp.Format(constant.DateFormat),
		Timestamp: timestamp.Format(constant.TimeFormat),
	}
	mockutl.LoggerMock(mockLogger)

	t.Run("success get hourly snapshot", func(t *testing.T) {

		mockRepo.EXPECT().GetDelegationSnapshotByValidator(gomock.Any(), querier.GetDelegationSnapshotByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDelegationSnapshotByValidatorRow{
			{
				DelegatorAddress: response.Address,
				AmountUatom:      response.Amount,
				Timestamp:        timestamp,
				ChangeUatom:      response.Change,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetCountDelegationSnapshotByValidator(gomock.Any(), request.ValidatorAddress).Return(int64(1), nil).Times(1)

		resp := validatorSvcMock.GetHourlySnapshot(ctx, request)

		assert.NotEmpty(t, resp)
		assert.Equal(t, response, resp.Data[0])
	})

	t.Run("success get hourly snapshot (from cache)", func(t *testing.T) {
		resp := validatorSvcMock.GetHourlySnapshot(ctx, request)

		assert.NotEmpty(t, resp)
		assert.Equal(t, response, resp.Data[0])
	})

	t.Run("failed get count hourly snapshot", func(t *testing.T) {
		cacheSvc.DelByPrefix(ctx, constant.ValidatorHourlySnapshotCacheKey)

		mockRepo.EXPECT().GetDelegationSnapshotByValidator(gomock.Any(), querier.GetDelegationSnapshotByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDelegationSnapshotByValidatorRow{
			{
				DelegatorAddress: response.Address,
				AmountUatom:      response.Amount,
				Timestamp:        timestamp,
				ChangeUatom:      response.Change,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetCountDelegationSnapshotByValidator(gomock.Any(), request.ValidatorAddress).Return(int64(0), errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    fmt.Sprintf("invalid request|%s", "failed to get hourly snapshot"),
		}, func() {
			validatorSvcMock.GetHourlySnapshot(ctx, request)
		})
	})

	t.Run("failed get hourly snapshot", func(t *testing.T) {
		mockRepo.EXPECT().GetDelegationSnapshotByValidator(gomock.Any(), querier.GetDelegationSnapshotByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDelegationSnapshotByValidatorRow{}, errInvalidReq).Times(1)

		mockRepo.EXPECT().GetCountDelegationSnapshotByValidator(gomock.Any(), request.ValidatorAddress).Return(int64(1), nil).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    fmt.Sprintf("invalid request|%s", "failed to get hourly snapshot"),
		}, func() {
			validatorSvcMock.GetHourlySnapshot(ctx, request)
		})
	})

}

func TestGetDailySnapshot(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	validatorSvcMock, mockRepo, mockLogger, cacheSvc := initValidatorSvc(t, ctrl, config)

	request := dto.GetDailySnapshotRequest{
		ValidatorAddress: "cosmosvaloper1...",
		Limit:            10,
		Page:             1,
	}
	timestamp := time.Now()
	response := dto.GetDailySnapshotResponse{
		Address: "cosmos1...",
		Date:    timestamp.Format(constant.DateFormat),
		Total:   1000,
	}
	mockutl.LoggerMock(mockLogger)

	t.Run("success get daily snapshot", func(t *testing.T) {
		mockRepo.EXPECT().GetDailyAggregateByValidator(gomock.Any(), querier.GetDailyAggregateByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDailyAggregateByValidatorRow{
			{
				DelegatorAddress: response.Address,
				Date:             timestamp,
				TotalAmount:      response.Total,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetCountDailyAggregateByValidator(gomock.Any(), request.ValidatorAddress).Return(int64(1), nil).Times(1)

		resp := validatorSvcMock.GetDailySnapshot(ctx, request)

		assert.NotEmpty(t, resp)
		assert.Equal(t, response, resp.Data[0])
	})

	t.Run("success get daily snapshot (from cache)", func(t *testing.T) {
		resp := validatorSvcMock.GetDailySnapshot(ctx, request)

		assert.NotEmpty(t, resp)
		assert.Equal(t, response, resp.Data[0])
	})

	t.Run("failed get count daily snapshot", func(t *testing.T) {
		cacheSvc.DelByPrefix(ctx, constant.ValidatorDailySnapshotCacheKey)

		mockRepo.EXPECT().GetDailyAggregateByValidator(gomock.Any(), querier.GetDailyAggregateByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDailyAggregateByValidatorRow{
			{
				DelegatorAddress: response.Address,
				Date:             timestamp,
				TotalAmount:      response.Total,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetCountDailyAggregateByValidator(gomock.Any(), request.ValidatorAddress).Return(int64(0), errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    fmt.Sprintf("invalid request|%s", "failed to get daily aggregate by validator"),
		}, func() {
			validatorSvcMock.GetDailySnapshot(ctx, request)
		})
	})

	t.Run("failed get daily snapshot", func(t *testing.T) {
		mockRepo.EXPECT().GetDailyAggregateByValidator(gomock.Any(), querier.GetDailyAggregateByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDailyAggregateByValidatorRow{}, errInvalidReq).Times(1)

		mockRepo.EXPECT().GetCountDailyAggregateByValidator(gomock.Any(), request.ValidatorAddress).Return(int64(1), nil).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    fmt.Sprintf("invalid request|%s", "failed to get daily aggregate by validator"),
		}, func() {
			validatorSvcMock.GetDailySnapshot(ctx, request)
		})
	})

}

func TestGetDelegatorHistory(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &utils.BaseConfig{}
	utils.LoadBaseConfig("../config", "test", config)
	validatorSvcMock, mockRepo, mockLogger, cacheSvc := initValidatorSvc(t, ctrl, config)

	request := dto.GetDelegatorHistoryRequest{
		ValidatorAddress: "cosmosvaloper1...",
		Limit:            10,
		Page:             1,
	}
	timestamp := time.Now()
	response := dto.GetDelegatorHistoryResponse{
		Timestamp: timestamp.Format(constant.TimeFormat),
		Amount:    1000,
		Change:    1000,
	}
	mockutl.LoggerMock(mockLogger)

	t.Run("success get delegator history", func(t *testing.T) {
		mockRepo.EXPECT().GetDelegatorHistoryByValidator(gomock.Any(), querier.GetDelegatorHistoryByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			DelegatorAddress: request.DelegatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDelegatorHistoryByValidatorRow{
			{
				Timestamp:   timestamp,
				AmountUatom: response.Amount,
				ChangeUatom: response.Change,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetCountDelegatorHistoryByValidator(gomock.Any(), querier.GetCountDelegatorHistoryByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			DelegatorAddress: request.DelegatorAddress,
		}).Return(int64(1), nil).Times(1)

		resp := validatorSvcMock.GetDelegatorHistory(ctx, request)

		assert.NotEmpty(t, resp)
		assert.Equal(t, response, resp.Data[0])
	})

	t.Run("success get delegator history (from cache)", func(t *testing.T) {
		resp := validatorSvcMock.GetDelegatorHistory(ctx, request)

		assert.NotEmpty(t, resp)
		assert.Equal(t, response, resp.Data[0])
	})

	t.Run("failed get count delegator history", func(t *testing.T) {
		cacheSvc.DelByPrefix(ctx, constant.ValidatorDelegatorHistoryCacheKey)

		mockRepo.EXPECT().GetDelegatorHistoryByValidator(gomock.Any(), querier.GetDelegatorHistoryByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			DelegatorAddress: request.DelegatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDelegatorHistoryByValidatorRow{
			{
				Timestamp:   timestamp,
				AmountUatom: response.Amount,
				ChangeUatom: response.Change,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().GetCountDelegatorHistoryByValidator(gomock.Any(), querier.GetCountDelegatorHistoryByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			DelegatorAddress: request.DelegatorAddress,
		}).Return(int64(0), errInvalidReq).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    fmt.Sprintf("invalid request|%s", "failed to get delegator history by validator"),
		}, func() {
			validatorSvcMock.GetDelegatorHistory(ctx, request)
		})
	})

	t.Run("failed get delegator history", func(t *testing.T) {
		mockRepo.EXPECT().GetDelegatorHistoryByValidator(gomock.Any(), querier.GetDelegatorHistoryByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			DelegatorAddress: request.DelegatorAddress,
			Limit:            request.Limit,
			Offset:           dto.GetOffSet(request.Page, request.Limit),
		}).Return([]querier.GetDelegatorHistoryByValidatorRow{}, errInvalidReq).Times(1)

		mockRepo.EXPECT().GetCountDelegatorHistoryByValidator(gomock.Any(), querier.GetCountDelegatorHistoryByValidatorParams{
			ValidatorAddress: request.ValidatorAddress,
			DelegatorAddress: request.DelegatorAddress,
		}).Return(int64(1), nil).Times(1)

		assert.PanicsWithValue(t, utils.AppError{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    fmt.Sprintf("invalid request|%s", "failed to get delegator history by validator"),
		}, func() {
			validatorSvcMock.GetDelegatorHistory(ctx, request)
		})
	})

}
