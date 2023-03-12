package api

import (
	"math"
	"strconv"

	"github.com/auribuo/novasearch/data"
	"github.com/auribuo/novasearch/types"
	"github.com/gin-gonic/gin"
)

type GalaxyResponse struct {
	Total     int              `json:"total"`
	Galaxies  []types.Galaxy   `json:"galaxies"`
	Viewports []types.Viewport `json:"viewports"`
}

type GalaxyFilterRequest types.BaseData

// AllGalaxies godoc
// @Summary Get all galaxies.
// @Description Get all galaxies that were fetched from the database.
// @ID galaxies
// @Tags galaxies
// @param limit query int false "Limit the number of galaxies returned."
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} GalaxyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /galaxies [get]
func AllGalaxies(context *gin.Context) {
	limitString := context.Query("limit")
	if limitString == "" {
		limitString = "-1"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		errorResponse(context, "invalid limit")
		return
	}

	galaxies := data.R().Galaxies()

	if limit > 0 {
		limit = int(math.Min(float64(len(galaxies)), float64(limit)))
		galaxies = galaxies[:limit]
	}

	context.JSON(200, GalaxyResponse{
		Total:    len(galaxies),
		Galaxies: galaxies,
	})
}

// Galaxy godoc
// @Summary Get a galaxy.
// @Description Get a galaxy by its ID.
// @ID galaxy
// @Tags galaxies
// @Accept  application/json
// @Produce  application/json
// @Param id path int true "Galaxy ID"
// @Success 200 {object} GalaxyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /galaxies/{id} [get]
func Galaxy(context *gin.Context) {
	id, err := idFromPath(context)
	if err != nil {
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	galaxy := data.R().Galaxy(id)
	context.JSON(200, GalaxyResponse{
		Total:    1,
		Galaxies: []types.Galaxy{galaxy},
	})
}

// FilterGalaxies godoc
// @Summary Filter galaxies.
// @Description Filter galaxies by a set of parameters. The filter used is only the situational. No algorithm is used.
// @ID filter-galaxies
// @Tags galaxies
// @param limit query int false "Limit the number of galaxies returned."
// @Accept  application/json
// @Produce  application/json
// @Param request body GalaxyFilterRequest true "Galaxy filter request"
// @Success 200 {object} GalaxyResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /galaxies [post]
func FilterGalaxies(context *gin.Context) {
	limitString := context.Query("limit")
	if limitString == "" {
		limitString = "-1"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		errorResponse(context, "invalid limit")
		return
	}

	var request GalaxyFilterRequest
	if !requestValid(context, &request) {
		return
	}

	galaxies := data.R().GalaxiesFiltered(types.BaseData(request))

	if limit > 0 {
		limit = int(math.Min(float64(len(galaxies)), float64(limit)))
		galaxies = galaxies[:limit]
	}

	viewports := data.R().ViewportsFiltered(galaxies, types.BaseData(request))

	context.JSON(200, GalaxyResponse{
		Total:     len(galaxies),
		Galaxies:  galaxies,
		Viewports: viewports,
	})
}
