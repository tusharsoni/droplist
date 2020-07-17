package main

import (
	"shoot/pkg/audience"

	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/csql"
	"go.uber.org/fx"
)

func main() {
	app := copper.NewHTTPApp(
		clogger.StdFx,
		csql.Fx,

		fx.Provide(NewConfig),

		audience.Fx,

		fx.Invoke(
			audience.RunMigrations,
		),
	)

	app.Run()
}
