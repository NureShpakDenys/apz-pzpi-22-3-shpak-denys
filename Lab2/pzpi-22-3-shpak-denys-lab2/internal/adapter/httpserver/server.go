// Package httpserver provides the http server and the routes for the API.
package httpserver // import "wayra/internal/adapter/httpserver"

import (
	"log/slog"
	"wayra/internal/adapter/config"
	"wayra/internal/adapter/httpserver/handlers"
	"wayra/internal/adapter/httpserver/handlers/admin"
	"wayra/internal/adapter/httpserver/handlers/company"
	"wayra/internal/adapter/httpserver/middlewares"
	"wayra/internal/core/port/services"

	_ "wayra/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates a new gin router with all the routes and middlewares
// log: logger
// cfg: config
// authHandler: handler for the auth routes
// companyHandler: handler for the company routes
// userHandler: handler for the user routes
// routeHanler: handler for the route routes
// waypointHandler: handler for the waypoint routes
// authService: service to validate the token
// sensorDataHandler: handler for the sensor data routes
// deliveryHandler: handler for the delivery routes
// productHandler: handler for the product routes
// adminHandler: handler for the admin routes
// returns: *gin.Engine
func NewRouter(
	log *slog.Logger,
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	companyHandler *company.CompanyHandler,
	userHandler *handlers.UserHandler,
	routeHanler *handlers.RouteHandler,
	waypointHandler *handlers.WaypointHandler,
	authService services.AuthService,
	sensorDataHandler *handlers.SensorDataHandler,
	deliveryHandler *handlers.DeliveryHandler,
	productHandler *handlers.ProductHandler,
	adminHandler *admin.AdminHandler,
	userCompanyHandler *handlers.UserCompanyHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.DocExpansion("false"),
		ginSwagger.DefaultModelsExpandDepth(0),
		ginSwagger.InstanceName("swagger"),
	))

	r.GET("device-config/:waypoint_id", waypointHandler.GetDeviceConfig)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.RegisterUser)
		auth.POST("/login", authHandler.LoginUser)
	}

	sensorData := r.Group("/sensor-data")
	{
		sensorData.POST("/", sensorDataHandler.AddSensorData)
		sensorData.GET("/:sensor_data_id", sensorDataHandler.GetSensorData)
		sensorData.PUT("/:sensor_data_id", sensorDataHandler.UpdateSensorData)
		sensorData.DELETE("/:sensor_data_id", sensorDataHandler.DeleteSensorData)
	}

	r.Use(middlewares.AuthMiddleware(log, authService))

	routes := r.Group("/routes")
	{
		routes.POST("/", routeHanler.CreateRoute)
		routes.GET("/:route_id", routeHanler.GetRoute)
		routes.PUT("/:route_id", routeHanler.UpdateRoute)
		routes.DELETE("/:route_id", routeHanler.DeleteRoute)
		routes.GET("/:route_id/get-sensor-data", routeHanler.GetSensorData)

		routes.GET("/:route_id/weather-alert", routeHanler.GetWeatherAlert)
	}

	waypoints := r.Group("/waypoints")
	{
		waypoints.POST("/", waypointHandler.AddWaypoint)
		waypoints.GET("/:waypoint_id", waypointHandler.GetWaypoint)
		waypoints.PUT("/:waypoint_id", waypointHandler.UpdateWaypoint)
		waypoints.DELETE("/:waypoint_id", waypointHandler.DeleteWaypoint)
		waypoints.GET("/", waypointHandler.GetWaypoints)
	}

	r.POST("/auth/logout", authHandler.LogoutUser)

	user := r.Group("/user")
	{
		user.GET("/:id", userHandler.GetUser)
		user.PUT("/:id", userHandler.UpdateUser)
		user.DELETE("/:id", userHandler.DeleteUser)
	}

	r.GET("/users", userHandler.GetUsers)

	company := r.Group("/company")
	{
		company.POST("/", companyHandler.RegisterCompany)
		company.GET("/:company_id", companyHandler.GetCompany)
		company.GET("/:company_id/users", userCompanyHandler.GetCompanyUsers)
		company.GET("/", companyHandler.GetCompanies)
		company.PUT("/:company_id", companyHandler.UpdateCompany)
		company.DELETE("/:company_id", companyHandler.DeleteCompany)

		company.POST("/:company_id/add-user", companyHandler.AddUserToCompany)
		company.PUT("/:company_id/update-user", companyHandler.UpdateUserInCompany)
		company.DELETE("/:company_id/remove-user", companyHandler.RemoveUserFromCompany)
	}

	deliveries := r.Group("/delivery")
	{
		deliveries.POST("/", deliveryHandler.CreateDelivery)
		deliveries.GET("/:delivery_id", deliveryHandler.GetDelivery)
		deliveries.PUT("/:delivery_id", deliveryHandler.UpdateDelivery)
		deliveries.DELETE("/:delivery_id", deliveryHandler.DeleteDelivery)
	}

	products := r.Group("/products")
	{
		products.POST("/", productHandler.AddProduct)
		products.GET("/:product_id", productHandler.GetProduct)
		products.PUT("/:product_id", productHandler.UpdateProduct)
		products.DELETE("/:product_id", productHandler.DeleteProduct)
	}

	analytics := r.Group("/analytics")
	{
		analytics.GET("/:delivery_id/optimal-route", routeHanler.GetOptimalRoute)
		analytics.GET("/:delivery_id/optimal-back-route", routeHanler.GetOptimalBackRoute)
	}

	admin := r.Group("/admin")
	{
		admin.POST("/backup", adminHandler.BackupDatabase)
		admin.POST("/restore", adminHandler.RestoreDatabase)
		admin.POST("/change-role", adminHandler.UpdateUserRole)
		admin.GET("/health", adminHandler.HealthCheckHandler)
		admin.GET("/logs", adminHandler.GetServerLogs)
		admin.GET("/system-configs", adminHandler.GetSystemConfigs)
		admin.PUT("/system-configs", adminHandler.UpdateSystemConfigs)
		admin.GET("/db-status", adminHandler.GetDBStatus)
		admin.POST("/optimize", adminHandler.OptimizeDatabase)
		admin.POST("/clear-logs", adminHandler.ClearLogs)
		admin.POST("/send-config", adminHandler.SendConfigToGlobalStorage)
	}

	return r
}
