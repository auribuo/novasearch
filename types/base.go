package types

type BaseData struct {
	StartTime        IsoTime    `json:"startTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Location         Location   `json:"location" validate:"required"`
	Fov              Fov        `json:"fov" validate:"required"`
	Hemisphere       Hemisphere `json:"hemisphere" validate:"required,oneof=W E"`
	MinHeight        float64    `json:"minHeight" default:"30"`
	MaxSemiMajorAxis float64    `json:"maxSemiMajorAxis" default:"10"`
	MaxSemiMinorAxis float64    `json:"maxSemiMinorAxis" default:"10"`
	IncludeViewports bool       `json:"includeViewports" default:"false"`
}
