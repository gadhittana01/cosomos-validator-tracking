package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gadhittana01/cosmos-validation-tracking/dto"
	"github.com/gadhittana01/cosmos-validation-tracking/service"
	mocksvc "github.com/gadhittana01/cosmos-validation-tracking/service/mock"
	"github.com/gadhittana01/cosmos-validation-tracking/utils"
	mockutl "github.com/gadhittana01/cosmos-validation-tracking/utils/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewValidatorHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	validatorMock := mocksvc.NewMockValidatorSvc(ctrl)
	loggerMock := mockutl.NewMockLoggerSvc(ctrl)

	type args struct {
		service service.ValidatorSvc
		logger  utils.LoggerSvc
	}

	tests := []struct {
		name string
		args args
		want *ValidatorHandlerImpl
	}{
		{
			args: args{
				service: validatorMock,
				logger:  loggerMock,
			},
			want: &ValidatorHandlerImpl{
				validatorService: validatorMock,
				logger:           loggerMock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewValidatorHandler(tt.args.service, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValidatorHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHourlyDelegationSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	validatorAddress := "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500"
	page := 1
	limit := 10

	sampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/api/v1/validators/%s/delegations/hourly?page=%d&limit=%d", validatorAddress, page, limit), strings.NewReader(``))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/api/v1/validators/%s/delegations/hourly?page=%d&limit=test", validatorAddress, page), strings.NewReader(``))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.ValidatorSvc
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "success get hourly delegation snapshot",
			fields: func() fields {
				validatorMock := mocksvc.NewMockValidatorSvc(ctrl)

				validatorMock.EXPECT().GetHourlySnapshot(gomock.Any(), dto.GetHourlySnapshotRequest{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetHourlySnapshotResponse]{
					Total:      1,
					IsLoadMore: false,
					Data: []dto.GetHourlySnapshotResponse{
						{
							Address:   validatorAddress,
							Amount:    100,
							Change:    100,
							Date:      "2021-01-01",
							Timestamp: "2021-01-01",
						},
					},
				}).Times(1)

				return fields{
					service: validatorMock,
				}
			},
			args: args{
				w:   sampleResp,
				req: sampleReq,
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			fields: func() fields {
				validatorMock := mocksvc.NewMockValidatorSvc(ctrl)

				validatorMock.EXPECT().GetHourlySnapshot(gomock.Any(), dto.GetHourlySnapshotRequest{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetHourlySnapshotResponse]{
					Total:      0,
					IsLoadMore: false,
					Data:       []dto.GetHourlySnapshotResponse{},
				}).Times(0)

				return fields{
					service: validatorMock,
				}
			},
			args: args{
				w:   invalidSampleResp,
				req: invalidSampleReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields()
			i := ValidatorHandlerImpl{
				validatorService: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.GetHourlyDelegationSnapshot(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.GetHourlyDelegationSnapshot(tt.args.w, tt.args.req)
				})
			}
		})
	}
}

func TestGetDailyDelegationSnapshot(t *testing.T) {
	ctrl := gomock.NewController(t)
	validatorAddress := "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500"
	page := 1
	limit := 10

	sampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/api/v1/validators/%s/delegations/daily?page=%d&limit=%d", validatorAddress, page, limit), strings.NewReader(``))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/api/v1/validators/%s/delegations/daily?page=%d&limit=test", validatorAddress, page), strings.NewReader(``))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.ValidatorSvc
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "success get daily delegation snapshot",
			fields: func() fields {
				validatorMock := mocksvc.NewMockValidatorSvc(ctrl)

				validatorMock.EXPECT().GetDailySnapshot(gomock.Any(), dto.GetDailySnapshotRequest{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetDailySnapshotResponse]{
					Total:      1,
					IsLoadMore: false,
					Data: []dto.GetDailySnapshotResponse{
						{
							Address: validatorAddress,
							Date:    "2021-01-01",
							Total:   100,
						},
					},
				}).Times(1)

				return fields{
					service: validatorMock,
				}
			},
			args: args{
				w:   sampleResp,
				req: sampleReq,
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			fields: func() fields {
				validatorMock := mocksvc.NewMockValidatorSvc(ctrl)

				validatorMock.EXPECT().GetDailySnapshot(gomock.Any(), dto.GetDailySnapshotRequest{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetDailySnapshotResponse]{
					Total:      0,
					IsLoadMore: false,
					Data:       []dto.GetDailySnapshotResponse{},
				}).Times(0)

				return fields{
					service: validatorMock,
				}
			},
			args: args{
				w:   invalidSampleResp,
				req: invalidSampleReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields()
			i := ValidatorHandlerImpl{
				validatorService: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.GetDailyDelegationSnapshot(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.GetDailyDelegationSnapshot(tt.args.w, tt.args.req)
				})
			}
		})
	}
}

func TestGetDelegatorHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	validatorAddress := "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500"
	delegatorAddress := "cosmos1360qkbsgysnhjeddlwqaqwgj84vq4z8a4g0500"
	page := 1
	limit := 10

	sampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/api/v1/validators/%s/delegator/%s/history?page=%d&limit=%d", validatorAddress, delegatorAddress, page, limit), strings.NewReader(``))
	sampleResp := httptest.NewRecorder()

	invalidSampleReq := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8000/api/v1/validators/%s/delegator/%s/history?page=%d&limit=test", validatorAddress, delegatorAddress, page), strings.NewReader(``))
	invalidSampleResp := httptest.NewRecorder()

	type fields struct {
		service service.ValidatorSvc
	}

	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
	}{
		{
			name: "success get delegator history",
			fields: func() fields {
				validatorMock := mocksvc.NewMockValidatorSvc(ctrl)

				validatorMock.EXPECT().GetDelegatorHistory(gomock.Any(), dto.GetDelegatorHistoryRequest{
					Page:   int32(page),
					Limit:  int32(limit),
					SortBy: "date",
				}).Return(dto.PaginationResp[dto.GetDelegatorHistoryResponse]{
					Total:      1,
					IsLoadMore: false,
					Data: []dto.GetDelegatorHistoryResponse{
						{
							Timestamp: "2021-01-01",
							Amount:    100,
							Change:    100,
						},
					},
				}).Times(1)

				return fields{
					service: validatorMock,
				}
			},
			args: args{
				w:   sampleResp,
				req: sampleReq,
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			fields: func() fields {
				validatorMock := mocksvc.NewMockValidatorSvc(ctrl)

				validatorMock.EXPECT().GetDelegatorHistory(gomock.Any(), dto.GetDelegatorHistoryRequest{
					Page:  int32(page),
					Limit: int32(limit),
				}).Return(dto.PaginationResp[dto.GetDelegatorHistoryResponse]{
					Total:      0,
					IsLoadMore: false,
					Data:       []dto.GetDelegatorHistoryResponse{},
				}).Times(0)

				return fields{
					service: validatorMock,
				}
			},
			args: args{
				w:   invalidSampleResp,
				req: invalidSampleReq,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := tt.fields()
			i := ValidatorHandlerImpl{
				validatorService: field.service,
			}

			if tt.wantErr {
				assert.Panics(t, func() {
					i.GetDelegatorHistory(tt.args.w, tt.args.req)
				})
			} else {
				assert.NotPanics(t, func() {
					i.GetDelegatorHistory(tt.args.w, tt.args.req)
				})
			}
		})
	}
}
