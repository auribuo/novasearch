package types

type CalculationData struct {
	BaseData
	GridSize        float64 `json:"gridSize" validate:"required"`
	SearchRadius    float64 `json:"searchRadius" validate:"required"`
	MaxSearchTime   int     `json:"maxSearchTime" validate:"required"`
	StartPointCount int     `json:"startPointCount" validate:"required"`
}
