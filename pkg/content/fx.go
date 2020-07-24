package content

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSvc,
	NewSQLRepo,

	NewRouter,
	NewListTemplatesRoute,
	NewCreateTemplateRoute,
)
