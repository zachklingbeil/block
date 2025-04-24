package circuit

func (c *Circuit) Remove(key any) {
	delete(c.Map, key)
}
func (c *Circuit) Clear() {
	c.Map = make(map[any]any)
}
func (c *Circuit) Size() int {
	return len(c.Map)
}
func (c *Circuit) Contains(key any) bool {
	_, ok := c.Map[key]
	return ok
}
func (c *Circuit) Keys() []any {
	keys := make([]any, 0, len(c.Map))
	for key := range c.Map {
		keys = append(keys, key)
	}
	return keys
}
func (c *Circuit) Values() []any {
	values := make([]any, 0, len(c.Map))
	for _, value := range c.Map {
		values = append(values, value)
	}
	return values
}
func (c *Circuit) ForEach(fn func(key any, value any)) {
	for key, value := range c.Map {
		fn(key, value)
	}
}

func (c *Circuit) Filter(fn func(key any, value any) bool) map[any]any {
	filtered := make(map[any]any)
	for key, value := range c.Map {
		if fn(key, value) {
			filtered[key] = value
		}
	}
	return filtered
}
func (c *Circuit) MapKeys(fn func(key any) any) []any {
	keys := make([]any, 0, len(c.Map))
	for key := range c.Map {
		keys = append(keys, fn(key))
	}
	return keys
}
func (c *Circuit) MapValues(fn func(value any) any) []any {
	values := make([]any, 0, len(c.Map))
	for _, value := range c.Map {
		values = append(values, fn(value))
	}
	return values
}
func (c *Circuit) MapEntries(fn func(key any, value any) (any, any)) map[any]any {
	entries := make(map[any]any)
	for key, value := range c.Map {
		newKey, newValue := fn(key, value)
		entries[newKey] = newValue
	}
	return entries
}
func (c *Circuit) Reduce(fn func(acc any, key any, value any) any, initialValue any) any {
	acc := initialValue
	for key, value := range c.Map {
		acc = fn(acc, key, value)
	}
	return acc
}
func (c *Circuit) GroupBy(fn func(key any, value any) any) map[any][]any {
	grouped := make(map[any][]any)
	for key, value := range c.Map {
		groupKey := fn(key, value)
		grouped[groupKey] = append(grouped[groupKey], value)
	}
	return grouped
}
func (c *Circuit) Sort(fn func(a any, b any) bool) []any {
	sorted := make([]any, 0, len(c.Map))
	for _, value := range c.Map {
		sorted = append(sorted, value)
	}
	for i := range len(sorted) - 1 {
		for j := i + 1; j < len(sorted); j++ {
			if fn(sorted[i], sorted[j]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
func (c *Circuit) SortByKey(fn func(a any, b any) bool) []any {
	sorted := make([]any, 0, len(c.Map))
	for key := range c.Map {
		sorted = append(sorted, key)
	}
	for i := range len(sorted) - 1 {
		for j := i + 1; j < len(sorted); j++ {
			if fn(sorted[i], sorted[j]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
func (c *Circuit) SortByValue(fn func(a any, b any) bool) []any {
	sorted := make([]any, 0, len(c.Map))
	for _, value := range c.Map {
		sorted = append(sorted, value)
	}
	for i := range len(sorted) - 1 {
		for j := i + 1; j < len(sorted); j++ {
			if fn(sorted[i], sorted[j]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
func (c *Circuit) SortByKeyValue(fn func(a any, b any) bool) []any {
	sorted := make([]any, 0, len(c.Map))
	for key, value := range c.Map {
		sorted = append(sorted, key, value)
	}
	for i := range len(sorted) - 1 {
		for j := i + 1; j < len(sorted); j++ {
			if fn(sorted[i], sorted[j]) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
