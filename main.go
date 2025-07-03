package main

import (
	"embed"

	"github.com/CrowderSoup/emoji-match/config"
	"github.com/CrowderSoup/emoji-match/server"
	"github.com/CrowderSoup/emoji-match/utils"
	"go.uber.org/fx"
)

//go:embed all:dist
var staticAssets embed.FS

func main() {
	bundle := fx.Options(
		config.Module,
		utils.Module,
		server.Module,
		fx.Provide(func() embed.FS {
			return staticAssets
		}),
	)

	app := fx.New(
		bundle,
		fx.Invoke(utils.InvokeLogger, utils.InvokeTracerProvider),
		fx.Invoke(server.Run),
	)

	app.Run()

	<-app.Done()
}
