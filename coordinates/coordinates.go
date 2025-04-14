package coordinates

import (
	"fmt"
	"time"
)

// Creates map for each Coord field, preserving the int/uint types, data will be assigned to the innermost map
func (c *Coordinates) CreateNestedMap(data any) {
	for _, coord := range c.Coords {
		c.NestedMap[coord.Year][coord.Month][coord.Day][coord.Hour][coord.Minute][coord.Second][coord.Millisecond][coord.Index] = data
	}
}

// Add a coordinate to the in-memory slice
func (c *Coordinates) AddCoordinate(timestamp int64, index int64) {
	t := time.UnixMilli(timestamp)
	coord := Coord{
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       uint16(index),
	}
	c.Coords = append(c.Coords, coord)
	fmt.Printf("%d.%d.%d.%d.%d.%d.%d.%d\n",
		coord.Year,
		coord.Month,
		coord.Day,
		coord.Hour,
		coord.Minute,
		coord.Second,
		coord.Millisecond,
		coord.Index,
	)
}
