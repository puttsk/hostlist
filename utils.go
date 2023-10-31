package hostlist

// CartesianProduct creates a Cartesian product of lists 'a'.
//
// For example:
//	Cartesian product of `[[0,1]]` is [[0], [1]]
//	Cartesian product of `[[0,1], [a,b]]` is [[0,a], [0,b], [1,a], [1,b]]
//
func CartesianProduct(a [][]interface{}) [][]interface{} {
	permutation := [][]interface{}{}
	// a = [0,1]
	// b = [[a],[b]]
	if len(a) == 1 {
		for _, first := range a[0] {
			permutation = append(permutation, []interface{}{first})
		}
	} else {
		b := CartesianProduct(a[1:])
		for _, first := range a[0] {
			for _, second := range b {
				permutation = append(permutation, append([]interface{}{first}, second...))
			}
		}
	}
	return permutation
}
