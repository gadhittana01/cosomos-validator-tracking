package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gadhittana01/cosmos-validation-tracking/constant"
	querier "github.com/gadhittana01/cosmos-validation-tracking/db/repository"
	mockrepo "github.com/gadhittana01/cosmos-validation-tracking/db/repository/mock"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	mockutl "github.com/gadhittana01/cosmos-validation-tracking/utils/mock"
	"github.com/gadhittana01/cosmos-validation-tracking/utils/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var errInvalidReq = errors.New("invalid request")

func initValidatorScheduler(t *testing.T, ctrl *gomock.Controller) (
	ValidatorScheduler, *mockrepo.MockRepository, *utils.BaseConfig, *mockutl.MockLoggerSvc, *mockutl.MockHTTPClient,
) {
	mockRepo := mockrepo.NewMockRepository(ctrl)
	config := utils.CheckAndSetConfig("../config", "test")
	mockLogger := mockutl.NewMockLoggerSvc(ctrl)
	mockHTTPClient := mockutl.NewMockHTTPClient(ctrl)
	cacheSvc := utils.InitCacheSvc(t, config, mockLogger)

	return NewValidatorScheduler(mockRepo, config, mockLogger, mockHTTPClient, cacheSvc), mockRepo, config, mockLogger, mockHTTPClient
}

func TestSchedulerForHourlyCollectValidatorData(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	validatorScheduler, mockRepo, config, mockLogger, mockHTTPClient := initValidatorScheduler(t, ctrl)
	mockutl.LoggerMock(mockLogger)
	retryCount := constant.RetryCount + 1

	t.Run("success collect hourly validator data", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo)

		mockHTTPClient.EXPECT().Get(gomock.Any(), config.CosmosAPIURL).Return(&types.HTTPResponse{
			StatusCode: 200,
			Body: `{
				"delegation_responses": [
					{
						"delegation": {
							"delegator_address": "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
							"validator_address": "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
							"shares": "8003.200796626260454171"
						},
						"balance": {
							"denom": "uatom",
							"amount": "8000"
						}
					}
				]
			}`,
			Headers: map[string][]string{},
		}, nil).Times(1)

		mockRepo.EXPECT().GetDelegationSnapshotByValidatorAndDelegator(gomock.Any(), querier.GetDelegationSnapshotByValidatorAndDelegatorParams{
			ValidatorAddress: "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
			DelegatorAddress: "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
		}).Return(querier.GetDelegationSnapshotByValidatorAndDelegatorRow{
			AmountUatom: 8000,
		}, nil).Times(1)

		mockRepo.EXPECT().CreateDelegationSnapshot(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDelegationSnapshotParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDelegationSnapshotParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.AmountUatom)
			assert.Equal(t, int64(0), arg.ChangeUatom)
			return uuid.New(), nil
		}).Times(1)

		validatorScheduler.SchedulerForHourlyCollectValidatorData(ctx)
		time.Sleep(1 * time.Millisecond)
	})

	t.Run("error get validator data", func(t *testing.T) {
		mockrepo.SetupMockTxPoolWithRetry(ctrl, mockRepo, retryCount, true)

		mockHTTPClient.EXPECT().Get(gomock.Any(), config.CosmosAPIURL).Return(&types.HTTPResponse{
			StatusCode: 400,
			Body:       ``,
			Headers:    map[string][]string{},
		}, errInvalidReq).Times(retryCount)

		mockRepo.EXPECT().GetDelegationSnapshotByValidatorAndDelegator(gomock.Any(), querier.GetDelegationSnapshotByValidatorAndDelegatorParams{
			ValidatorAddress: "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
			DelegatorAddress: "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
		}).Return(querier.GetDelegationSnapshotByValidatorAndDelegatorRow{
			AmountUatom: 8000,
		}, nil).Times(0)

		mockRepo.EXPECT().CreateDelegationSnapshot(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDelegationSnapshotParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDelegationSnapshotParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.AmountUatom)
			assert.Equal(t, int64(0), arg.ChangeUatom)
			return uuid.New(), nil
		}).Times(0)

		validatorScheduler.SchedulerForHourlyCollectValidatorData(ctx)
		time.Sleep(1000 * time.Millisecond)
	})

	t.Run("failed get delegation snapshot", func(t *testing.T) {
		mockrepo.SetupMockTxPoolWithRetry(ctrl, mockRepo, retryCount, true)

		mockHTTPClient.EXPECT().Get(gomock.Any(), config.CosmosAPIURL).Return(&types.HTTPResponse{
			StatusCode: 200,
			Body: `{
				"delegation_responses": [
					{
						"delegation": {
							"delegator_address": "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
							"validator_address": "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
							"shares": "8003.200796626260454171"
						},
						"balance": {
							"denom": "uatom",
							"amount": "8000"
						}
					}
				]
			}`,
			Headers: map[string][]string{},
		}, nil).Times(retryCount)

		mockRepo.EXPECT().GetDelegationSnapshotByValidatorAndDelegator(gomock.Any(), querier.GetDelegationSnapshotByValidatorAndDelegatorParams{
			ValidatorAddress: "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
			DelegatorAddress: "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
		}).Return(querier.GetDelegationSnapshotByValidatorAndDelegatorRow{}, errInvalidReq).Times(retryCount)

		mockRepo.EXPECT().CreateDelegationSnapshot(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDelegationSnapshotParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDelegationSnapshotParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.AmountUatom)
			assert.Equal(t, int64(0), arg.ChangeUatom)
			return uuid.New(), nil
		}).Times(0)

		validatorScheduler.SchedulerForHourlyCollectValidatorData(ctx)
		time.Sleep(1000 * time.Millisecond)
	})

	t.Run("failed create delegation snapshot", func(t *testing.T) {
		mockrepo.SetupMockTxPoolWithRetry(ctrl, mockRepo, retryCount, true)

		mockHTTPClient.EXPECT().Get(gomock.Any(), config.CosmosAPIURL).Return(&types.HTTPResponse{
			StatusCode: 200,
			Body: `{
				"delegation_responses": [
					{
						"delegation": {
							"delegator_address": "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
							"validator_address": "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
							"shares": "8003.200796626260454171"
						},
						"balance": {
							"denom": "uatom",
							"amount": "8000"
						}
					}
				]
			}`,
			Headers: map[string][]string{},
		}, nil).Times(retryCount)

		mockRepo.EXPECT().GetDelegationSnapshotByValidatorAndDelegator(gomock.Any(), querier.GetDelegationSnapshotByValidatorAndDelegatorParams{
			ValidatorAddress: "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
			DelegatorAddress: "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
		}).Return(querier.GetDelegationSnapshotByValidatorAndDelegatorRow{
			AmountUatom: 8000,
		}, nil).Times(retryCount)

		mockRepo.EXPECT().CreateDelegationSnapshot(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDelegationSnapshotParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDelegationSnapshotParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.AmountUatom)
			assert.Equal(t, int64(0), arg.ChangeUatom)
			return uuid.New(), errInvalidReq
		}).Times(retryCount)

		validatorScheduler.SchedulerForHourlyCollectValidatorData(ctx)
		time.Sleep(1000 * time.Millisecond)
	})

}

func TestSchedulerForDailyCollectValidatorData(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	validatorScheduler, mockRepo, _, mockLogger, _ := initValidatorScheduler(t, ctrl)
	mockutl.LoggerMock(mockLogger)
	retryCount := constant.RetryCount + 1

	t.Run("success collect daily validator data", func(t *testing.T) {
		mockrepo.SetupMockTxPool(ctrl, mockRepo)

		mockRepo.EXPECT().GetLatestDelegationSnapshot(gomock.Any()).Return([]querier.GetLatestDelegationSnapshotRow{
			{
				ValidatorAddress: "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
				DelegatorAddress: "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
				AmountUatom:      8000,
			},
		}, nil).Times(1)

		mockRepo.EXPECT().CreateDailyAggregate(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDailyAggregateParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDailyAggregateParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.TotalAmount)
			return uuid.New(), nil
		}).Times(1)

		validatorScheduler.SchedulerForDailyCollectValidatorData(ctx)
		time.Sleep(1 * time.Millisecond)
	})

	t.Run("error get latest delegation snapshot", func(t *testing.T) {
		mockrepo.SetupMockTxPoolWithRetry(ctrl, mockRepo, retryCount, true)

		mockRepo.EXPECT().GetLatestDelegationSnapshot(gomock.Any()).Return([]querier.GetLatestDelegationSnapshotRow{}, errInvalidReq).Times(retryCount)

		mockRepo.EXPECT().CreateDailyAggregate(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDailyAggregateParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDailyAggregateParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.TotalAmount)
			return uuid.New(), nil
		}).Times(0)

		validatorScheduler.SchedulerForDailyCollectValidatorData(ctx)
		time.Sleep(1000 * time.Millisecond)
	})

	t.Run("failed create daily aggregate", func(t *testing.T) {
		mockrepo.SetupMockTxPoolWithRetry(ctrl, mockRepo, retryCount, true)

		mockRepo.EXPECT().GetLatestDelegationSnapshot(gomock.Any()).Return([]querier.GetLatestDelegationSnapshotRow{
			{
				ValidatorAddress: "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
				DelegatorAddress: "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500",
				AmountUatom:      8000,
			},
		}, nil).Times(retryCount)

		mockRepo.EXPECT().CreateDailyAggregate(gomock.Any(), gomock.AssignableToTypeOf(querier.CreateDailyAggregateParams{})).DoAndReturn(func(ctx context.Context, arg querier.CreateDailyAggregateParams) (uuid.UUID, error) {
			assert.Equal(t, "cosmosvaloper1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.ValidatorAddress)
			assert.Equal(t, "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500", arg.DelegatorAddress)
			assert.Equal(t, int64(8000), arg.TotalAmount)
			return uuid.New(), errInvalidReq
		}).Times(retryCount)

		validatorScheduler.SchedulerForDailyCollectValidatorData(ctx)
		time.Sleep(1000 * time.Millisecond)
	})
}
