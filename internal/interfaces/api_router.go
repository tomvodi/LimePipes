package interfaces

import "github.com/gin-gonic/gin"

//go:generate mockgen -source api_router.go -destination ./mocks/mock_api_router.go

type ApiRouter interface {
	GetEngine() *gin.Engine
}
