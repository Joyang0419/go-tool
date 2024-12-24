package main

import (
	"time"

	"go.uber.org/fx"

	"go-tool/web"
)

func main() {
	fx.New(
		fx.Supply(web.NewServerConfig(8081, 5*time.Second)),
		web.RegisterController(web.NewHealthController),
		fx.Invoke(web.NewServer),
	).Run()
}
