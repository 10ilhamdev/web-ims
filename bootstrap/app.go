package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"

	"ims/config"
	"ims/routes"
	"ims/app/http/middleware"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithSeeders(Seeders).
		WithMigrations(Migrations).
		WithRouting(func() {
			routes.Web()
			routes.Grpc()
		}).
		WithMiddleware(func(handler configuration.Middleware) {
			handler.Append(
				middleware.StartSession(),
			)
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
