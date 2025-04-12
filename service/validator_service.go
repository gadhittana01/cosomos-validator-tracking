package service

import (
	"context"
	"net/http"

	"github.com/gadhittana01/cosmos-validation-tracking/constant"
	querier "github.com/gadhittana01/cosmos-validation-tracking/db/repository"
	"github.com/gadhittana01/cosmos-validation-tracking/dto"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type (
	PaginationValidatorSnapshotResp         = dto.PaginationResp[dto.GetHourlySnapshotResponse]
	PaginationValidatorDailySnapshotResp    = dto.PaginationResp[dto.GetDailySnapshotResponse]
	PaginationValidatorDelegatorHistoryResp = dto.PaginationResp[dto.GetDelegatorHistoryResponse]
)

type ValidatorSvc interface {
	GetHourlySnapshot(ctx context.Context, req dto.GetHourlySnapshotRequest) PaginationValidatorSnapshotResp
	GetDailySnapshot(ctx context.Context, req dto.GetDailySnapshotRequest) PaginationValidatorDailySnapshotResp
	GetDelegatorHistory(ctx context.Context, req dto.GetDelegatorHistoryRequest) PaginationValidatorDelegatorHistoryResp
}

type validatorSvc struct {
	repo     querier.Repository
	logger   utils.LoggerSvc
	cacheSvc utils.CacheSvc
}

func NewValidatorSvc(repo querier.Repository, logger utils.LoggerSvc, cacheSvc utils.CacheSvc) ValidatorSvc {
	return &validatorSvc{
		repo:     repo,
		logger:   logger,
		cacheSvc: cacheSvc,
	}
}

func (v *validatorSvc) GetHourlySnapshot(ctx context.Context, req dto.GetHourlySnapshotRequest) dto.PaginationResp[dto.GetHourlySnapshotResponse] {
	resp, err := utils.GetOrSetData(v.cacheSvc, utils.BuildCacheKey(constant.ValidatorHourlySnapshotCacheKey, "", "", req), func() (dto.PaginationResp[dto.GetHourlySnapshotResponse], error) {
		ewg := errgroup.Group{}
		var delegationSnapshot []querier.GetDelegationSnapshotByValidatorRow
		var countDelegationSnapshot int64
		var err1, err2 error

		ewg.Go(func() error {
			delegationSnapshot, err1 = v.repo.GetDelegationSnapshotByValidator(ctx, querier.GetDelegationSnapshotByValidatorParams{
				ValidatorAddress: req.ValidatorAddress,
				Limit:            req.Limit,
				Offset:           dto.GetOffSet(req.Page, req.Limit),
			})
			if err1 != nil {
				return err1
			}

			return nil
		})

		ewg.Go(func() error {
			countDelegationSnapshot, err2 = v.repo.GetCountDelegationSnapshotByValidator(ctx, req.ValidatorAddress)
			if err2 != nil {
				return err2
			}

			return nil
		})

		if err := ewg.Wait(); err != nil {
			return dto.PaginationResp[dto.GetHourlySnapshotResponse]{}, utils.CustomErrorWithTrace(err, "failed to get hourly snapshot", http.StatusUnprocessableEntity)
		}

		return dto.ToPaginationResp(lo.Map(delegationSnapshot, func(item querier.GetDelegationSnapshotByValidatorRow, _ int) dto.GetHourlySnapshotResponse {
			return dto.GetHourlySnapshotResponse{
				Address:   item.DelegatorAddress,
				Amount:    item.AmountUatom,
				Change:    item.ChangeUatom,
				Date:      item.Timestamp.Format(constant.DateFormat),
				Timestamp: item.Timestamp.Format(constant.TimeFormat),
			}
		}), int(req.Page), int(req.Limit), int(countDelegationSnapshot)), nil
	})
	utils.PanicIfAppError(err, "failed to get hourly snapshot", http.StatusUnprocessableEntity)

	return resp

}

func (v *validatorSvc) GetDailySnapshot(ctx context.Context, req dto.GetDailySnapshotRequest) dto.PaginationResp[dto.GetDailySnapshotResponse] {
	resp, err := utils.GetOrSetData(v.cacheSvc, utils.BuildCacheKey(constant.ValidatorDailySnapshotCacheKey, "", "", req), func() (dto.PaginationResp[dto.GetDailySnapshotResponse], error) {
		ewg := errgroup.Group{}
		var delegationSnapshot []querier.GetDailyAggregateByValidatorRow
		var countDelegationSnapshot int64
		var err1, err2 error

		ewg.Go(func() error {
			delegationSnapshot, err1 = v.repo.GetDailyAggregateByValidator(ctx, querier.GetDailyAggregateByValidatorParams{
				ValidatorAddress: req.ValidatorAddress,
				Limit:            req.Limit,
				Offset:           dto.GetOffSet(req.Page, req.Limit),
			})
			if err1 != nil {
				return err1
			}

			return nil
		})

		ewg.Go(func() error {
			countDelegationSnapshot, err2 = v.repo.GetCountDailyAggregateByValidator(ctx, req.ValidatorAddress)
			if err2 != nil {
				return err2
			}

			return nil
		})

		if err := ewg.Wait(); err != nil {
			return dto.PaginationResp[dto.GetDailySnapshotResponse]{}, utils.CustomErrorWithTrace(err, "failed to get daily aggregate by validator", http.StatusUnprocessableEntity)
		}

		return dto.ToPaginationResp(lo.Map(delegationSnapshot, func(item querier.GetDailyAggregateByValidatorRow, _ int) dto.GetDailySnapshotResponse {
			return dto.GetDailySnapshotResponse{
				Address: item.DelegatorAddress,
				Date:    item.Date.Format(constant.DateFormat),
				Total:   item.TotalAmount,
			}
		}), int(req.Page), int(req.Limit), int(countDelegationSnapshot)), nil
	})
	utils.PanicIfAppError(err, "failed to get daily snapshot", http.StatusUnprocessableEntity)

	return resp
}

func (v *validatorSvc) GetDelegatorHistory(ctx context.Context, req dto.GetDelegatorHistoryRequest) dto.PaginationResp[dto.GetDelegatorHistoryResponse] {

	resp, err := utils.GetOrSetData(v.cacheSvc, utils.BuildCacheKey(constant.ValidatorDelegatorHistoryCacheKey, "", "", req), func() (dto.PaginationResp[dto.GetDelegatorHistoryResponse], error) {
		ewg := errgroup.Group{}
		var delegationSnapshot []querier.GetDelegatorHistoryByValidatorRow
		var countDelegationSnapshot int64
		var err1, err2 error

		ewg.Go(func() error {
			delegationSnapshot, err1 = v.repo.GetDelegatorHistoryByValidator(ctx, querier.GetDelegatorHistoryByValidatorParams{
				ValidatorAddress: req.ValidatorAddress,
				DelegatorAddress: req.DelegatorAddress,
				SortBy:           req.SortBy,
				Limit:            req.Limit,
				Offset:           dto.GetOffSet(req.Page, req.Limit),
			})
			if err1 != nil {
				return err1
			}

			return nil
		})

		ewg.Go(func() error {
			countDelegationSnapshot, err2 = v.repo.GetCountDelegatorHistoryByValidator(ctx, querier.GetCountDelegatorHistoryByValidatorParams{
				ValidatorAddress: req.ValidatorAddress,
				DelegatorAddress: req.DelegatorAddress,
			})
			if err2 != nil {
				return err2
			}

			return nil
		})

		if err := ewg.Wait(); err != nil {
			return dto.PaginationResp[dto.GetDelegatorHistoryResponse]{}, utils.CustomErrorWithTrace(err, "failed to get delegator history by validator", http.StatusUnprocessableEntity)
		}

		return dto.ToPaginationResp(lo.Map(delegationSnapshot, func(item querier.GetDelegatorHistoryByValidatorRow, _ int) dto.GetDelegatorHistoryResponse {
			return dto.GetDelegatorHistoryResponse{
				Timestamp: item.Timestamp.Format(constant.TimeFormat),
				Amount:    item.AmountUatom,
				Change:    item.ChangeUatom,
			}
		}), int(req.Page), int(req.Limit), int(countDelegationSnapshot)), nil
	})
	utils.PanicIfAppError(err, "failed to get delegator history by validator", http.StatusUnprocessableEntity)

	return resp
}
