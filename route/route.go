package route

import (
	"github.com/gin-gonic/gin"
	"database/sql"

	"question-1/service"
)

type Controller struct {
	service *service.Service
}

func NewController(db *sql.DB) *Controller {
	return &Controller{
		service: service.NewService(db),
	}
}

func (c *Controller) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/users")
	{
		api.POST("/create", c.service.CreateUser)
		api.POST("/generateotp", c.service.GenerateOTP)
		api.POST("/verifyotp", c.service.VerifyOTP)
	}
}