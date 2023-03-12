package api

import (
	"time"

	"github.com/auribuo/novasearch/data"
	"github.com/auribuo/novasearch/filter"
	"github.com/auribuo/novasearch/types"
	"github.com/gin-gonic/gin"
)

type GalaxyCalculateRequest types.CalculationData

type GalaxyCalculateResponse struct {
	TotalQuality float64        `json:"totalQuality"`
	Total        int            `json:"total"`
	Time         types.TimeStat `json:"time"`
	Galaxies     []types.Galaxy `json:"galaxies"`
}

func CalculateGalaxyPath(context *gin.Context) {
	strategy, err := strategyFromPath(context)
	if err != nil {
		errorResponse(context, err.Error())
		return
	}

	var request GalaxyCalculateRequest
	if !requestValid(context, &request) {
		return
	}

	galaxies := data.R().GalaxiesFiltered(types.CalculationData(request).BaseData)

	ratables := galaxies.Ratable()

	var filtered []types.Ratable
	var usedTime time.Duration
	if strategy == "rng" {
		filtered, usedTime = filter.Random(ratables, 5*time.Hour)
	} else if strategy == "alg" {
		filtered, usedTime = filter.Generic(ratables, types.CalculationData(request), data.R())
	}

	totalQuality := 0.0
	res := make([]types.Galaxy, 0)

	for _, galaxy := range filtered {
		galaxy.Reset()
		res = append(res, *galaxy.(*types.Galaxy))
		totalQuality += galaxy.Quality()
	}

	context.JSON(200, GalaxyCalculateResponse{
		TotalQuality: totalQuality,
		Total:        len(res),
		Time: types.TimeStat{
			TotalTime: 5 * time.Hour,
			UsedTime:  usedTime,
			Diff:      5*time.Hour - usedTime,
		},
		Galaxies: res,
	})
}

func CalculateViewportPath(context *gin.Context) {
	strategy, err := strategyFromPath(context)
	if err != nil {
		errorResponse(context, err.Error())
		return
	}
	context.JSON(200, gin.H{
		"message": strategy,
	})
}
