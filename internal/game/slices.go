package game

func containsInt(slice []int, v int) bool {
	for _, x := range slice {
		if x == v {
			return true
		}
	}
	return false
}

func inRange(idx, length int) bool {
	return idx >= 0 && idx < length
}

func removeInt(slice []int, v int) []int {
	for i, x := range slice {
		if x == v {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func appendIntIfMissing(slice []int, v int) []int {
	if containsInt(slice, v) {
		return slice
	}
	return append(slice, v)
}

func validTerritoryID(territories []*Territory, id int) bool {
	return inRange(id, len(territories))
}
