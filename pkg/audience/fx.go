package audience

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSvc,
	NewSQLRepo,

	NewRouter,
	NewCreateListRoute,
	NewListContactsRoute,
	NewCreateContactsRoute,
	NewAddContactsToListRoute,
)
