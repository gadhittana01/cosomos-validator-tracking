package dto

type GetHourlySnapshotResponse struct {
	Address   string `json:"address"`
	Amount    int64  `json:"amount"`
	Change    int64  `json:"change"`
	Date      string `json:"date"`
	Timestamp string `json:"timestamp"`
}

type GetDailySnapshotResponse struct {
	Address string `json:"address"`
	Date    string `json:"date"`
	Total   int64  `json:"total"`
}

type GetDelegatorHistoryResponse struct {
	Timestamp string `json:"timestamp"`
	Amount    int64  `json:"amount"`
	Change    int64  `json:"change"`
}
