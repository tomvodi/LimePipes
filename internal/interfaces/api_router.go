package interfaces

import "github.com/gin-gonic/gin"

type ApiRouter interface {
	GetEngine() *gin.Engine
}
