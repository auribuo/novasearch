package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/auribuo/novasearch/validation"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func errorResponse(context *gin.Context, message string) {
	context.AbortWithStatusJSON(400, ErrorResponse{
		Message: message,
	})
}

func idFromPath(context *gin.Context) (int, error) {
	id := context.Param("id")
	if id == "" {
		return -1, fmt.Errorf("missing required param id")
	}
	id = strings.TrimSuffix(id, "/")
	return strconv.Atoi(id)
}

func strategyFromPath(context *gin.Context) (string, error) {
	strategy := context.Param("strategy")
	if strategy == "" {
		return "", fmt.Errorf("missing required param strategy")
	}
	strategy = strings.TrimSuffix(strategy, "/")
	if strategy != "alg" && strategy != "rng" {
		return "", fmt.Errorf("invalid strategy %s", strategy)
	}
	return strategy, nil
}

func requestValid(ctx *gin.Context, request any) bool {
	if err := ctx.ShouldBindJSON(request); err != nil {
		errorResponse(ctx, err.Error())
		return false
	}
	err := validation.FillDefaults(request)
	if err != nil {
		errorResponse(ctx, err.Error())
		return false
	}
	if err := validation.Validate(request); err != nil {
		errorResponse(ctx, err.Error())
		return false
	}
	return true
}
