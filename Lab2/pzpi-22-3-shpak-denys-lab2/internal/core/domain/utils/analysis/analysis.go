// Package analysis contains the functions to analyze the data of the deliveries
// and predict the delivery speed
package analysis // import "wayra/internal/core/domain/utils/analysis"

import (
	"fmt"
	"wayra/internal/core/domain/models"
	utilsMath "wayra/internal/core/domain/utils/math"
)

// DeliveryMetrics is a struct that contains the metrics of a delivery
type DeliveryMetrics struct {
	Temperature float64 // Temperature in Celsius
	Humidity    float64 // Humidity in percentage
	WindSpeed   float64 // Wind speed in m/s
	TotalWeight float64 // Total weight in kg

	DeliverySpeed float64 // Delivery speed in km/h
}

// PredictData is a struct that contains the data to predict the delivery speed
type PredictData struct {
	Distance float64 // Distance in km
	Speed    float64 // Speed in km/h
	Time     float64 // Time in hours
}

// PredictedData is a struct that contains the predicted data
// data: the predicted data
// return: the predicted data
func LinearRegression(data []DeliveryMetrics) []float64 {
	n := len(data)

	// X
	temperature := make([]float64, n)
	humidity := make([]float64, n)
	windSpeed := make([]float64, n)
	totalWeight := make([]float64, n)

	// Y
	deliverySpeed := make([]float64, n)

	for i, metrics := range data {
		temperature[i] = metrics.Temperature
		humidity[i] = metrics.Humidity
		windSpeed[i] = metrics.WindSpeed
		totalWeight[i] = metrics.TotalWeight

		deliverySpeed[i] = metrics.DeliverySpeed
	}

	// X^2
	temperatureSquared := utilsMath.Square(temperature)
	humiditySquared := utilsMath.Square(humidity)
	windSpeedSquared := utilsMath.Square(windSpeed)
	totalWeightSquared := utilsMath.Square(totalWeight)
	fmt.Printf("totalWeightSquared: %v\n", totalWeightSquared)

	// Y^2
	deliverySpeedSquared := utilsMath.Square(deliverySpeed)

	// XY
	temperatureDeliverySpeed := utilsMath.Multiply(temperature, deliverySpeed)
	humidityDeliverySpeed := utilsMath.Multiply(humidity, deliverySpeed)
	windSpeedDeliverySpeed := utilsMath.Multiply(windSpeed, deliverySpeed)
	totalWeightDeliverySpeed := utilsMath.Multiply(totalWeight, deliverySpeed)
	fmt.Printf("totalWeightDeliverySpeed: %v\n", totalWeightDeliverySpeed)

	// X sum
	sumTemperature := utilsMath.Sum(temperature)
	sumHumidity := utilsMath.Sum(humidity)
	sumWindSpeed := utilsMath.Sum(windSpeed)
	sumTotalWeight := utilsMath.Sum(totalWeight)
	fmt.Printf("sumTotalWeight: %v\n", sumTotalWeight)

	// Y sum
	sumDeliverySpeed := utilsMath.Sum(deliverySpeed)

	// X^2 sum
	sumTemperatureSquared := utilsMath.Sum(temperatureSquared)
	sumHumiditySquared := utilsMath.Sum(humiditySquared)
	sumWindSpeedSquared := utilsMath.Sum(windSpeedSquared)
	sumTotalWeightSquared := utilsMath.Sum(totalWeightSquared)
	fmt.Printf("sumTotalWeightSquared: %v\n", sumTotalWeightSquared)

	// Y^2 sum
	_ = utilsMath.Sum(deliverySpeedSquared)

	// XY sum
	sumTemperatureDeliverySpeed := utilsMath.Sum(temperatureDeliverySpeed)
	sumHumidityDeliverySpeed := utilsMath.Sum(humidityDeliverySpeed)
	sumWindSpeedDeliverySpeed := utilsMath.Sum(windSpeedDeliverySpeed)
	sumTotalWeightDeliverySpeed := utilsMath.Sum(totalWeightDeliverySpeed)
	fmt.Printf("sumTotalWeightDeliverySpeed: %v\n", sumTotalWeightDeliverySpeed)

	// AvgX
	avgTemperature := utilsMath.Mean(temperature)
	avgHumidity := utilsMath.Mean(humidity)
	avgWindSpeed := utilsMath.Mean(windSpeed)
	avgTotalWeight := utilsMath.Mean(totalWeight)
	fmt.Printf("avgTotalWeight: %v\n", avgTotalWeight)

	// AvgY
	avgDeliverySpeed := utilsMath.Mean(deliverySpeed)

	// Betas
	betaTemperature := CalculateBeta(sumTemperature, sumDeliverySpeed, sumTemperatureDeliverySpeed, sumTemperatureSquared, n)
	betaHumidity := CalculateBeta(sumHumidity, sumDeliverySpeed, sumHumidityDeliverySpeed, sumHumiditySquared, n)
	betaWindSpeed := CalculateBeta(sumWindSpeed, sumDeliverySpeed, sumWindSpeedDeliverySpeed, sumWindSpeedSquared, n)
	betaTotalWeight := CalculateBeta(sumTotalWeight, sumDeliverySpeed, sumTotalWeightDeliverySpeed, sumTotalWeightSquared, n)

	fmt.Printf("Beta Total Weight: %f\n", betaTotalWeight)

	// Free term
	beta0 := avgDeliverySpeed - betaTemperature*avgTemperature - betaHumidity*avgHumidity - betaWindSpeed*avgWindSpeed - betaTotalWeight*avgTotalWeight

	// Calculate errors
	temperatureErrors := make([]float64, n)
	humidityErrors := make([]float64, n)
	windSpeedErrors := make([]float64, n)
	totalWeightErrors := make([]float64, n)

	for i := 0; i < n; i++ {
		temperatureErrors[i] = deliverySpeed[i] - beta0 - betaTemperature*temperature[i] - betaHumidity*humidity[i] - betaWindSpeed*windSpeed[i]
		humidityErrors[i] = deliverySpeed[i] - beta0 - betaTemperature*temperature[i] - betaHumidity*humidity[i] - betaWindSpeed*windSpeed[i]
		windSpeedErrors[i] = deliverySpeed[i] - beta0 - betaTemperature*temperature[i] - betaHumidity*humidity[i] - betaWindSpeed*windSpeed[i]
		totalWeightErrors[i] = deliverySpeed[i] - beta0 - betaTemperature*temperature[i] - betaHumidity*humidity[i] - betaWindSpeed*windSpeed[i]
	}

	// Calculate average error
	avgTemperatureError := utilsMath.Mean(temperatureErrors)
	avgHumidityError := utilsMath.Mean(humidityErrors)
	avgWindSpeedError := utilsMath.Mean(windSpeedErrors)
	avgTotalWeightError := utilsMath.Mean(totalWeightErrors)

	error := avgTemperatureError + avgHumidityError + avgWindSpeedError + avgTotalWeightError

	return []float64{beta0, betaTemperature, betaHumidity, betaWindSpeed, betaTotalWeight, error}
}

// CalculateBeta calculates the beta value of regression
// sumX: sum of the X values
// sumY: sum of the Y values
// sumXY: sum of the X*Y values
// sumXSquared: sum of the X^2 values
// n: number of elements
// return: the beta value
func CalculateBeta(sumX, sumY, sumXY, sumXSquared float64, n int) float64 {
	if sumX == 0 && sumXSquared == 0 {
		return 0
	}

	return (sumX*sumY - float64(n)*sumXY) / (sumX*sumX - float64(n)*sumXSquared)
}

// Predict predicts the delivery speed
// coeffs: the coefficients of the regression
// input: the input data
// totalWeight: the total weight of the delivery
// return: the predicted delivery speed
func Predict(coeffs []float64, input models.SensorData, totalWeight float64) float64 {
	features := []float64{
		input.Temperature,
		input.Humidity,
		input.WindSpeed,
		totalWeight,
	}

	result := coeffs[0]
	for i := 0; i < len(features); i++ {
		result += coeffs[i+1] * features[i]
	}

	result += coeffs[len(coeffs)-1]

	return result
}
