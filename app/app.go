package app

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jiaqi-yin/go-url-shortener/config"
	"github.com/jiaqi-yin/go-url-shortener/domain/shortlink"
	"github.com/jiaqi-yin/go-url-shortener/utils"
)

type App struct {
	Router *gin.Engine
	Config *config.Config
}

func (app *App) Init() {
	app.Config = config.LoadConfig()
	app.Router = gin.Default()
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.GET("/ping", app.ping)
	app.Router.POST("/api/shorten", app.createShortlink)
	app.Router.GET("/:shortlink", app.redirect)
}

func (app *App) Run() {
	log.Fatal(app.Router.Run(app.Config.ServerAddr))
}

func (app *App) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
	return
}

func (app *App) createShortlink(c *gin.Context) {
	var req shortlink.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restErr := utils.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	eid, err := app.Config.ShortlinkService.Shorten(req.URL)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, shortlink.ShortlinkResponse{Shortlink: eid})
	return
}

func (app *App) redirect(c *gin.Context) {
	var eid shortlink.EncodedID
	if err := c.ShouldBindUri(&eid); err != nil {
		restErr := utils.NewBadRequestError("cannot bind uri")
		c.JSON(restErr.Status(), restErr)
		return
	}

	var validShortlink = regexp.MustCompile(`^[a-zA-Z0-9]{1,11}$`)
	if !validShortlink.MatchString(eid.Shortlink) {
		restErr := utils.NewBadRequestError("invalid uri")
		c.JSON(restErr.Status(), restErr)
		return
	}

	url, err := app.Config.ShortlinkService.Unshorten(eid.Shortlink)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
	return
}
