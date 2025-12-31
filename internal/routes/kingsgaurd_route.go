package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func KGRouters(rg *gin.RouterGroup, provider *di.Provider) {
	kg := rg.Group("/kingsguard")
	{
		kg.GET("/polices", provider.KingsGuard.Controller.List)
		kg.POST("/polices", provider.KingsGuard.Controller.Create)
		kg.DELETE("/polices/:id", provider.KingsGuard.Controller.Delete)
	}
}
