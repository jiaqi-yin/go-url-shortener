package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jiaqi-yin/go-url-shortener/controllers/ping"
)

type App struct {
	Router *gin.Engine
}

func (app *App) Init() {
	app.Router = gin.Default()
	// app.Router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 	// custom format
	// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.ClientIP,
	// 		param.TimeStamp.Format(time.RFC1123),
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 		param.ErrorMessage,
	// 	)
	// }))
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.GET("/ping", ping.PingController.Ping)
}

func (app *App) Run(addr string) {
	log.Fatal(app.Router.Run(addr))
}
