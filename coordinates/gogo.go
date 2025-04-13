package coordinates

import "github.com/zachklingbeil/factory"

type Coordinates struct {
	Coords    []Coord `json:"coordinates"`
	Factory   *factory.Factory
	NestedMap map[uint8]map[uint8]map[uint8]map[uint8]map[uint8]map[uint8]map[uint16]map[uint16]any `json:"nested_map"`
}

type Coord struct {
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
}

func NewCoordinates(factory *factory.Factory) *Coordinates {
	// Define the levels for the nested map
	levels := []int{256, 12, 31, 24, 60, 60, 1000, 1000} // Corresponding to year, month, day, hour, minute, second, millisecond, index

	// Initialize the nested map using the helper functions
	nestedMap := initializeNestedMap(levels).(map[uint8]map[uint8]map[uint8]map[uint8]map[uint8]map[uint8]map[uint16]map[uint16]any)

	return &Coordinates{
		Coords:    make([]Coord, 0),
		NestedMap: nestedMap,
		Factory:   factory,
	}
}

// Helper function to initialize 0-based indexes
func indexZero(levels []int) any {
	if len(levels) == 0 {
		return nil
	}
	currentLevel := make(map[uint16]any)
	for i := uint16(0); i < uint16(levels[0]); i++ {
		currentLevel[i] = indexZero(levels[1:])
	}
	return currentLevel
}

// Helper function to initialize 1-based indexes
func indexOne(levels []int) any {
	if len(levels) == 0 {
		return nil
	}
	currentLevel := make(map[uint16]any)
	for i := uint16(1); i <= uint16(levels[0]); i++ {
		currentLevel[i] = indexOne(levels[1:])
	}
	return currentLevel
}

// Main helper function to initialize the nested map structure
func initializeNestedMap(levels []int) any {
	if len(levels) == 0 {
		return nil
	}
	if len(levels) == 8 { // Year
		return indexZero(levels)
	} else if len(levels) == 7 || len(levels) == 6 { // Month or Day
		return indexOne(levels)
	} else { // Hour, Minute, Second, Millisecond, Index
		return indexZero(levels)
	}
}
