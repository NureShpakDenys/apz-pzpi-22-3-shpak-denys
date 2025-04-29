// Package digcontainer provides a function to build a dig container with all the dependencies
package digcontainer // import "wayra/internal/digcontainer"

import (
	"wayra/internal/adapter/config"
	"wayra/internal/adapter/httpserver"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/adapter/httpserver/handlers/admin"
	"wayra/internal/adapter/httpserver/handlers/company"
	"wayra/internal/adapter/repository"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"
	"wayra/internal/core/service"

	"log/slog"

	"go.uber.org/dig"
)

// BuildContainer creates a dig container with all the dependencies
// Returns a pointer to the container
func BuildContainer() *dig.Container {
	container := dig.New()

	// Config
	container.Provide(config.MustLoad)
	container.Provide(slog.Default)

	// Database
	container.Provide(repository.NewGORMDB)

	// Repositories
	container.Provide(repository.NewRepository[models.Company])
	container.Provide(repository.NewRepository[models.Delivery])
	container.Provide(repository.NewRepository[models.Product])
	container.Provide(repository.NewRepository[models.ProductCategory])
	container.Provide(repository.NewRepository[models.Role])
	container.Provide(repository.NewRepository[models.Route])
	container.Provide(repository.NewRepository[models.SensorData])
	container.Provide(repository.NewRepository[models.User])
	container.Provide(repository.NewRepository[models.Waypoint])
	container.Provide(repository.NewRepository[models.UserCompany])
	container.Provide(repository.NewRepository[models.Log])

	// Services
	container.Provide(service.NewCompanyService)
	container.Provide(service.NewDeliveryService)
	container.Provide(service.NewProductService)
	container.Provide(service.NewProductCategoryService)
	container.Provide(service.NewRouteService)
	container.Provide(service.NewSensorDataService)
	container.Provide(service.NewUserService)
	container.Provide(func(
		userService services.UserService,
		cfg *config.Config,
	) services.AuthService {
		return service.NewAuthService(userService, cfg.AuthConfig.SecretKey, cfg.AuthConfig.TokenExpiry)
	})
	container.Provide(service.NewWaypointService)
	container.Provide(service.NewUserCompanyService)
	container.Provide(service.NewLogService)
	container.Provide(service.NewAdminService)

	// Handlers
	container.Provide(func(authService services.AuthService, cfg *config.Config, logService services.LogService, userService services.UserService) *handlers.AuthHandler {
		return handlers.NewAuthHandler(authService, cfg.AuthConfig.TokenExpiry, logService, userService)
	})

	container.Provide(company.NewCompanyHandler)
	container.Provide(handlers.NewUserHandler)
	container.Provide(handlers.NewRoutesHandler)
	container.Provide(handlers.NewSensorDataHandler)
	container.Provide(handlers.NewWaypointHandler)
	container.Provide(handlers.NewDeliveryHandler)
	container.Provide(handlers.NewProductHandler)
	container.Provide(admin.NewAdminHandler)

	// HTTP Server
	container.Provide(httpserver.NewRouter)

	return container
}
