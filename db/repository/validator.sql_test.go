package querier

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

var errQuery = errors.New("error query")

func TestCreateDailyAggregate(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := CreateDailyAggregateParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		DelegatorAddress: "test@gmail.com",
		Date:             time.Now(),
		TotalAmount:      100,
	}
	id := uuid.New()

	t.Run("success create daily aggregate", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createDailyAggregate)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress, req.Date, req.TotalAmount).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(id))

		res, err := q.CreateDailyAggregate(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, id, res)
	})

	t.Run("failed create daily aggregate", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createDailyAggregate)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress, req.Date, req.TotalAmount).
			WillReturnError(errQuery)

		res, err := q.CreateDailyAggregate(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestCreateDelegationSnapshot(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := CreateDelegationSnapshotParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
		AmountUatom:      100,
		ChangeUatom:      100,
		Timestamp:        time.Now(),
	}
	id := uuid.New()

	t.Run("success create delegation snapshot", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createDelegationSnapshot)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress, req.AmountUatom, req.ChangeUatom, req.Timestamp).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(id))

		res, err := q.CreateDelegationSnapshot(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, id, res)
	})

	t.Run("failed create delegation snapshot", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(createDelegationSnapshot)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress, req.AmountUatom, req.ChangeUatom, req.Timestamp).
			WillReturnError(errQuery)

		res, err := q.CreateDelegationSnapshot(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetCountDailyAggregateByValidator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c"
	totalCount := int64(1)

	t.Run("success get count daily aggregate by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getCountDailyAggregateByValidator)).
			WithArgs(req).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(totalCount))

		res, err := q.GetCountDailyAggregateByValidator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, totalCount, res)
	})

	t.Run("failed get count daily aggregate by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getCountDailyAggregateByValidator)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.GetCountDailyAggregateByValidator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetCountDelegationSnapshotByValidator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c"
	totalCount := int64(1)

	t.Run("success get count delegation snapshot by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getCountDelegationSnapshotByValidator)).
			WithArgs(req).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(totalCount))

		res, err := q.GetCountDelegationSnapshotByValidator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, totalCount, res)
	})

	t.Run("failed get count delegation snapshot by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getCountDelegationSnapshotByValidator)).
			WithArgs(req).
			WillReturnError(errQuery)

		res, err := q.GetCountDelegationSnapshotByValidator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetCountDelegatorHistoryByValidator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := GetCountDelegatorHistoryByValidatorParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
	}
	totalCount := int64(1)

	t.Run("success get count delegator history by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getCountDelegatorHistoryByValidator)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress).
			WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(totalCount))

		res, err := q.GetCountDelegatorHistoryByValidator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, totalCount, res)
	})

	t.Run("failed get count delegator history by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getCountDelegatorHistoryByValidator)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress).
			WillReturnError(errQuery)

		res, err := q.GetCountDelegatorHistoryByValidator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetDailyAggregateByValidator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := GetDailyAggregateByValidatorParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		Limit:            10,
		Offset:           0,
	}

	response := []GetDailyAggregateByValidatorRow{
		{
			DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
			Date:             time.Now(),
			TotalAmount:      100,
		},
	}

	t.Run("success get daily aggregate by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDailyAggregateByValidator)).
			WithArgs(req.ValidatorAddress, req.Limit, req.Offset).
			WillReturnRows(pgxmock.NewRows([]string{"delegator_address", "date", "total_amount"}).
				AddRow(response[0].DelegatorAddress, response[0].Date, response[0].TotalAmount))

		res, err := q.GetDailyAggregateByValidator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, response, res)
	})

	t.Run("failed get daily aggregate by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDailyAggregateByValidator)).
			WithArgs(req.ValidatorAddress, req.Limit, req.Offset).
			WillReturnError(errQuery)

		res, err := q.GetDailyAggregateByValidator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetDelegationSnapshotByValidator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := GetDelegationSnapshotByValidatorParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		Limit:            10,
		Offset:           0,
	}

	response := []GetDelegationSnapshotByValidatorRow{
		{
			DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
			AmountUatom:      100,
			Timestamp:        time.Now(),
			ChangeUatom:      100,
		},
	}

	t.Run("success get delegation snapshot by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDelegationSnapshotByValidator)).
			WithArgs(req.ValidatorAddress, req.Limit, req.Offset).
			WillReturnRows(pgxmock.NewRows([]string{"delegator_address", "amount_uatom", "timestamp", "change_uatom"}).
				AddRow(response[0].DelegatorAddress, response[0].AmountUatom, response[0].Timestamp, response[0].ChangeUatom))

		res, err := q.GetDelegationSnapshotByValidator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, response, res)
	})

	t.Run("failed get delegation snapshot by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDelegationSnapshotByValidator)).
			WithArgs(req.ValidatorAddress, req.Limit, req.Offset).
			WillReturnError(errQuery)

		res, err := q.GetDelegationSnapshotByValidator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetDelegationSnapshotByValidatorAndDelegator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := GetDelegationSnapshotByValidatorAndDelegatorParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
	}

	response := GetDelegationSnapshotByValidatorAndDelegatorRow{
		ID:               uuid.New(),
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
		AmountUatom:      100,
		Timestamp:        time.Now(),
		ChangeUatom:      100,
	}

	t.Run("success get delegation snapshot by validator and delegator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDelegationSnapshotByValidatorAndDelegator)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress).
			WillReturnRows(pgxmock.NewRows([]string{"id", "validator_address", "delegator_address", "amount_uatom", "change_uatom", "timestamp"}).
				AddRow(response.ID, response.ValidatorAddress, response.DelegatorAddress, response.AmountUatom, response.ChangeUatom, response.Timestamp))

		res, err := q.GetDelegationSnapshotByValidatorAndDelegator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, response, res)
	})

	t.Run("failed get delegation snapshot by validator and delegator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDelegationSnapshotByValidatorAndDelegator)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress).
			WillReturnError(errQuery)

		res, err := q.GetDelegationSnapshotByValidatorAndDelegator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetDelegatorHistoryByValidator(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	req := GetDelegatorHistoryByValidatorParams{
		ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
		DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
		Limit:            10,
		Offset:           0,
		SortBy:           "date",
	}

	response := []GetDelegatorHistoryByValidatorRow{
		{
			Timestamp:   time.Now(),
			AmountUatom: 100,
			ChangeUatom: 100,
		},
	}

	t.Run("success get delegator history by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDelegatorHistoryByValidator)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress, req.Limit, req.Offset, req.SortBy).
			WillReturnRows(pgxmock.NewRows([]string{"timestamp", "amount_uatom", "change_uatom"}).
				AddRow(response[0].Timestamp, response[0].AmountUatom, response[0].ChangeUatom))

		res, err := q.GetDelegatorHistoryByValidator(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, response, res)
	})

	t.Run("failed get delegator history by validator", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getDelegatorHistoryByValidator)).
			WithArgs(req.ValidatorAddress, req.DelegatorAddress, req.Limit, req.Offset, req.SortBy).
			WillReturnError(errQuery)

		res, err := q.GetDelegatorHistoryByValidator(ctx, req)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestGetLatestDelegationSnapshot(t *testing.T) {
	mockDB, _ := pgxmock.NewPool()
	defer mockDB.Close()
	q := NewRepository(mockDB)
	ctx := context.Background()

	response := []GetLatestDelegationSnapshotRow{
		{
			ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
			DelegatorAddress: "cosmos1pxlmxuzdams3e9j54gdvaell0npa2j695r90jv",
			AmountUatom:      100,
		},
	}

	t.Run("success get latest delegation snapshot", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getLatestDelegationSnapshot)).
			WillReturnRows(pgxmock.NewRows([]string{"validator_address", "delegator_address", "amount_uatom"}).
				AddRow(response[0].ValidatorAddress, response[0].DelegatorAddress, response[0].AmountUatom))

		res, err := q.GetLatestDelegationSnapshot(ctx)
		assert.NoError(t, err)
		assert.Equal(t, response, res)
	})

	t.Run("failed get latest delegation snapshot", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(getLatestDelegationSnapshot)).
			WillReturnError(errQuery)

		res, err := q.GetLatestDelegationSnapshot(ctx)
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}
