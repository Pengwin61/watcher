package utils

func RemoveSlice[T comparable](slice []T, i int) []T {

	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
