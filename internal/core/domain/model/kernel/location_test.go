package kernel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Location_Valid(t *testing.T) {
	// Data
	x, y := 5, 7

	// Steps
	location, err := NewLocation(x, y)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, x, location.X())
	assert.Equal(t, y, location.Y())
}

func Test_Location_Error_OutOfBounds(t *testing.T) {
	tests := map[string]struct {
		x        int
		y        int
		expected error
	}{
		"x too low":    {x: 0, y: 5, expected: ErrorInvalidCoordinate},
		"x too high":   {x: 11, y: 5, expected: ErrorInvalidCoordinate},
		"y too low":    {x: 5, y: 0, expected: ErrorInvalidCoordinate},
		"y too high":   {x: 5, y: 11, expected: ErrorInvalidCoordinate},
		"both invalid": {x: 0, y: 11, expected: ErrorInvalidCoordinate},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Date
			_, err := NewLocation(test.x, test.y)

			// Assert
			assert.EqualError(t, err, test.expected.Error())
		})
	}
}

func Test_Location_DistanceTo(t *testing.T) {
	// Data
	location1, _ := NewLocation(2, 3)
	location2, _ := NewLocation(5, 6)

	// Steps
	distance := location1.DistanceTo(location2)

	// Assert
	assert.Equal(t, 6, distance) // |5-2| + |6-3| = 3 + 3 = 6
}

func Test_Location_Equals(t *testing.T) {
	tests := map[string]struct {
		x1, y1 int
		x2, y2 int
		want   bool
	}{
		"equal coordinates":     {x1: 3, y1: 3, x2: 3, y2: 3, want: true},
		"different coordinates": {x1: 3, y1: 3, x2: 4, y2: 3, want: false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Steps
			location1, _ := NewLocation(test.x1, test.y1)
			location2, _ := NewLocation(test.x2, test.y2)

			// Assert
			assert.Equal(t, test.want, location1.Equals(location2))
		})
	}
}

func Test_Location_CreateRandom(t *testing.T) {
	// Steps
	randomLocation, _ := CreateRandom()

	// Assert
	assert.GreaterOrEqual(t, randomLocation.X(), 1)
	assert.LessOrEqual(t, randomLocation.X(), 10)
	assert.GreaterOrEqual(t, randomLocation.Y(), 1)
	assert.LessOrEqual(t, randomLocation.Y(), 10)
}
