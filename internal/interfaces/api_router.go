package interfaces

import "github.com/gin-gonic/gin"

type APIRouter interface {
	GetEngine() *gin.Engine
}
