package scheduler

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gadhittana01/cosmos-validation-tracking/constant"
	querier "github.com/gadhittana01/cosmos-validation-tracking/db/repository"
	"github.com/gadhittana01/cosmos-validation-tracking/scheduler/message"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ValidatorScheduler interface {
	SchedulerForHourlyCollectValidatorData(ctx context.Context)
	SchedulerForDailyCollectValidatorData(ctx context.Context)
}

type ValidatorSchedulerImpl struct {
	repo       querier.Repository
	config     *utils.BaseConfig
	logger     utils.LoggerSvc
	httpClient utils.HTTPClient
	cache      utils.CacheSvc
}

func NewValidatorScheduler(
	repo querier.Repository,
	config *utils.BaseConfig,
	logger utils.LoggerSvc,
	httpClient utils.HTTPClient,
	cache utils.CacheSvc,
) ValidatorScheduler {
	return &ValidatorSchedulerImpl{
		repo:       repo,
		config:     config,
		logger:     logger,
		httpClient: httpClient,
		cache:      cache,
	}
}

func (s *ValidatorSchedulerImpl) SchedulerForHourlyCollectValidatorData(ctx context.Context) {
	s.logger.Info("Scheduler for collect validator data")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := utils.ExecTxPoolWithRetry(ctx, s.repo.GetDB(), constant.RetryCount, func(tx pgx.Tx) error {
			repoTx := s.repo.WithTx(tx)

			response, err := s.httpClient.Get(ctx, s.config.CosmosAPIURL)
			if err != nil {
				s.logger.Error("Error getting validator data", zap.Error(err))
				return err
			}

			var data message.CosmosAPIResponse
			err = json.Unmarshal([]byte(response.Body), &data)
			if err != nil {
				s.logger.Error("Error unmarshalling validator data", zap.Error(err))
				return err
			}

			for _, delegation := range data.DelegationResponses {
				delegationSnapshot, err := repoTx.GetDelegationSnapshotByValidatorAndDelegator(ctx, querier.GetDelegationSnapshotByValidatorAndDelegatorParams{
					ValidatorAddress: delegation.Delegation.ValidatorAddress,
					DelegatorAddress: delegation.Delegation.DelegatorAddress,
				})
				if err != nil && err != pgx.ErrNoRows {
					s.logger.Error("Error getting delegation snapshot", zap.Error(err))
					return err
				}

				currentUatom, err := strconv.ParseInt(delegation.Balance.Amount, 10, 64)
				if err != nil {
					s.logger.Error("Error parsing current uatom", zap.Error(err))
					return err
				}

				changeUatom := currentUatom - delegationSnapshot.AmountUatom

				_, err = repoTx.CreateDelegationSnapshot(ctx, querier.CreateDelegationSnapshotParams{
					ValidatorAddress: delegation.Delegation.ValidatorAddress,
					DelegatorAddress: delegation.Delegation.DelegatorAddress,
					AmountUatom:      currentUatom,
					ChangeUatom:      changeUatom,
					Timestamp:        utils.GetCurrentTimeInJakarta(),
				})

				if err != nil {
					s.logger.Error("Error creating delegation snapshot", zap.Error(err))
					return err
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Error("Error executing transaction", zap.Error(err))
		}

		s.cache.ClearCaches([]string{constant.ValidatorHourlySnapshotCacheKey}, "")
		s.cache.ClearCaches([]string{constant.ValidatorDelegatorHistoryCacheKey}, "")
		s.logger.Info("Successfully collected hourly validator data")
	}()
}

func (s *ValidatorSchedulerImpl) SchedulerForDailyCollectValidatorData(ctx context.Context) {
	s.logger.Info("Scheduler for collect validator data")

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := utils.ExecTxPoolWithRetry(ctx, s.repo.GetDB(), constant.RetryCount, func(tx pgx.Tx) error {
			repoTx := s.repo.WithTx(tx)

			delegationSnapshot, err := repoTx.GetLatestDelegationSnapshot(ctx)
			if err != nil {
				s.logger.Error("Error getting latest delegation snapshot", zap.Error(err))
				return err
			}

			for _, delegation := range delegationSnapshot {
				_, err = repoTx.CreateDailyAggregate(ctx, querier.CreateDailyAggregateParams{
					ValidatorAddress: delegation.ValidatorAddress,
					DelegatorAddress: delegation.DelegatorAddress,
					Date:             utils.GetCurrentTimeInJakarta(),
					TotalAmount:      delegation.AmountUatom,
				})
				if err != nil {
					s.logger.Error("Error creating daily aggregate", zap.Error(err))
					return err
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Error("Error executing transaction", zap.Error(err))
		}

		s.cache.ClearCaches([]string{constant.ValidatorDailySnapshotCacheKey}, "")
		s.logger.Info("Successfully collected daily validator data")
	}()
}
