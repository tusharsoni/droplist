package profile

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSvc,
	NewSQLRepo,
	NewSecrets,

	NewRouter,
	NewGetProfileRoute,
	NewSaveProfileRoute,
)
