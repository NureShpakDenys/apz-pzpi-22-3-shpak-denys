// Package math provides mathematical functions for the domain layer.
// This package is used to perform mathematical operations on data.
package math // import "wayra/internal/core/domain/utils/math"

import "math"

// Transpose returns the transpose of a matrix.
// matrix: a 2D slice of float64.
// returns: a 2D slice of float64.
func Transpose(matrix [][]float64) [][]float64 {
	rows := len(matrix)
	cols := len(matrix[0])
	result := make([][]float64, cols)
	for i := range result {
		result[i] = make([]float64, rows)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result[j][i] = matrix[i][j]
		}
	}
	return result
}

// MultiplyMatrices returns the product of two matrices.
// a: a 2D slice of float64.
// b: a 2D slice of float64.
// returns: a 2D slice of float64.
func MultiplyMatrices(a, b [][]float64) [][]float64 {
	rowsA := len(a)
	colsA := len(a[0])
	colsB := len(b[0])
	result := make([][]float64, rowsA)
	for i := range result {
		result[i] = make([]float64, colsB)
	}
	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

// Inverse returns the inverse of a matrix.
// matrix: a 2D slice of float64.
// returns: a 2D slice of float64.
func Inverse(matrix [][]float64) [][]float64 {

	n := len(matrix)
	augmented := make([][]float64, n)
	for i := range augmented {
		augmented[i] = make([]float64, 2*n)
		copy(augmented[i], matrix[i])
		augmented[i][n+i] = 1
	}

	for i := 0; i < n; i++ {

		maxRow := i
		for j := i + 1; j < n; j++ {
			if math.Abs(augmented[j][i]) > math.Abs(augmented[maxRow][i]) {
				maxRow = j
			}
		}

		augmented[i], augmented[maxRow] = augmented[maxRow], augmented[i]

		pivot := augmented[i][i]
		for j := 0; j < 2*n; j++ {
			augmented[i][j] /= pivot
		}

		for j := i + 1; j < n; j++ {
			factor := augmented[j][i]
			for k := 0; k < 2*n; k++ {
				augmented[j][k] -= factor * augmented[i][k]
			}
		}
	}

	for i := n - 1; i >= 0; i-- {
		for j := i - 1; j >= 0; j-- {
			factor := augmented[j][i]
			for k := 0; k < 2*n; k++ {
				augmented[j][k] -= factor * augmented[i][k]
			}
		}
	}

	inverse := make([][]float64, n)
	for i := range inverse {
		inverse[i] = make([]float64, n)
		copy(inverse[i], augmented[i][n:])
	}

	return inverse
}

// HaversineDistance returns the distance between two points on the Earth's surface.
// lat1: latitude of the first point.
// lon1: longitude of the first point.
// lat2: latitude of the second point.
// lon2: longitude of the second point.
// returns: a float64 - the distance between the two points.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// Mean returns the mean of a slice of float64.
// data: a slice of float64.
// returns: a float64 - the mean of the data.
func Mean(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum / float64(len(data))
}

// StdDev returns the standard deviation of a slice of float64.
// data: a slice of float64.
// returns: a float64 - the standard deviation of the data.
func StdDev(data []float64) float64 {
	mean := Mean(data)
	sumSquares := 0.0
	for _, value := range data {
		diff := value - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}

// Square returns the square of each element in a slice of float64.
// data: a slice of float64.
// returns: a slice of float64.
func Square(data []float64) []float64 {
	result := make([]float64, len(data))
	for i, value := range data {
		result[i] = value * value
	}
	return result
}

// Multiply returns the element-wise product of two slices of float64.
// a: a slice of float64.
// b: a slice of float64.
// returns: a slice of float64.
func Multiply(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range result {
		result[i] = a[i] * b[i]
	}
	return result
}

// Sum returns the sum of a slice of float64.
// data: a slice of float64.
// returns: a float64 - the sum of the data.
func Sum(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	return sum
}
