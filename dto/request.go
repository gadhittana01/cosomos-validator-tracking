package dto

type GetHourlySnapshotRequest struct {
	ValidatorAddress string `json:"validatorAddress" validate:"required"`
	Limit            int32  `json:"limit" validate:"required"`
	Page             int32  `json:"page" validate:"required"`
}

type GetDailySnapshotRequest struct {
	ValidatorAddress string `json:"validatorAddress" validate:"required"`
	Limit            int32  `json:"limit" validate:"required"`
	Page             int32  `json:"page" validate:"required"`
}

type GetDelegatorHistoryRequest struct {
	ValidatorAddress string `json:"validatorAddress" validate:"required"`
	DelegatorAddress string `json:"delegatorAddress" validate:"required"`
	SortBy           string `json:"sortBy" validate:"required"`
	Limit            int32  `json:"limit" validate:"required"`
	Page             int32  `json:"page" validate:"required"`
}
