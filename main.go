package main

import (
	"github.com/CrowderSoup/emoji-match/config"
	"github.com/CrowderSoup/emoji-match/server"
	"github.com/CrowderSoup/emoji-match/utils"
	"go.uber.org/fx"
)

func main() {
	bundle := fx.Options(
		config.Module,
		utils.Module,
		server.Module,
	)

	app := fx.New(
		bundle,
		fx.Invoke(utils.InvokeLogger, utils.InvokeTracerProvider),
		fx.Invoke(server.Run),
	)

	app.Run()

	<-app.Done()
}
