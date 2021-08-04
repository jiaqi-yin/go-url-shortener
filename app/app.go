package app

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jiaqi-yin/go-url-shortener/config"
	"github.com/jiaqi-yin/go-url-shortener/controllers/ping"
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

type EncodedID struct {
	Shortlink string `uri:"shortlink" binding:"required"`
}

func (app *App) initializeRoutes() {
	app.Router.GET("/ping", ping.PingController.Ping)
	app.Router.POST("/api/shorten", app.createShortlink)
	app.Router.GET("/:shortlink", app.unshorten)
}

func (app *App) Run(addr string) {
	log.Fatal(app.Router.Run(addr))
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

	eid, err := app.Config.S.Shorten(req.URL)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, shortlink.ShortlinkResponse{Shortlink: eid})
	return
}

func (app *App) unshorten(c *gin.Context) {
	var eid EncodedID
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

	url, err := app.Config.S.Unshorten(eid.Shortlink)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
	return
}
