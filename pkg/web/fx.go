package web

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewRouter,
	NewAppRoute,
	NewStaticRoute,
)
