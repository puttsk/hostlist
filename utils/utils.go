package utils

// CartesianProduct creates a list of Cartesian products from a set of arrays.
// Returns an empty list if `a` is `nil` or empty list.
//
// For example:
//
//	Cartesian product of `[[0,1]]` is [[0], [1]]
//	Cartesian product of `[[0,1], [a,b]]` is [[0,a], [0,b], [1,a], [1,b]]
func CartesianProduct[T any](a [][]T) [][]T {
	// If 'a' is an empty array, return empty list
	if len(a) == 0 {
		return [][]T{}
	}

	product := make([][]T, len(a[0]))
	for i, v := range a[0] {
		product[i] = []T{v}
	}
	//product := [][]T{{}}
	for _, next := range a[1:] {
		nextLen := len(next)
		newProduct := make([][]T, len(product)*nextLen)
		for i, c := range product {
			for j, v := range next {
				newProduct[i*nextLen+j] = append(c, v)
			}
		}
		product = newProduct
	}
	return product
}
