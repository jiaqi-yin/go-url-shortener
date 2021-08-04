package main

import "github.com/jiaqi-yin/go-url-shortener/app"

func main() {
	app := app.App{}
	app.Init()
	app.Run()
}
