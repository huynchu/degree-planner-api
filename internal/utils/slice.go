package utils

func Insert[T any](slice []T, index int, value T) []T {
	// Expand the slice by one element
	var tmp T
	slice = append(slice, tmp)

	// Shift the elements after the insertion point
	copy(slice[index+1:], slice[index:])

	// Insert the value at the specified index
	slice[index] = value

	return slice
}

func Move[T any](slice []T, from int, to int) []T {
	if from == to {
		return slice
	}

	value := slice[from]

	// Remove the element at the from index
	ret := make([]T, 0)
	ret = append(ret, slice[:from]...)
	ret = append(ret, slice[from+1:]...)

	// Insert the element at the to index
	ret = Insert(ret, to, value)

	return ret
}
