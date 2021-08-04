package shortener

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ShortenerController shortenerControllerInterface = &shortenerController{}
)

type shortenerControllerInterface interface {
	CreateShortlink(*gin.Context)
	Unshorten(*gin.Context)
}

type shortenerController struct{}

func (controller *shortenerController) CreateShortlink(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (controller *shortenerController) Unshorten(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
