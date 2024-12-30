package ordmap

// Set sets a value in the map, assigning it the highest index + 1.
func (om *OrderedMap[K, V]) Set(key K, v V) {
	Set(om, key, Value[V]{V: v}, getIndex[V], setIndex[V])
}

// Set is a helper function to set a value in the map, assigning it the highest index + 1.
func Set[M ~map[K]V, K comparable, V any](
	m *M, key K, v V,
	getIndex func(V) int,
	setIndex func(*V, int),
) {
	// check if the map is nil and create it if it is
	if *m == nil {
		setIndex(&v, 1)
		*m = M{key: v}
		return
	}

	highestIdx := 0
	for _, v := range *m {
		if idx := getIndex(v); idx > highestIdx {
			highestIdx = idx
		}
	}

	setIndex(&v, highestIdx+1)
	(*m)[key] = v
}