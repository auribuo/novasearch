package routes

import (
	"github.com/auribuo/novasearch/data"
	"github.com/auribuo/novasearch/types"
	"github.com/gin-gonic/gin"
)

type GalaxyResponse struct {
	Total    int            `json:"total"`
	Galaxies []types.Galaxy `json:"galaxies"`
}

type GalaxyFilterRequest types.BaseData

func AllGalaxies(context *gin.Context) {
	galaxies := data.Galaxies()
	context.JSON(200, GalaxyResponse{
		Total:    len(galaxies),
		Galaxies: galaxies,
	})
}

func Galaxy(context *gin.Context) {
	id, err := idFromPath(context)
	if err != nil {
		context.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	galaxy := data.Galaxy(id)
	context.JSON(200, GalaxyResponse{
		Total:    1,
		Galaxies: []types.Galaxy{galaxy},
	})
}

func FilterGalaxies(context *gin.Context) {
	var request GalaxyFilterRequest
	context.BindJSON(&request)

	galaxies := data.GalaxiesFiltered(types.BaseData(request))

	context.JSON(200, GalaxyResponse{
		Total:    len(galaxies),
		Galaxies: galaxies,
	})
}
