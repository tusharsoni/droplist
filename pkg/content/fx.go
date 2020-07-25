package content

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSvc,
	NewSQLRepo,

	NewRouter,
	NewGetTemplateRoute,
	NewListTemplatesRoute,
	NewDeleteTemplateRoute,
	NewCreateTemplateRoute,
	NewUpdateTemplateRoute,
	NewPreviewTemplateHTMLRoute,
)
