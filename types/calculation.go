package types

type CalculationData struct {
	BaseData
	CalculateViewports bool    `json:"calculateViewports"`
	GridSize           float64 `json:"gridSize"`
	SearchRadius       float64 `json:"searchRadius"`
	MaxSearchTime      float64 `json:"maxSearchTime"`
	StartPointCount    int     `json:"startPointCount"`
}
