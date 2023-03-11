package routes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func idFromPath(context *gin.Context) (int, error) {
	id := context.Param("id")
	if id == "" {
		return -1, fmt.Errorf("missing required param id")
	}
	id = strings.TrimSuffix(id, "/")
	return strconv.Atoi(id)
}
