package service // import "wayra/internal/core/service"

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/domain/utils/analysis"
	utilsMath "wayra/internal/core/domain/utils/math"
	utilsTime "wayra/internal/core/domain/utils/time"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"

	"log/slog"
)

// RouteService is a struct that defines the service for the Route model
type RouteService struct {
	*GenericService[models.Route]                                    // Embedding the GenericService struct for the Route model
	waypointRepository            port.Repository[models.Waypoint]   // Repository for the Waypoint model
	deliveryRepository            port.Repository[models.Delivery]   // Repository for the Delivery model
	sensorDataRepository          port.Repository[models.SensorData] // Repository for the SensorData model
}

// NewRouteService is a function that creates a new RouteService instance
// repo: Repository for the Route model
// waypointRepository: Repository for the Waypoint model
// deliveryRepository: Repository for the Delivery model
// sensorDataRepository: Repository for the SensorData model
// Returns a pointer to the RouteService instance
func NewRouteService(
	repo port.Repository[models.Route],
	waypointRepository port.Repository[models.Waypoint],
	deliveryRepository port.Repository[models.Delivery],
	sensorDataRepository port.Repository[models.SensorData],
) services.RouteService {
	return &RouteService{
		GenericService:       NewGenericService(repo),
		waypointRepository:   waypointRepository,
		deliveryRepository:   deliveryRepository,
		sensorDataRepository: sensorDataRepository,
	}
}

// GetOptimalRoute is a function that returns the optimal route for a delivery
// ctx: Context for the request
// delivery: Delivery for which the optimal route is to be found
// includeWeight: Boolean to include weight in the calculation
// considerPerishable: Boolean to consider perishable products in the calculation
// Returns the additional message, predict data, coefficients, optimal route, and error
func (s *RouteService) GetOptimalRoute(
	ctx context.Context,
	delivery *models.Delivery,
	includeWeight bool,
	considerPerishable bool,
) (string, *analysis.PredictData, []float64, models.Route, error) {
	var predictData analysis.PredictData
	var coeffs []float64
	var optimalRoute *models.Route
	var additionalMessage string

	routes, err := s.Where(ctx, &models.Route{CompanyID: delivery.CompanyID})
	if err != nil || len(routes) == 0 {
		return "", nil, nil, models.Route{}, errors.New("no routes found for the company")
	}

	minDeliveryTime := math.MaxFloat64
	minHumidity := 85.0
	maxTemperature := 0.0

	var deliveryMetrics []analysis.DeliveryMetrics
	for _, route := range routes {
		waypoints, err := s.waypointRepository.Where(ctx, &models.Waypoint{RouteID: route.ID})
		if err != nil {
			return "", nil, nil, models.Route{}, err
		}

		deliveries, err := s.deliveryRepository.Where(ctx, &models.Delivery{
			Status:  "completed",
			RouteID: route.ID,
		})
		if err != nil {
			return "", nil, nil, models.Route{}, err
		}

		for _, delivery := range deliveries {
			data := CalculateRouteMetrics(delivery, waypoints, includeWeight)
			if data == nil {
				continue
			}

			deliveryMetrics = append(deliveryMetrics, *data)
		}
	}

	coeffs = analysis.LinearRegression(deliveryMetrics)

	for _, route := range routes {
		waypoints, err := s.waypointRepository.Where(ctx, &models.Waypoint{RouteID: route.ID})
		if err != nil {
			return "", nil, nil, models.Route{}, err
		}

		latestSensorData := []models.SensorData{}
		for _, waypoint := range waypoints {
			latestSensorData = append(
				latestSensorData,
				waypoint.SensorData[len(waypoint.SensorData)-1],
			)
		}

		avgTemp := 0.0
		avgHumidity := 0.0
		avgWindSpeed := 0.0
		avgPressure := 0.0

		for _, sensorData := range latestSensorData {
			avgTemp += sensorData.Temperature
			avgHumidity += sensorData.Humidity
			avgWindSpeed += sensorData.WindSpeed
			avgPressure += sensorData.MeanPressure
		}
		avgData := models.SensorData{
			Temperature:  avgTemp / float64(len(latestSensorData)),
			Humidity:     avgHumidity / float64(len(latestSensorData)),
			WindSpeed:    avgWindSpeed / float64(len(latestSensorData)),
			MeanPressure: avgPressure / float64(len(latestSensorData)),
		}

		totalWeight := 0.0
		isPerishable := false

		for _, product := range delivery.Products {
			if includeWeight {
				totalWeight += product.Weight
			}
			if product.ProductCategory.IsPerishable {
				isPerishable = true
			}
		}

		predictedSpeed := analysis.Predict(coeffs, avgData, totalWeight)

		var distance float64

		for i := 0; i < len(waypoints)-1; i++ {
			distance += utilsMath.HaversineDistance(
				waypoints[i].Latitude,
				waypoints[i].Longitude,
				waypoints[i+1].Latitude,
				waypoints[i+1].Longitude,
			)
		}

		time := distance / predictedSpeed

		if isPerishable && considerPerishable {
			additionalMessage = "Recommended route depends on perishable products"
			if time < minDeliveryTime {
				minDeliveryTime = time
				optimalRoute = &route

				predictData = analysis.PredictData{
					Distance: distance,
					Speed:    predictedSpeed,
					Time:     time,
				}
			}
		} else {
			if avgData.Humidity < minHumidity && avgData.Temperature > maxTemperature {
				additionalMessage = "Recommended route based on safety conditions"
				if time < minDeliveryTime {
					minDeliveryTime = time
					optimalRoute = &route

					predictData = analysis.PredictData{
						Distance: distance,
						Speed:    predictedSpeed,
						Time:     time,
					}
				}
			} else {
				additionalMessage = "Recommended route depends on speed of the route"
				if time < minDeliveryTime {
					minDeliveryTime = time
					optimalRoute = &route

					predictData = analysis.PredictData{
						Distance: distance,
						Speed:    predictedSpeed,
						Time:     time,
					}
				}
			}
		}
	}

	return additionalMessage, &predictData, coeffs, *optimalRoute, nil
}

// CalculateRouteMetrics is a function that calculates the metrics for a delivery route
// delivery: Delivery for which the metrics are to be calculated
// waypoints: Waypoints for the delivery route
// includeWeight: Boolean to include weight in the calculation
// Returns the calculated metrics for the delivery route
func CalculateRouteMetrics(delivery models.Delivery, waypoints []models.Waypoint, includeWeight bool) *analysis.DeliveryMetrics {
	speedData := analysis.DeliveryMetrics{}
	totalDistance := 0.0
	for i := 0; i < len(waypoints)-1; i++ {
		totalDistance += utilsMath.HaversineDistance(
			waypoints[i].Latitude,
			waypoints[i].Longitude,
			waypoints[i+1].Latitude,
			waypoints[i+1].Longitude,
		)
	}

	sensorData := []models.SensorData{}

	for _, waypoint := range waypoints {
		sensorData = append(sensorData, waypoint.SensorData...)
	}

	tempSum := 0.0
	humiditySum := 0.0
	windSpeedSum := 0.0
	totalWeight := 0.0
	count := 0

	for _, product := range delivery.Products {
		totalWeight += product.Weight
	}

	duration, err := utilsTime.ParseDuration(delivery.Duration)
	if err != nil {
		slog.Error("Error parsing, duration: ", slog.Any("error", err.Error()))
		return nil
	}

	for _, sensorData := range sensorData {
		oneHourBefore := delivery.Date.Add(-1 * time.Hour)
		oneHourAfter := delivery.Date.Add(time.Hour)

		if sensorData.Date.After(oneHourBefore) && sensorData.Date.Before(oneHourAfter.Add(duration)) {
			tempSum += sensorData.Temperature
			humiditySum += sensorData.Humidity
			windSpeedSum += sensorData.WindSpeed
			count++
		}
	}

	speedData = analysis.DeliveryMetrics{
		Temperature:   tempSum / float64(count),
		Humidity:      humiditySum / float64(count),
		WindSpeed:     windSpeedSum / float64(count),
		DeliverySpeed: totalDistance / duration.Hours(),
	}

	if includeWeight {
		speedData.TotalWeight = totalWeight
	} else {
		speedData.TotalWeight = 0.0
	}

	return &speedData
}

// GetWeatherAlert is a function that returns the weather alerts for a route
// ctx: Context for the request
// route: Route for which the weather alerts are to be found
// Returns the weather alerts and error
func (s *RouteService) GetWeatherAlert(
	ctx context.Context,
	route models.Route,
) ([]models.WeatherAlert, error) {
	alerts := []models.WeatherAlert{}

	if len(route.Waypoints) == 0 {
		return nil, errors.New("no waypoints found for the route")
	}

	latestSensorData := []models.SensorData{}
	for _, waypoint := range route.Waypoints {
		latestSensorData = append(
			latestSensorData,
			waypoint.SensorData[len(waypoint.SensorData)-1],
		)
	}

	if len(latestSensorData) == 0 {
		return nil, errors.New("no sensor data available for the route")
	}

	existingAlertTypes := make(map[string]bool)

	for _, data := range latestSensorData {
		if data.Temperature < 0 && data.Humidity > 80 {
			if !existingAlertTypes["Ice Alert"] {
				alerts = append(alerts, models.WeatherAlert{
					Type:    "Ice Alert",
					Message: "Potential ice formation detected due to low temperature and high humidity.",
					Details: fmt.Sprintf("Temperature: %.2f°C, Humidity: %.2f%%", data.Temperature, data.Humidity),
				})
				existingAlertTypes["Ice Alert"] = true
			}
		}

		if data.WindSpeed > 20 {
			if !existingAlertTypes["Storm Alert"] {
				alerts = append(alerts, models.WeatherAlert{
					Type:    "Storm Alert",
					Message: "High wind speed detected, potential storm risk.",
					Details: fmt.Sprintf("Wind Speed: %.2f m/s", data.WindSpeed),
				})
				existingAlertTypes["Storm Alert"] = true
			}
		}

		if data.MeanPressure < 980 {
			if !existingAlertTypes["Low Pressure Alert"] {
				alerts = append(alerts, models.WeatherAlert{
					Type:    "Low Pressure Alert",
					Message: "Low atmospheric pressure detected, potential severe weather conditions.",
					Details: fmt.Sprintf("Pressure: %.2f hPa", data.MeanPressure),
				})
				existingAlertTypes["Low Pressure Alert"] = true
			}
		}

		if data.Temperature > 35 {
			if !existingAlertTypes["Heat Alert"] {
				alerts = append(alerts, models.WeatherAlert{
					Type:    "Heat Alert",
					Message: "High temperature detected, risk of heat-related issues.",
					Details: fmt.Sprintf("Temperature: %.2f°C", data.Temperature),
				})
				existingAlertTypes["Heat Alert"] = true
			}
		}

		if data.Humidity < 20 {
			if !existingAlertTypes["Low Humidity Alert"] {
				alerts = append(alerts, models.WeatherAlert{
					Type:    "Low Humidity Alert",
					Message: "Low humidity detected, risk of dry conditions.",
					Details: fmt.Sprintf("Humidity: %.2f%%", data.Humidity),
				})
				existingAlertTypes["Low Humidity Alert"] = true
			}
		}

		if data.WindSpeed > 30 && data.Temperature < 5 {
			if !existingAlertTypes["Cold Storm Alert"] {
				alerts = append(alerts, models.WeatherAlert{
					Type:    "Cold Storm Alert",
					Message: "High wind speed combined with low temperature detected, risk of severe cold storm.",
					Details: fmt.Sprintf("Wind Speed: %.2f m/s, Temperature: %.2f°C", data.WindSpeed, data.Temperature),
				})
				existingAlertTypes["Cold Storm Alert"] = true
			}
		}
	}

	return alerts, nil
}
